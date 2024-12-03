package main

import (
	"fmt"
	"os"

	"github.com/reverbdotcom/sbx/errr"
	"github.com/reverbdotcom/sbx/parser"
)

func main() {
	cmdfn, err := parser.Parse(os.Args)
	onError(err)

	fn := *cmdfn
	out, err := fn()
	onError(err)

	fmt.Println(out)
}

func onError(err error) {
	if err != nil {
		fmt.Println(errr.New(err.Error()))
		os.Exit(1)
	}
}
