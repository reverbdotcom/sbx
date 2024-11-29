package commands

import (
	"github.com/reverbdotcom/sbx/dash"
	"github.com/reverbdotcom/sbx/down"
	"github.com/reverbdotcom/sbx/graphiql"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/up"
	"github.com/reverbdotcom/sbx/version"
	"github.com/reverbdotcom/sbx/web"
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
      version
      down ( not yet supported )

DESCRIPTION

  command     shorthand     description

  up          u             spin up an orchestra sandbox.
  help                      display this help message.
  name        n             display the name of the current branch.
  web                       open the site in a browser.
  dash                      open the dashboard in a browser.
  graphiql                  open graphql user interface in a browser.
  graphql                   alias for 'graphiql'.
  down                      tear down the sandbox.
  version     v             show the version of the sbx cli.

USAGE:

  sbx up
  sbx name
`

func help() (string, error) {
	return Help, nil
}

func Commands() map[string]RunFn {
	return map[string]RunFn{
		"up":       up.Run,
		"u":        up.Run,
		"help":     help,
		"name":     name.Run,
		"n":        name.Run,
		"web":      web.Run,
		"dash":     dash.Run,
		"graphiql": graphiql.Run,
		"graphql":  graphiql.Run,
		"down":     down.Run,
		"version":  version.Run,
		"v":        version.Run,
	}
}
