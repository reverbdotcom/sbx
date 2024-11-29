package version

import (
	"runtime/debug"
	"testing"
)

func TestRun(t *testing.T) {
	buildInfoFn = func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{
				Path:    "github.com/username/repo",
				Version: "v1.0.0",
			},
		}, true
	}

	t.Run("it returns version", func(t *testing.T) {
		got, err := Run()
		want := "Module: github.com/username/repo, Version: v1.0.0\n"

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("it returns error", func(t *testing.T) {
		buildInfoFn = func() (*debug.BuildInfo, bool) {
			return nil, false
		}

		_, err := Run()

		want := "failed to retrieve build info"
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}
