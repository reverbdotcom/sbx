package parser

import (
	"errors"
	"os"

	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/up"
)

func Parse() (*cli.CmdFn, error) {
	cmd, err := command()

	if err != nil {
		return nil, err
	}

	return cmdfn(*cmd)
}

func command() (command *string, err error) {
	if len(os.Args) < 2 {
		return nil, errors.New("Please provide at least one argument")
	}

	cmd := os.Args[1]

	return &cmd, nil
}

func cmdfn(command string) (*cli.CmdFn, error) {
	cmd, ok := map[string]cli.CmdFn{
		"up": up.Run,
	}[command]

	if !ok {
		return nil, errors.New("Command not found")
	}

	return &cmd, nil
}
