package parser

import (
	"errors"

	"github.com/reverbdotcom/sbx/commands"
)

func Parse(args []string) (*commands.RunFn, error) {
	cmd, err := command(args)

	if err != nil {
		return nil, err
	}

	return cmdfn(*cmd)
}

func command(args []string) (command *string, err error) {
	if len(args) < 2 {
		return nil, errors.New("command required")
	}

	cmd := args[1]

	return &cmd, nil
}

func cmdfn(command string) (*commands.RunFn, error) {
	cmd, ok := commands.Commands()[command]

	if !ok {
		return nil, errors.New("command not found")
	}

	return &cmd, nil
}
