package github

import (
	"context"
	"errors"
	"reflect"
	"testing"

	gh "github.com/google/go-github/v67/github"
)

func TestTeardownSandbox(t *testing.T) {
	type args struct {
		owner            string
		repo             string
		workflowFileName string
		event            gh.CreateWorkflowDispatchEventRequest
	}

	t.Run("it dispatches the teardown", func(t *testing.T) {
		got := args{}
		createWorkflowDispatchEventByFileName = func(ctx context.Context, owner, repo, workflowFileName string, event gh.CreateWorkflowDispatchEventRequest) (*gh.Response, error) {
			got = args{owner, repo, workflowFileName, event}

			return nil, nil
		}

		want := args{
			owner:            "reverbdotcom",
			repo:             "flux-system",
			workflowFileName: "orchestra-down.yaml",
			event: gh.CreateWorkflowDispatchEventRequest{
				Ref: "main",
				Inputs: map[string]interface{}{
					"sandbox": "sandbox-blake-julian-kevin",
				},
			},
		}

		err := TeardownSandbox("sandbox-blake-julian-kevin")

		if err != nil {
			t.Errorf("got %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it errs", func(t *testing.T) {
		createWorkflowDispatchEventByFileName = func(_ context.Context, owner, repo, workflowFileName string, event gh.CreateWorkflowDispatchEventRequest) (*gh.Response, error) {
			_ = args{owner, repo, workflowFileName, event}

			return nil, errors.New("some error")
		}

		want := errors.New("some error")
		err := TeardownSandbox("sandbox-blake-julian-kevin")

		if err.Error() != want.Error() {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}

func TestFindSandboxRun(t *testing.T) {
	type args struct {
		owner            string
		repo             string
		workflowFileName string
		opts             *gh.ListWorkflowRunsOptions
	}

	t.Run("it finds the sandbox run", func(t *testing.T) {
		got := args{}
		listWorkflowRunsByFileName = func(ctx context.Context, owner, repo, workflowFileName string, opts *gh.ListWorkflowRunsOptions) (*gh.WorkflowRuns, *gh.Response, error) {
			got = args{owner, repo, workflowFileName, opts}
			runs := &gh.WorkflowRuns{}

			return runs, nil, nil
		}

		want := args{
			owner:            "reverbdotcom",
			repo:             "flux-system",
			workflowFileName: "conductor-on-orchestra.yml",
			opts:             &gh.ListWorkflowRunsOptions{HeadSHA: "sha"},
		}

		runs, err := FindSandboxRun("flux-system", "sha")

		if err != nil {
			t.Errorf("got %v", err)
		}

		if runs == nil {
			t.Errorf("got %v, want not nil", runs)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it errs", func(t *testing.T) {
		listWorkflowRunsByFileName = func(_ context.Context, owner, repo, workflowFileName string, opts *gh.ListWorkflowRunsOptions) (*gh.WorkflowRuns, *gh.Response, error) {
			_ = args{owner, repo, workflowFileName, opts}

			return nil, nil, errors.New("some error")
		}

		want := errors.New("some error")
		_, err := FindSandboxRun("flux-system", "sha")

		if err.Error() != want.Error() {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}
