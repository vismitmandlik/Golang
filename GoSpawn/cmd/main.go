// package main

// import (
// 	"fmt"
// 	"time"
// )

// // fetchDeviceMetrics simulates fetching CPU metrics from a device.
// func fetchDeviceMetrics() int {
// 	return 75 // Example CPU usage percentage
// }

// func main() {
// 	ticker := time.NewTicker(4 * time.Second) // Poll every 10 seconds
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			cpuUsage := fetchDeviceMetrics()
// 			fmt.Printf("CPU Usage: %d%%\n", cpuUsage) // Print CPU metrics to standard output
// 		}
// 	}
// }
