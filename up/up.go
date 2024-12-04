package up

import (
	"errors"
	"fmt"
	"strings"

	"github.com/reverbdotcom/sbx/check"
	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/run"
	"github.com/reverbdotcom/sbx/summary"
)

const info = `»»»
Name:       %s
SHA:        %s

Deploy:     %s
Dash:       %s
Logs:       %s

Host:       %s
Graphiql:   %s
`

const noChanges = "up-to-date"

var cmdFn = cli.Cmd
var nameFn = name.Name
var htmlUrlFn = run.HtmlUrl
var summaryFn = summary.Print
var ensureOrchestraFn = check.EnsureOrchestra

func Run() (string, error) {
	err := ensureOrchestraFn()

	if err != nil {
		return "", err
	}

	fmt.Println("deploying...")
	fmt.Println()

	yes, err := isMain()

	if err != nil {
		return "", err
	}

	if yes {
		return "", errors.New("cannot deploy from main branch")
	}

	name, err := nameFn()

	if err != nil {
		return "", err
	}

	out, err := deploy(name, false)

	if err != nil && strings.Contains(err.Error(), noChanges) {
		out, err = deploy(name, true)
	}

	if err != nil {
		return out, err
	}

	err = summaryFn(name)

	if err != nil {
		return "", err
	}

	return out, nil
}

func deploy(name string, noopCommit bool) (string, error) {
	out, err := makeLocal(name, noopCommit)

	if err != nil {
		return out, err
	}

	out, err = pushRemote(name)

	if err != nil {
		return out, err
	}

	return "", nil
}

func makeLocal(name string, noopCommit bool) (string, error) {
	if noopCommit {
		out, err := cmdFn("git", "commit", "--allow-empty", "-m", "'sandbox is up-to-date, noop to trigger'")

		if err != nil {
			return out, err
		}
	}

	out, err := cmdFn("git", "branch", "-f", name, "HEAD")

	if err != nil {
		return out, err
	}

	return out, nil
}

func pushRemote(name string) (string, error) {
	out, err := cmdFn("git", "push", "-f", "origin", name)

	if err != nil {
		return out, err
	}

	if strings.Contains(out, noChanges) {
		return "", errors.New(noChanges)
	}

	return out, nil
}

func isMain() (bool, error) {
	out, err := cmdFn("git", "branch", "--show-current")

	if err != nil {
		return false, err
	}

	yes := strings.TrimSpace(out) == "main"

	return yes, nil
}
