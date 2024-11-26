package up

import (
	"github.com/reverbdotcom/sbx/cli"
)

func Run() error {
	cli.Cmd("git", "statusdasd")
	return nil
}
