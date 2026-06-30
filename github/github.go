package github

import (
	"context"
	"errors"
	"fmt"
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

var getContents = _getContents

func _getContents(ctx context.Context, owner, repo, path string, opts *github.RepositoryContentGetOptions) (*github.RepositoryContent, []*github.RepositoryContent, *github.Response, error) {
	client, err := Client()

	if err != nil {
		return nil, nil, nil, err
	}

	return client.Repositories.GetContents(ctx, owner, repo, path, opts)
}

// GetFileContents fetches and decodes a single file from a reverbdotcom repo.
func GetFileContents(repo, path string) (string, error) {
	ctx := context.Background()
	opts := &github.RepositoryContentGetOptions{Ref: "main"}

	fileContent, _, _, err := getContents(ctx, owner, repo, path, opts)
	if err != nil {
		return "", err
	}

	if fileContent == nil {
		return "", fmt.Errorf("%s/%s/%s is not a file", owner, repo, path)
	}

	return fileContent.GetContent()
}
