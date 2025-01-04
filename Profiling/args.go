package main

import (
	"fmt"
	"os"
)

func args() {
	if len(os.Args) < 2 {
		fmt.Println("No arguments provided")
		return
	}
	fmt.Printf("hello %s\n", os.Args[1])
}
