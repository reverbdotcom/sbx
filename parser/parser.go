package parser

import (
	"fmt"
	"errors"
	"os"

	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/up"
)

const help = `
Usage: sbx <up|down|name>
`

func Parse() (*cli.CmdFn, error) {
	cmd, err := command()

	if err != nil {
		return nil, err
	}

	return cmdfn(*cmd)
}

func command() (command *string, err error) {
	if len(os.Args) < 2 {
		return nil, errr("command required")
	}

	cmd := os.Args[1]

	return &cmd, nil
}

func cmdfn(command string) (*cli.CmdFn, error) {
	cmd, ok := map[string]cli.CmdFn{
		"up": up.Run,
	}[command]

	if !ok {
		return nil, errr("command not found")
	}

	return &cmd, nil
}

func errr(message string) error {
  return errors.New(fmt.Sprintf("%s\n%s", message, help))
}
