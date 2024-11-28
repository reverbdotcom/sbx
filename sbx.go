package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/reverbdotcom/sbx/check"
	"github.com/reverbdotcom/sbx/errr"
	"github.com/reverbdotcom/sbx/parser"
)

func main() {
	has, err := check.OnOrchestra()

	onError(err)

	if !has {
		onError(errors.New("This project is not on Orchestra."))
	}

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
