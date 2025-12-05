package check

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestHasGithubToken(t *testing.T) {
	original := os.Getenv("GITHUB_TOKEN")

	t.Run("when GITHUB_TOKEN is set", func(t *testing.T) {
		os.Setenv("GITHUB_TOKEN", "foo")
		if !HasGithubToken() {
			t.Errorf("expected HasGithubToken() to be true")
		}
		os.Setenv("GITHUB_TOKEN", original)
	})

	t.Run("when GITHUB_TOKEN is not set", func(t *testing.T) {
		os.Setenv("GITHUB_TOKEN", "")
		if HasGithubToken() {
			t.Errorf("expected HasGithubToken() to be false")
		}
		os.Setenv("GITHUB_TOKEN", original)
	})
}

func TestOnOrchestra(t *testing.T) {
	t.Run("when conductor-on-orchestra.yml exists", func(t *testing.T) {
		want := "sbx/.github/workflows/conductor-on-orchestra.yml"

		fileExists = func(path string) bool {
			if !strings.Contains(path, want) {
				t.Errorf("got %v, want %v", path, want)
			}

			return true
		}

		on, err := OnOrchestra()

		if !on {
			t.Errorf("expected to be true")
		}

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("when conductor-on-orchestra.yml does not exist", func(t *testing.T) {
		want := "sbx/.github/workflows/conductor-on-orchestra.yml"

		fileExists = func(path string) bool {
			if !strings.Contains(path, want) {
				t.Errorf("got %v, want %v", path, want)
			}

			return false
		}

		on, err := OnOrchestra()

		if on {
			t.Errorf("expected to be false")
		}

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("when run outside of a git repository", func(t *testing.T) {

		gitDir = func() (string, error) {
			return "exit status 128", errors.New("fatal: not a git repository (or any of the parent directories): .git")
		}

		err := EnsureOrchestra()

		if err != nil {
			if !strings.Contains(err.Error(), "This project is not on Orchestra") {
				t.Errorf("got %v, want 'This project is not on Orchestra'", err)
			}
		}
	})
}
