package commands

import (
	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/up"
)

const Help = `
Usage: sbx <help|up|name>

Examples:
  sbx up
  sbx name
`

func help() (string, error) {
	return Help, nil
}

func Commands() map[string]cli.CmdFn {
	return map[string]cli.CmdFn{
		"up":   up.Run,
		"name": name.Run,
		"help": help,
	}
}
