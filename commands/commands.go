package commands

import (
	"github.com/reverbdotcom/sbx/dash"
	"github.com/reverbdotcom/sbx/down"
	"github.com/reverbdotcom/sbx/graphiql"
	"github.com/reverbdotcom/sbx/help"
	"github.com/reverbdotcom/sbx/logs"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/summary"
	"github.com/reverbdotcom/sbx/up"
	"github.com/reverbdotcom/sbx/version"
	"github.com/reverbdotcom/sbx/web"
)

type RunFn func(cmdArgs []string) (string, error)

func Commands() map[string]RunFn {
	return map[string]RunFn{
		"up":       up.Run,
		"u":        up.Run,
		"help":     help.Run,
		"h":        help.Run,
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
	}
}
