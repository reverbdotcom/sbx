package parser

import (
	"github.com/reverbdotcom/sbx/commands"
	"github.com/reverbdotcom/sbx/errr"
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
		return nil, errr.NewWithHelp("command required")
	}

	cmd := args[1]

	return &cmd, nil
}

func cmdfn(command string) (*commands.RunFn, error) {
	cmd, ok := commands.Commands()[command]

	if !ok {
		return nil, errr.NewWithHelp("command not found")
	}

	return &cmd, nil
}
