package parser

import (
	"errors"
	"fmt"

	"github.com/reverbdotcom/sbx/check"
	"github.com/reverbdotcom/sbx/commands"
	"golang.org/x/exp/slices"
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
		return nil, errr("command required")
	}

	cmd := args[1]

	return &cmd, nil
}

var ensureOrchestra = check.EnsureOrchestra

func cmdfn(command string) (*commands.RunFn, error) {
	general := []string{
		"help",
		"h",
		"version",
		"v",
		"env",
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
	return errors.New(fmt.Sprintf("%s\n\n\n%s", message, commands.Help))
}
