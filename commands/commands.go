package commands

import (
	"github.com/reverbdotcom/sbx/beta"
	"github.com/reverbdotcom/sbx/console"
	"github.com/reverbdotcom/sbx/dash"
	"github.com/reverbdotcom/sbx/db"
	"github.com/reverbdotcom/sbx/down"
	"github.com/reverbdotcom/sbx/env"
	"github.com/reverbdotcom/sbx/graphiql"
	"github.com/reverbdotcom/sbx/k8s"
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

USAGE
  sbx <command> [flags]

SANDBOX MANAGEMENT
  command     shorthand     description
  up          u             spins up an orchestra sandbox
  down                      tears down an orchestra sandbox
  name        n             shows the sandbox name
  info        i             shows the summary of the sandbox
  beta                      spins up an orchestra sandbox with an isolated data layer

BROWSER TOOLS
  command     shorthand     description
	progress    p             opens deployment progress in a browser
  dash        d             opens the datadog dashboard in a browser
  logs        l             opens the datadog logs in a browser
  web         w             opens the site in a browser
  graphiql    g             opens graphql user interface in a browser
	db                        database console explorer. Use 'sbx db help' for subcommands.

CLI TOOLS
  command     shorthand     description
  k8s                       kubernetes resources explorer. Use 'sbx k8s help' for subcommands
  console     c             opens a core rails console

CONFIGURATION & INFO
  command     shorthand     description
  env         e             shows the configured environment variables the sbx cli
  version     v             shows the version of the sbx cli
  help        h             shows this help message

EXAMPLES
  sbx up
  sbx dash
  sbx k8s help
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
		"db":       db.Run,
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
		"k8s":      k8s.Run,
		"console":  console.Run,
		"c":        console.Run,
	}
}
