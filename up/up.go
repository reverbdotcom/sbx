package up

import (
	"errors"
	"fmt"
	"strings"

	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/name"
)

const info = `🚀 deploying orchestra sandbox: [ %s ]

»» commands:
»» sbx help	: display this help message
»» sbx dash	: open the dashboard in a browser
»» sbx name	: display the sandbox name
`

const noChanges = "up-to-date"

var cmdFn = cli.Cmd
var nameFn = name.Name

func Run() (string, error) {
	name, err := nameFn()

	if err != nil {
		return "", err
	}

	out, err := makeLocal(name)

	if err != nil {
		return out, err
	}

	out, err = pushRemote(name, true)

	if err != nil {
		return out, err
	}

	fmt.Printf(info, name)

	return "", nil
}

func makeLocal(name string) (string, error) {
	out, err := cmdFn("git", "branch", "-f", name, "HEAD")

	if err != nil {
		return out, err
	}

	return out, nil
}

func pushRemote(name string, noop bool) (string, error) {
	out, err := cmdFn("git", "push", "origin", name)

	if err != nil {
		return out, err
	}

  if noop && strings.Contains(out, noChanges) {
    out, err := noopCommit()

    if err != nil {
      return out, err
    }

    out, err = pushRemote(name, false)

    if err != nil {
      return out, err
    }
  } else if strings.Contains(out, noChanges) {
    return out, errors.New(name + " is up to date, make a new commit")
  }


	return out, nil
}

func noopCommit() (string, error) {
	out, err := cmdFn("git", "commit", "--allow-empty", "-m", "'sandbox is up-to-date, noop commit to trigger deploy'")

	if err != nil {
		return out, err
	}

	return out, nil
}
