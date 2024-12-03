package version

import (
	"testing"
)

func TestRun(t *testing.T) {
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
}
