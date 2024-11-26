package main

import (
	"fmt"
	"os"

	"github.com/reverbdotcom/sbx/parser"
)

func main() {
	cmdfn, err := parser.Parse()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fn := *cmdfn

	fn()
}
