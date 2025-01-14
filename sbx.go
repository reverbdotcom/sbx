package main

import (
	"fmt"
	"os"

	"github.com/reverbdotcom/sbx/env"
	"github.com/reverbdotcom/sbx/errr"
	"github.com/reverbdotcom/sbx/parser"
)

func main() {
	err := env.Verify()
	onError(err)

	cmdfn, err := parser.Parse(os.Args)
	onError(err)

	fn := *cmdfn
	out, err := fn()
	fmt.Println(out)
	onError(err)
}

func onError(err error) {
	if err != nil {
		fmt.Println(errr.New(err.Error()))
		os.Exit(1)
	}
}
