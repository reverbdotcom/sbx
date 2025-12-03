package commands

import (
	"github.com/reverbdotcom/sbx/beta"
	"github.com/reverbdotcom/sbx/dash"
	"github.com/reverbdotcom/sbx/down"
	"github.com/reverbdotcom/sbx/env"
	"github.com/reverbdotcom/sbx/graphiql"
	"github.com/reverbdotcom/sbx/logs"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/pods"
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
      progress
      env
      beta
      pods

DESCRIPTION

  command     shorthand     description

  help                      shows the help message.
  up          u             spins up an orchestra sandbox.
  down                      tears down an orchestra sandbox.
  name        n             shows the sandbox name.
  dash        d             opens the dashboard in a browser.
  logs        l             opens the logs in a browser.
  web         w             opens the site in a browser.
  graphiql    g             opens graphql user interface in a browser.
  version     v             shows the version of the sbx cli.
  info        i             shows the summary of the sandbox.
  progress    p             opens deployment progress in a browser.
  env         e             shows the configured environment variables for sbx.
  beta                      spins up an orchestra sandbox with an isolated data layer
  pods                      opens the kubernetes pod view in a browser.

USAGE:
  sbx <command> [flags]
  sbx up
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
		"web":      web.Open,
		"progress": web.OpenProgress,
		"p":        web.OpenProgress,
		"w":        web.Open,
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
		"env":      env.Run,
		"e":        env.Run,
		"beta":     beta.Run,
		"pods":     pods.Run,
	}
}
