package main

import (
	"fmt"
	"os"

	"github.com/reverbdotcom/sbx/env"
	"github.com/reverbdotcom/sbx/errr"
	"github.com/reverbdotcom/sbx/parser"
	"github.com/reverbdotcom/sbx/tui"
)

// tuiCommands are commands handled by the TUI instead of the legacy CLI.
var tuiCommands = map[string]bool{
	"up": true, "u": true,
	"version": true, "v": true,
	"name": true, "n": true,
	"help": true, "h": true,
	"info": true, "i": true,
}

func main() {
	err := env.Verify()
	onError(err)

	// No args: launch TUI with main menu.
	if len(os.Args) < 2 {
		err := tui.Run("")
		onError(err)
		return
	}

	cmd := os.Args[1]

	// TUI-handled commands: launch TUI with that command's output,
	// then return to the main menu.
	if tuiCommands[cmd] {
		err := tui.Run(cmd)
		onError(err)
		return
	}

	// All other commands: use existing CLI parser.
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
