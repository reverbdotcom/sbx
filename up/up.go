package up

import (
  "strings"

	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/name"
)

const noChanges = "up-to-date"

func Run() (string, error) {
  name, err := name.Name()

  if err != nil {
    return "", err
  }

  branchOut, branchErr := cli.Cmd(
    "git",
    "branch",
    "-f",
    name,
    "HEAD",
  )

  if branchErr != nil {
    return branchOut, branchErr
  }

  out, err := cli.Cmd(
    "git",
    "push",
    "origin",
    name,
  )

  if err != nil {
    return out, err
  }

  if strings.Contains(out, noChanges) {
    return "sandbox branch up to date, make a new commit", nil
  }

  return out, nil
}
