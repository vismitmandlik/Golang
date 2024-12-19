package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// Shared data structures
type Credentials struct {
	Username string `json:"username"`

	Password string `json:"password"`
}

type Device struct {
	ID string `json:"_id"`

	IP string `json:"ip"`

	Credentials Credentials `json:"credentials"`

	Port int `json:"port"`
}

type Metrics struct {
	DeviceID string `json:"deviceId"`

	IP string `json:"ip"`

	CPUUsage string `json:"cpuUsage"`

	MemoryUsage string `json:"memoryUsage"`

	DiskUsage string `json:"diskUsage"`

	Timestamp int64 `json:"timestamp"`
}

type IpCredential struct {
	Ip string `json:"ip"`

	Port int `json:"port"`

	Credentials []Credentials `json:"credentials"`
}

// Utility functions
func parseCommandLineArgs() (string, string, error) {

	if len(os.Args) < 3 {

		return "", "", fmt.Errorf("insufficient arguments provided")
	}

	eventName := os.Args[1]

	input := os.Args[2]

	return eventName, input, nil
}

// Poller logic
func runPoller(devices []Device) {

	var wg sync.WaitGroup

	for _, device := range devices {

		wg.Add(1)

		go fetchMetrics(device, &wg)
	}

	wg.Wait()
}

func fetchMetrics(device Device, wg *sync.WaitGroup) {

	defer wg.Done()

	client, err := sshConnect(device)

	if err != nil {

		fmt.Printf("Error connecting to device %s: %v\n", device.ID, err)

		return
	}

	defer client.Close()

	cpuUsage, err := getCPUUsage(device, client)

	if err != nil {

		fmt.Printf("Error getting CPU usage for device %s: %v\n", device.ID, err)

		return
	}

	memoryUsage, err := getMemoryUsage(device, client)

	if err != nil {

		fmt.Printf("Error getting memory usage for device %s: %v\n", device.ID, err)

		return
	}

	diskUsage, err := getDiskUsage(device, client)

	if err != nil {

		fmt.Printf("Error getting disk usage for device %s: %v\n", device.ID, err)

		return
	}

	metrics := Metrics{

		DeviceID: device.ID,

		IP: device.IP,

		CPUUsage: cpuUsage,

		MemoryUsage: memoryUsage,

		DiskUsage: diskUsage,

		Timestamp: time.Now().UnixMilli(),
	}

	metricsJSON, _ := json.Marshal(metrics)

	fmt.Println(string(metricsJSON))
}

// Discovery logic
func runDiscovery(ipCred IpCredential) {

	var wg sync.WaitGroup

	var mux sync.Mutex

	successfulIPs := make(map[string]bool)

	for _, cred := range ipCred.Credentials {

		wg.Add(1)

		go func(ip string, port int, username, password string) {

			defer wg.Done()

			if success := trySSH(ip, port, username, password); success {

				mux.Lock()

				defer mux.Unlock()

				if !successfulIPs[ip] {

					successfulIPs[ip] = true

					successfulCredential := Credentials{Username: username, Password: password}

					jsonData, _ := json.Marshal(successfulCredential)

					fmt.Printf("Successful login for IP %s: %s\n", ip, string(jsonData))
				}
			}
		}(ipCred.Ip, ipCred.Port, cred.Username, cred.Password)
	}

	wg.Wait()

	fmt.Println("Processing complete.")
}

// Shared SSH utilities
func sshConnect(device Device) (*ssh.Client, error) {

	cfg := &ssh.ClientConfig{

		User: device.Credentials.Username,

		Auth: []ssh.AuthMethod{

			ssh.Password(device.Credentials.Password),
		},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", device.IP, device.Port), cfg)

	if err != nil {

		return nil, fmt.Errorf("failed to connect to device %s: %w", device.ID, err)
	}

	return client, nil
}

func trySSH(ip string, port int, username, password string) bool {

	config := &ssh.ClientConfig{

		User: username,

		Auth: []ssh.AuthMethod{

			ssh.Password(password),
		},

		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", ip, port)

	client, err := ssh.Dial("tcp", addr, config)

	if err != nil {

		return false
	}

	defer client.Close()

	return true
}

func getCPUUsage(device Device, client *ssh.Client) (string, error) {

	cmd := "top -b -n 1 | grep 'Cpu(s)' | awk '{usage=100-$8; printf(\"%.2f\\n\", usage)}'"

	return runSshCommand(client, cmd)
}

func getMemoryUsage(device Device, client *ssh.Client) (string, error) {

	cmd := "free | grep Mem | awk '{usage=($3/$2)*100; printf(\"%.2f\\n\", usage)}'"

	return runSshCommand(client, cmd)
}

func getDiskUsage(device Device, client *ssh.Client) (string, error) {

	cmd := "df --total | tail -1 | awk '{print $5}'"

	return runSshCommand(client, cmd)
}

func runSshCommand(client *ssh.Client, cmd string) (string, error) {

	session, err := client.NewSession()

	if err != nil {

		return "", fmt.Errorf("failed to create session: %v", err)
	}

	defer session.Close()

	output, err := session.CombinedOutput(cmd)

	if err != nil {

		return "", fmt.Errorf("error running cmd '%s': %s", cmd, err)
	}

	return strings.TrimSpace(string(output)), nil
}

func main() {

	eventName, input, err := parseCommandLineArgs()

	if err != nil {

		fmt.Printf("Error: %v\n", err)

		return
	}

	switch eventName {

	case "poller":

		var devices []Device

		err := json.Unmarshal([]byte(input), &devices)

		if err != nil {

			fmt.Printf("Error parsing devices: %v\n", err)

			return
		}

		runPoller(devices)

	case "discovery":

		var ipCred IpCredential

		err := json.Unmarshal([]byte(input), &ipCred)

		if err != nil {

			fmt.Printf("Error parsing IP credentials: %v\n", err)

			return
		}

		runDiscovery(ipCred)

	default:

		fmt.Printf("Unknown event name: %s\n", eventName)
	}
}