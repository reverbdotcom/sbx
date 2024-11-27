package up

import (
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

  out, err := deploy(name)

  if err != nil {
    return out, err
  }

	fmt.Printf(info, name)

  return out, nil
}

func deploy(name string) (string, error) {
	out, err := makeLocal(name)

	if err != nil {
		return out, err
	}

	out, err = pushRemote(name)

	if err != nil {
		return out, err
	}

	return "", nil
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
    out, err := deployNoop(name)

		if err != nil {
			return out, err
		}
	}

	return out, nil
}

func deployNoop(name string) (string, error) {
	out, err := cmdFn("git", "commit", "--allow-empty", "-m", "'sandbox is up-to-date, noop to trigger'")

  if err != nil {
    return out, err
  }

  out, err = deploy(name)

	if err != nil {
		return out, err
	}

	return out, nil
}
