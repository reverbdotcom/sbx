package parser

import (
	"errors"
	"fmt"

	"github.com/reverbdotcom/sbx/check"
	"github.com/reverbdotcom/sbx/commands"
	"github.com/reverbdotcom/sbx/help"
	"golang.org/x/exp/slices"
)

func Parse(args []string) (*commands.RunFn, []string, error) {
	cmd, cmdArgs, err := command(args)

	if err != nil {
		return nil, nil, err
	}

	cmdFn, err := cmdfn(*cmd)
	if err != nil {
		return nil, nil, err
	}

	return cmdFn, cmdArgs, nil
}

func command(args []string) (command *string, commandArgs []string, err error) {
	if len(args) < 2 {
		return nil, nil, errr("command required")
	}

	cmd := args[1]

	cmdArgs := []string{}
	if len(args) > 2 {
		cmdArgs = args[2:]
	}

	return &cmd, cmdArgs, nil
}

var ensureOrchestra = check.EnsureOrchestra

func cmdfn(command string) (*commands.RunFn, error) {
	general := []string{
		"help",
		"h",
		"version",
		"v",
	}

	if !slices.Contains(general, command) {
		if err := ensureOrchestra(); err != nil {
			return nil, err
		}
	}

	cmd, ok := commands.Commands()[command]

	if !ok {
		return nil, errr("command not found")
	}

	return &cmd, nil
}

func errr(message string) error {
	return errors.New(fmt.Sprintf("%s\n\n\n%s", message, help.Help))
}
