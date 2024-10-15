package main

import (
    "fmt"
    "time"
)

func main() {
    stop := make(chan bool) // Channel to signal the goroutine to stop

    go func() {
        for {
            select {
            case <-stop:
                fmt.Println("Goroutine stopped")
                return // Exit the goroutine
            default:
                fmt.Println("Goroutine running")
                time.Sleep(500 * time.Millisecond) // Simulate work
            }
        }
    }()

    time.Sleep(2 * time.Second) // Let the goroutine run for a while

    // Send signal to stop the goroutine conditionally
    stop <- true 

	message := "Hello motadata !"
	fmt.Println(message)
    time.Sleep(1 * time.Second) // Give time for goroutine to stop
    fmt.Println("Main function finished")
}
