package parser

import (
	"errors"
	"fmt"

	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/up"
	"github.com/reverbdotcom/sbx/name"
)

const help = `
Usage: sbx <up|down|name>
`

func Parse(args []string) (*cli.CmdFn, error) {
	cmd, err := command(args)

	if err != nil {
		return nil, err
	}

	return cmdfn(*cmd)
}

func command(args []string) (command *string, err error) {
	if len(args) < 2 {
		return nil, errr("command required")
	}

	cmd := args[1]

	return &cmd, nil
}

func cmdfn(command string) (*cli.CmdFn, error) {
	cmd, ok := commands()[command]

	if !ok {
		return nil, errr("command not found")
	}

	return &cmd, nil
}


func commands() map[string]cli.CmdFn {
	return map[string]cli.CmdFn{
		"up": up.Run,
		"name": name.Run,
	}
}

func errr(message string) error {
	return errors.New(fmt.Sprintf("🚫 %s\n%s", message, help))
}
