package version

import (
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("it returns version", func(t *testing.T) {
    version = "v1.0.0"

		got, err := Run()
    want := "Version: v1.0.0\n"

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})
}
