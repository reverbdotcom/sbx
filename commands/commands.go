package commands

import (
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/up"
)

type RunFn func() (string, error)

const Help = `
Usage: sbx <help|up|name>

Examples:
  sbx up
  sbx name
`

func help() (string, error) {
	return Help, nil
}

func Commands() map[string]RunFn {
	return map[string]RunFn{
		"up":   up.Run,
		"name": name.Run,
		"help": help,
	}
}
