package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Device struct {
	ID          string      `json:"_id"`
	IP          string      `json:"ip"`
	Credentials Credentials `json:"credentials"`
	Port        int         `json:"port"`
}

type Metrics struct {
	DeviceID    string `json:"deviceId"`
	CPUUsage    string `json:"cpuUsage"`
	MemoryUsage string `json:"memoryUsage"`
	DiskUsage   string `json:"diskUsage"`
	Timestamp   int64  `json:"timestamp"`
}

// Parse command-line arguments to read device details
func parseCommandLineArgs() ([]Device, error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("no devices passed as arguments")
	}

	input := os.Args[1]
	var devices []Device

	err := json.Unmarshal([]byte(input), &devices)
	if err != nil {
		return nil, fmt.Errorf("failed to parse device details: %v", err)
	}
	return devices, nil
}

// Run SSH command and capture output
func runSshCommand(client *ssh.Client, cmd string) (string, error) {
	// Create a new session for each command
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", fmt.Errorf("error running cmd '%s': %s", cmd, err)
	}
	return string(output), nil
}

// Establish SSH connection
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

// Get CPU usage percentage using the provided command
func getCPUUsage(device Device, client *ssh.Client) (string, error) {
	cmd := "top -b -n 1 | grep 'Cpu(s)' | awk '{usage=100-$8; printf(\"%.2f\\n\", usage)}'"
	output, err := runSshCommand(client, cmd)
	if err != nil {
		return "", fmt.Errorf("error getting CPU usage for device %s: %v", device.ID, err)
	}
	return strings.TrimSpace(output), nil
}

// Get memory usage percentage using the provided command
func getMemoryUsage(device Device, client *ssh.Client) (string, error) {
	cmd := "free | grep Mem | awk '{usage=($3/$2)*100; printf(\"%.2f\\n\", usage)}'"
	output, err := runSshCommand(client, cmd)
	if err != nil {
		return "", fmt.Errorf("error getting memory usage for device %s: %v", device.ID, err)
	}
	return strings.TrimSpace(output), nil
}

// Get disk usage percentage using the provided command
func getDiskUsage(device Device, client *ssh.Client) (string, error) {
	cmd := "df --total | tail -1 | awk '{print $5}'"
	output, err := runSshCommand(client, cmd)
	if err != nil {
		return "", fmt.Errorf("error getting disk usage for device %s: %v", device.ID, err)
	}
	return strings.TrimSpace(output), nil
}

// Fetch metrics for a single device
func fetchMetrics(device Device, wg *sync.WaitGroup) {
	defer wg.Done()

	// Establish SSH connection
	client, err := sshConnect(device)
	if err != nil {
		log.Printf("Error connecting to device %s: %v", device.ID, err)
		return
	}
	defer client.Close()

	// Get CPU usage
	cpuUsage, err := getCPUUsage(device, client)
	if err != nil {
		log.Printf("Error getting CPU usage for device %s: %v", device.ID, err)
		return
	}

	// Get memory usage
	memoryUsage, err := getMemoryUsage(device, client)
	if err != nil {
		log.Printf("Error getting memory usage for device %s: %v", device.ID, err)
		return
	}

	// Get disk usage
	diskUsage, err := getDiskUsage(device, client)
	if err != nil {
		log.Printf("Error getting disk usage for device %s: %v", device.ID, err)
		return
	}

	// Prepare metrics structure
	metrics := Metrics{
		DeviceID:    device.ID,
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		DiskUsage:   diskUsage,
		Timestamp:   time.Now().UnixMilli(),
	}

	// Print metrics as JSON
	metricsJSON, _ := json.Marshal(metrics)
	fmt.Println(string(metricsJSON))
}

func main() {
	// Parse devices from command-line arguments
	devices, err := parseCommandLineArgs()
	if err != nil {
		log.Fatalf("Error parsing command-line arguments: %v", err)
	}

	var wg sync.WaitGroup

	// Process each device concurrently using goroutines
	for _, device := range devices {
		wg.Add(1)
		go fetchMetrics(device, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("All devices processed.")
}
