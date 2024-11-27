package up

import (
	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/name"
)

func Run() (string, error) {
	return cli.Cmd("git", "status")
}
