package main

import "fmt"

func main() {
	// Create a slice with a length of 3 and capacity of 5
	slice := make([]int, 3, 5)

	fmt.Println("Initial slice:", slice)
	fmt.Println("Length:", len(slice))    // 3
	fmt.Println("Capacity:", cap(slice))  // 5

	// Append more elements to the slice
	slice = append(slice, 10, 20)

	fmt.Println("Slice after appending:", slice)
	fmt.Println("Length after appending:", len(slice))   // 5
	fmt.Println("Capacity after appending:", cap(slice)) // 5

	// Append one more element, exceeding the current capacity
	slice = append(slice, 30)

	fmt.Println("Slice after exceeding capacity:", slice)
	fmt.Println("Length now:", len(slice))   // 6
	fmt.Println("Capacity now:", cap(slice)) // Capacity will increase (doubled)
}
