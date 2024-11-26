package name

import (
	"github.com/reverbdotcom/sbx/cli"
)

func Run() (string, error) {
  return currentBranch()
}

func currentBranch() (string, error) {
  out, err := cli.Cmd("git", "branch", "--show-current")

  if err != nil {
    return out, err
  }

  return out, nil
}
