package up

import (
	"fmt"
	"errors"
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

	exists, err := remoteExists(name)

	if err != nil {
		return "", err
	}

	if exists {
		out, err := resetRemote(name)

		if err != nil {
			return out, err
		}
	}

	out, err := makeLocal(name)

	if err != nil {
		return out, err
	}

	out, err = pushRemote(name)

	if err != nil {
		return out, err
	}

	fmt.Printf(info, name)

	return "", nil
}

func remoteExists(name string) (bool, error) {
	out, err := cmdFn("git", "ls-remote", "--heads", "origin", name)

	if err != nil {
		return false, err
	}

	if strings.Contains(out, name) {
		return true, nil
	}

	return false, nil
}

func resetRemote(name string) (string, error) {
	out, err := cmdFn("git", "push", "origin", "--delete", name)

	if err != nil {
		return out, err
	}

	return out, nil
}

func makeLocal(name string) (string, error) {
	out, err := cmdFn("git", "branch", "-f", name, "HEAD")

	if err != nil {
		return out, err
	}

	return out, nil
}

func pushRemote(name string) (string, error) {
	out, err := cmdFn("git", "push", "origin", name)

	if err != nil {
		return out, err
	}

  if strings.Contains(out, noChanges) {
    return out, errors.New(name + " is up to date, make a new commit")
  }

	return out, nil
}
