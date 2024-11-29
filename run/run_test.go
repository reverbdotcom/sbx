package run

import (
	"errors"
	"testing"

	"github.com/google/go-github/v67/github"
	"github.com/reverbdotcom/sbx/cli"
)

func TestHtmlUrl(t *testing.T) {
	wantRepo := "foooo"
	wantSHA := "abc123"
	headSHA = func() (string, error) {
		return wantSHA, nil
	}

	t.Run("it retruns html url", func(t *testing.T) {
		findRun = func(repo, sha string) (*github.WorkflowRuns, error) {
			if sha != wantSHA {
				t.Errorf("got %v", sha)
			}

			if repo != wantRepo {
				t.Errorf("got %v", repo)
			}

			wantCount := 1
			wantUrl := "http://some.url"
			return &github.WorkflowRuns{
				TotalCount: &wantCount,
				WorkflowRuns: []*github.WorkflowRun{
					{
						HTMLURL: &wantUrl,
					},
				},
			}, nil
		}

		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     wantRepo,
				Err:     nil,
			},
		}

		cmdFn = cli.MockCmd(t, wants)

		url, err := HtmlUrl()

		if err != nil {
			t.Errorf("got %v", err)
		}

		if url != "http://some.url" {
			t.Errorf("got %v", url)
		}
	})

	t.Run("it errs when current repo errs", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     "",
				Err:     errors.New("repo error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)

		_, err := HtmlUrl()

		if err.Error() != "repo error" {
			t.Errorf("got %v", err)
		}
	})

	t.Run("it errs when find run errs", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     wantRepo,
				Err:     nil,
			},
		}

		cmdFn = cli.MockCmd(t, wants)

		findRun = func(repo, sha string) (*github.WorkflowRuns, error) {
			return nil, errors.New("find run error")
		}

		_, err := HtmlUrl()

		if err.Error() != "find run error" {
			t.Errorf("got %v", err)
		}
	})

	t.Run("it errs when no runs", func(t *testing.T) {
		maxRetries = 0
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     wantRepo,
				Err:     nil,
			},
		}

		cmdFn = cli.MockCmd(t, wants)

		findRun = func(repo, sha string) (*github.WorkflowRuns, error) {
			if sha != wantSHA {
				t.Errorf("got %v", sha)
			}

			if repo != wantRepo {
				t.Errorf("got %v", repo)
			}

			wantCount := 0
			return &github.WorkflowRuns{
				TotalCount: &wantCount,
			}, nil
		}

		_, err := HtmlUrl()

		wantError := "no workflow runs found"
		if err.Error() != wantError {
			t.Errorf("got %v, want %v", err, wantError)
		}
	})
}
