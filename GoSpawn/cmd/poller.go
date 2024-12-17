package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
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
	DeviceID    string  `json:"deviceId"`
	CPUUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
	DiskUsage   float64 `json:"diskUsage"`
	Timestamp   int64   `json:"timestamp"`
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

// Establish SSH connection
func sshConnect(device Device) (*ssh.Session, error) {
	config := &ssh.ClientConfig{
		User: device.Credentials.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(device.Credentials.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", device.IP, device.Port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to device %s: %w", device.ID, err)
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session for device %s: %w", device.ID, err)
	}

	return session, nil
}

// Parse the `top` output to extract CPU, Memory, and Disk usage metrics
func parseTopOutput(output string) (float64, float64, float64) {
	var cpuUsage, memoryUsage, diskUsage float64
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "%Cpu(s):") {
			// Extract idle percentage and calculate CPU usage
			parts := strings.Fields(line)
			if len(parts) > 3 {
				idle, _ := strconv.ParseFloat(strings.TrimSuffix(parts[3], ","), 64)
				cpuUsage = 100.0 - idle
			}
		}

		if strings.HasPrefix(line, "MiB Mem :") {
			// Extract memory usage
			parts := strings.Fields(line)
			if len(parts) > 6 {
				totalMem, _ := strconv.ParseFloat(parts[1], 64)
				usedMem, _ := strconv.ParseFloat(parts[5], 64)
				if totalMem > 0 {
					memoryUsage = (usedMem / totalMem) * 100
				}
			}
		}

		if strings.HasPrefix(line, "MiB Swap:") {
			// Extract disk usage under Swap
			parts := strings.Fields(line)
			if len(parts) > 6 {
				totalSwap, _ := strconv.ParseFloat(parts[1], 64)
				usedSwap, _ := strconv.ParseFloat(parts[5], 64)
				if totalSwap > 0 {
					diskUsage = (usedSwap / totalSwap) * 100
				}
			}
		}
	}
	return cpuUsage, memoryUsage, diskUsage
}

// Fetch metrics for a single device
func fetchMetrics(device Device, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Connecting to device: %s (%s)", device.ID, device.IP)

	session, err := sshConnect(device)
	if err != nil {
		log.Printf("Error connecting to device %s: %v", device.ID, err)
		return
	}
	defer session.Close()

	// Run the `top` command to fetch metrics
	output, err := session.CombinedOutput("top -b -n 1")
	if err != nil {
		log.Printf("Error executing top command for device %s: %v", device.ID, err)
		return
	}
	fmt.Print("Output of top is ", output)

	// Parse the metrics from `top` output
	cpuUsage, memoryUsage, diskUsage := parseTopOutput(string(output))

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
