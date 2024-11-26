package up

import (
	"github.com/reverbdotcom/sbx/cli"
)

func Run() (string, error) {
	return cli.Cmd("git", "status")
}
