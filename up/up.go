package up

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/go-github/v67/github"
	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/errr"
	gh "github.com/reverbdotcom/sbx/github"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/run"
	"github.com/reverbdotcom/sbx/summary"
	"github.com/reverbdotcom/sbx/version"
)

const noChanges = "up-to-date"

var cmdFn = cli.Cmd
var nameFn = name.Name
var htmlUrlFn = run.HtmlUrl
var summaryFn = summary.Print
var latestReleaseFn = latestRelease

func Run() (string, error) {
	upgrade()

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

	yes, err := onSandbox(name)

	if err != nil {
		return "", err
	}

	if yes { // noop
		return "", nil
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

func onSandbox(name string) (bool, error) {
	out, err := cmdFn("git", "branch", "--show-current")

	if err != nil {
		return false, err
	}

	yes := strings.TrimSpace(out) == name

	return yes, nil
}

func latestRelease() (*github.RepositoryRelease, error) {
	client, err := gh.Client()

	if err != nil {
		return nil, err
	}

	release, _, err := client.Repositories.GetLatestRelease(context.Background(), "reverbdotcom", "sbx")

	return release, err
}

func newVersionAvailable() bool {
	release, err := latestReleaseFn()

	if err != nil {
		return false
	}

	latest := release.GetTagName()
	current := version.Get()

	return latest != current
}

func upgrade() {
	if !newVersionAvailable() {
		return
	}

	out, err := cmdFn("brew", "update")

	if err != nil {
		errr.Warning(fmt.Sprintf("sbx update: %s %s", out, err))
		return
	}

	out, err = cmdFn("arch", "-arm64", "brew", "upgrade", "sbx")

	if err != nil {
		errr.Warning(fmt.Sprintf("sbx upgrade: %s %s", out, err))
	}
}
