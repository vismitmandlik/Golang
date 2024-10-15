package main

import (
	"fmt"
)

func main() {
	// Create an empty slice of integers
	var numbers []int

	// Append elements to the slice
	numbers = append(numbers, 1)
	numbers = append(numbers, 2, 3, 4)

	// Print the slice
	fmt.Println("Slice:", numbers)

	// Access elements from the slice
	fmt.Println("First element:", numbers[0])
	fmt.Println("Second element:", numbers[1])

	// Modify an element in the slice
	numbers[1] = 20
	fmt.Println("Modified Slice:", numbers)

	// Length and capacity of the slice
	fmt.Println("Length of slice:", len(numbers))
	fmt.Println("Capacity of slice:", cap(numbers))

	// Slicing the slice
	subSlice := numbers[1:3]
	fmt.Println("Sub-slice:", subSlice)

	// Iterate over the slice
	fmt.Println("Iterating over slice:")
	for index, value := range numbers {
		fmt.Printf("Index: %d, Value: %d\n", index, value)
	}
}
