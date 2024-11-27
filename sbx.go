package main

import (
	"fmt"
	"os"

	"github.com/reverbdotcom/sbx/errr"
	"github.com/reverbdotcom/sbx/parser"
)

func main() {
	cmdfn, err := parser.Parse(os.Args)

	if err != nil {
		fmt.Println(errr.New(err.Error()))
		os.Exit(1)
	}

	fn := *cmdfn
	out, err := fn()
	fmt.Println(out)

	if err != nil {
		fmt.Println(errr.New(err.Error()))
		os.Exit(1)
	}
}
