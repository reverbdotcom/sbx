package run

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/go-github/v67/github"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/commit"
	"github.com/reverbdotcom/sbx/debug"
)

const owner = "reverbdotcom"
const workflow = "conductor-on-orchestra.yml"
const notFound = "no workflow runs found"

func HtmlUrl() (string, error) {
	sha, err := commit.HeadSHA()

	if err != nil {
		return "", err
	}

	run := &github.WorkflowRun{}
	maxRetries := 5

	for i := 0; i <= maxRetries; i++ {
		backoff := time.Duration(i*2) * time.Second
		run, err = currentRun(sha)

		if err != nil && err.Error() != notFound {
			return "", err
		}

		if run != nil {
			break
		}

		if i < maxRetries {
			if debug.On() {
				fmt.Printf("Waiting for run... sha: %v, attempt: %v, backoff: %v\n", sha, i+1, backoff)
			}
			time.Sleep(backoff)
		}
	}

	if run == nil {
		return "", errors.New(notFound)
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
		return nil, errors.New(notFound)
	}

	return runs.WorkflowRuns[0], nil
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
