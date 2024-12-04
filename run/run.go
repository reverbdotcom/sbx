package run

import (
	"errors"
	"github.com/google/go-github/v67/github"
	"path/filepath"
	"strings"

	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/commit"
	gh "github.com/reverbdotcom/sbx/github"
	"github.com/reverbdotcom/sbx/retries"
)

const notFound = "no workflow runs found"

var findRun = gh.FindSandboxRun
var cmdFn = cli.Cmd
var headSHA = commit.HeadSHA
var maxRetries = 5

func HtmlUrl() (string, error) {
	sha, err := headSHA()

	if err != nil {
		return "", err
	}

	run := &github.WorkflowRun{}
	err = retries.Backoff(maxRetries, 2, func() (bool, error) {
		run, err = currentRun(sha)

		if err != nil && err.Error() != notFound {
			return false, err
		}

		if run != nil {
			return true, nil
		}

		return false, nil
	})

	if err != nil && err.Error() == retries.ErrBackoffExhausted {
		return "", errors.New(notFound)
	}

	if err != nil {
		return "", err
	}

	if run == nil {
		return "", errors.New(notFound)
	}

	return *run.HTMLURL, nil
}

func currentRun(commitSHA string) (*github.WorkflowRun, error) {
	repo, err := currentRepo()

	if err != nil {
		return nil, err
	}

	runs, err := findRun(repo, commitSHA)

	if err != nil {
		return nil, err
	}

	if *runs.TotalCount == 0 {
		return nil, errors.New(notFound)
	}

	return runs.WorkflowRuns[0], nil
}

func currentRepo() (string, error) {
	out, err := cmdFn("git", "rev-parse", "--show-toplevel")

	if err != nil {
		return out, err
	}

	path := strings.TrimSpace(out)
	repo := filepath.Base(path)

	return repo, nil
}
