package help

import (
	"fmt"

	"github.com/reverbdotcom/sbx/up"
)

type HelpFn func() (string, error)

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

DESCRIPTION

  command     shorthand     description

  help        h             show the help message. use "sbx help <cmd>" for more info about that command.
  up          u             spin up an orchestra sandbox.
  down                      tear down an orchestra sandbox.
  name        n             show the sandbox name.
  dash        d             open the dashboard in a browser.
  logs        l             open the logs in a browser.
  web         w             open the site in a browser.
  graphiql    g             open graphql user interface in a browser.
  version     v             show the version of the sbx cli.
  info        i             show the summary of the sandbox.
  progress    p             open deployment progress in a browser.

USAGE:

  sbx up
  sbx name
`

func Run(cmdArgs []string) (string, error) {
	var totalArgs = len(cmdArgs)
	if totalArgs == 0 {
		return Help, nil
	}

	if totalArgs > 1 {
		return "", fmt.Errorf("too many arguments")
	}

	var cmd = cmdArgs[0]
	var helpFn = helps()[cmd]
	if helpFn == nil {
		return fmt.Sprintf("no help info available for command %s", cmd), nil
	}

	return helpFn()
}

func helps() map[string]HelpFn {
	return map[string]HelpFn{
		"up": up.Help,
	}
}
