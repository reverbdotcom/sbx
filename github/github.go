package github

import (
	"context"
	"errors"
	"os"

	"github.com/google/go-github/v67/github"
	"github.com/reverbdotcom/sbx/check"
)

const owner = "reverbdotcom"

func Client() (*github.Client, error) {
	if !check.HasGithubToken() {
		return nil, errors.New("Please set the GITHUB_TOKEN environment variable.")
	}

	return github.NewClient(nil).WithAuthToken(os.Getenv("GITHUB_TOKEN")), nil
}

var listWorkflowRunsByFileName = _listWorkflowRunsByFileName

func _listWorkflowRunsByFileName(ctx context.Context, owner, repo, workflowFileName string, opts *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error) {
	client, err := Client()

	if err != nil {
		return nil, nil, err
	}

	return client.Actions.ListWorkflowRunsByFileName(ctx, owner, repo, workflowFileName, opts)
}

func FindSandboxRun(repo, sha string) (*github.WorkflowRuns, error) {
	workflow := "conductor-on-orchestra.yml"
	ctx := context.Background()

	opts := &github.ListWorkflowRunsOptions{HeadSHA: sha}
	runs, _, err := listWorkflowRunsByFileName(ctx, owner, repo, workflow, opts)

	return runs, err
}

var createWorkflowDispatchEventByFileName = _createWorkflowDispatchEventByFileName

func _createWorkflowDispatchEventByFileName(ctx context.Context, owner, repo, workflowFileName string, event github.CreateWorkflowDispatchEventRequest) (*github.Response, error) {
	client, err := Client()

	if err != nil {
		return nil, err
	}

	return client.Actions.CreateWorkflowDispatchEventByFileName(ctx, owner, repo, workflowFileName, event)
}

func TeardownSandbox(name string) error {
	workflow := "orchestra-down.yaml"
	repo := "flux-system"
	ctx := context.Background()

	opts := github.CreateWorkflowDispatchEventRequest{
		Ref:    "main",
		Inputs: map[string]interface{}{"sandbox": name},
	}

	_, err := createWorkflowDispatchEventByFileName(ctx, owner, repo, workflow, opts)

	if err != nil {
		return err
	}

	return err
}
