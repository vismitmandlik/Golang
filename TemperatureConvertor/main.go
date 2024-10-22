package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
)

func main() {
	// Create a file to store the profile
	f, err := os.Create("mypro.prof")
	if err != nil {
		fmt.Println("Error creating profile:", err)
		return
	}
	pprof.StartCPUProfile(f)     // Start CPU profiling
	defer pprof.StopCPUProfile() // Stop when main exits

	var input int
	fmt.Print("Please enter Temperature in degrees: ")
	fmt.Scanf("%d", &input)
	go func() {

		log.Println(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", 8080), nil))
	}()
	output := input * 100
	fmt.Println("Temperature is:", output)

	// Simulate some work for profiling purposes
	for i := 0; i < 10000000; i++ {
		_ = i * 2
	}
	fmt.Println("Work done.")
}
