package main

import (
	"fmt"
	"time"
)

func task(name string, duration time.Duration) {
	fmt.Printf("%s starting...\n", name)
	time.Sleep(duration) // Simulate work with a sleep
	fmt.Printf("%s finished.\n", name)
}

func main() {
	start := time.Now()

	// Concurrent execution with goroutines
	go task("Task 1", 2*time.Second)
	go task("Task 2", 1*time.Second)
	go task("Task 3", 3*time.Second)

	// Wait for all tasks to complete
	time.Sleep(4 * time.Second) // Sleep enough time for all tasks to finish

	elapsed := time.Since(start)
	fmt.Printf("All tasks completed concurrently in %s\n", elapsed)
}
