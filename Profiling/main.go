package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof" // Import pprof to register its handlers
	"time"
)

func main() {
	// Start an HTTP server for pprof
	go func() {
		fmt.Println("Starting pprof server on localhost:6060")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Simulate some work
	for i := 0; i < 1000000; i++ {
		_ = i * 2
		time.Sleep(1 * time.Millisecond) // Simulate some delay
	}

	fmt.Println("Main function work done.")
}
