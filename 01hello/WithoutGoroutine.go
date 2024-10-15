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

	// Sequential execution (without goroutines)
	task("Task 1", 2*time.Second)
	task("Task 2", 1*time.Second)
	task("Task 3", 3*time.Second)

	elapsed := time.Since(start)
	fmt.Printf("All tasks completed sequentially in %s\n", elapsed)
}
