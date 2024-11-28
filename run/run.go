package run

import (
	"context"
	"errors"
	"github.com/google/go-github/v67/github"
	"os"
	"path/filepath"
	"strings"

	"github.com/reverbdotcom/sbx/cli"
)

const owner = "reverbdotcom"
const workflow = "conductor-on-orchestra.yml"

func HtmlUrl() (string, error) {
	sha, err := headSHA()

	if err != nil {
		return "", err
	}

	run, err := currentRun(sha)

	if err != nil {
		return "", err
	}

	return *run.HTMLURL, nil
}

func currentRun(commitSHA string) (*github.WorkflowRun, error) {
	ctx := context.Background()
	client := client()

	repo, err := currentRepo()

	if err != nil {
		return nil, err
	}

	opts := &github.ListWorkflowRunsOptions{HeadSHA: commitSHA}
	runs, _, err := client.Actions.ListWorkflowRunsByFileName(ctx, owner, repo, workflow, opts)

	if err != nil {
		return nil, err
	}

	if *runs.TotalCount == 0 {
		return nil, errors.New("no workflow runs found")
	}

	return nil, nil
}

func client() *github.Client {
	return github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN"))
}

func currentRepo() (string, error) {
	out, err := cli.Cmd("git", "rev-parse", "--show-toplevel")

	if err != nil {
		return out, err
	}

	path := strings.TrimSpace(out)
	repo := filepath.Base(path)

	return repo, nil
}

func headSHA() (string, error) {
	out, err := cli.Cmd("git", "rev-parse", "HEAD")

	if err != nil {
		return out, err
	}

	path := strings.TrimSpace(out)
	sha := filepath.Base(path)

	return sha, nil
}
