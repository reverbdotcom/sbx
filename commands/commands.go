package commands

import (
	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/up"
	"github.com/reverbdotcom/sbx/name"
)

func Commands() map[string]cli.CmdFn {
	return map[string]cli.CmdFn{
		"up": up.Run,
		"name": name.Run,
	}
}
