package commands

import (
	"github.com/reverbdotcom/sbx/dash"
	"github.com/reverbdotcom/sbx/down"
	"github.com/reverbdotcom/sbx/web"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/up"
	"github.com/reverbdotcom/sbx/graphiql"
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
      graphiql
      graphql
      down ( not yet supported )

DESCRIPTION

  up          spin up an orchestra sandbox.
  help        display this help message.
  name        display the name of the current branch.
  web         open the site in a browser.
  dash        open the dashboard in a browser.
  graphiql    open graphql user interface in a browser.
  graphql     alias for 'graphiql'.
  down        tear down the sandbox.

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
		"help": help,
		"name": name.Run,
		"web": web.Run,
		"dash": dash.Run,
		"graphiql": graphiql.Run,
		"graphql": graphiql.Run,
		"down": down.Run,
	}
}
