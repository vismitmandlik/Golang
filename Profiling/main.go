package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {

	go func() {

		log.Println(http.ListenAndServe("localhost:8080", nil))
	}()

	// Simulate some work in the main program
	for i := 0; i < 1000000; i++ {

		_ = i * 2

		time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("Main function work done.")
}
