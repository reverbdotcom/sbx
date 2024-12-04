package commands

import (
	"github.com/reverbdotcom/sbx/dash"
	"github.com/reverbdotcom/sbx/down"
	"github.com/reverbdotcom/sbx/graphiql"
	"github.com/reverbdotcom/sbx/logs"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/summary"
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
      down
      name
      dash
      logs
      web
      graphiql
      version
      info

DESCRIPTION

  command     shorthand     description

  help        h             show the help message.
  up          u             spin up an orchestra sandbox.
  down                      tear down an orchestra sandbox.
  name        n             show the sandbox name.
  dash        d             open the dashboard in a browser.
  logs        l             open the logs in a browser.
  web         w             open the site in a browser.
  graphiql    g             open graphql user interface in a browser.
  version     v             show the version of the sbx cli.
  info        i             show the summary of the sandbox.

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
		"g":        graphiql.Run,
		"down":     down.Run,
		"version":  version.Run,
		"v":        version.Run,
		"logs":     logs.Run,
		"l":        logs.Run,
		"info":     summary.Run,
		"i":        summary.Run,
	}
}
