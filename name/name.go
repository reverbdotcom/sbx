package name

import (
	"github.com/reverbdotcom/sbx/cli"
)

func Name() (*string, error) {
	return nil, nil
}

func branchName() (*string, error) {

	cli.Cmd("git", "branch", "--show-current")

	return nil, nil
}
