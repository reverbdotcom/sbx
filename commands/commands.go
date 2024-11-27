package commands

import (
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/up"
	"github.com/reverbdotcom/sbx/dash"
)

type RunFn func() (string, error)

const Help = `NAME
  sbx - orchestra cli

COMMANDS

  sbx up
      help
      name
      web
      dash
      gateway

DESCRIPTION

  up      spin up an orchestra sandbox.
  help    display this help message.
  name    display the name of the current branch.
  web     open the web interface.

USAGE:

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
		"dash": dash.Run,
		"help": help,
	}
}
