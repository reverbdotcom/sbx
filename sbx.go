package main

import (
	"fmt"
	"os"
)

func main() {
	// os.Args[0] is the program name itself
	if len(os.Args) < 2 {
		fmt.Println("Please provide at least one argument")
		os.Exit(1)
	}

	// Access the first positional argument
	firstArg := os.Args[1]
	fmt.Println("First argument:", firstArg)

	// Access all positional arguments (excluding the program name)
	allArgs := os.Args[1:]
	fmt.Println("All arguments:", allArgs)
}
