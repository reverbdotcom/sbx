package main

import (
	"fmt"
	"os"


  "github.com/reverbdotcom/sbx/up"
)

func main() {
	// os.Args[0] is the program name itself
	if len(os.Args) < 2 {
		fmt.Println("Please provide at least one argument")
		os.Exit(1)
	}

	// Access the first positional argument
	command := os.Args[1]
	fmt.Println("cmd:", command)


  up.Test()
}
