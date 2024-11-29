package commands

import (
	"github.com/reverbdotcom/sbx/dash"
	"github.com/reverbdotcom/sbx/down"
	"github.com/reverbdotcom/sbx/graphiql"
	"github.com/reverbdotcom/sbx/log"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/up"
	"github.com/reverbdotcom/sbx/version"
	"github.com/reverbdotcom/sbx/web"
)

type RunFn func() (string, error)

const Help = `NAME
  sbx - orchestra cli

COMMANDS

  sbx help
      up
      name
      dash
      logs
      web
      graphiql
      graphql
      version

DESCRIPTION

  command     shorthand     description

  help        h             display this help message.
  up          u             spin up an orchestra sandbox.
  name        n             display the name of the current branch.
  dash        d             open the dashboard in a browser.
  logs        l             open the logs in a browser.
  web         w             open the site in a browser.
  graphiql    g             open graphql user interface in a browser.
  graphql                   alias for 'graphiql'.
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
		"h":        help,
		"name":     name.Run,
		"n":        name.Run,
		"web":      web.Run,
		"w":        web.Run,
		"dash":     dash.Run,
		"d":        dash.Run,
		"graphiql": graphiql.Run,
		"graphql":  graphiql.Run,
		"g":        graphiql.Run,
		"down":     down.Run,
		"version":  version.Run,
		"v":        version.Run,
		"logs":     log.Run,
		"l":        log.Run,
	}
}
