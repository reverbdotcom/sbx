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
	validate()

	cmdfn, err := parser.Parse(os.Args)
	onError(err)

	fn := *cmdfn
	out, err := fn()
	onError(err)

	fmt.Println(out)
}

func validate() {
	has, err := check.OnOrchestra()

	onError(err)

	if !has {
		onError(errors.New("This project is not on Orchestra."))
	}

	if !check.HasGithubToken() {
		onError(errors.New("Please set the GITHUB_TOKEN environment variable."))
	}
}

func onError(err error) {
	if err != nil {
		fmt.Println(errr.New(err.Error()))
		os.Exit(1)
	}
}
