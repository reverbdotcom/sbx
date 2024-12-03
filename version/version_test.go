package version

import (
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("it returns version", func(t *testing.T) {
    version = "v1.0.0"
    checksum = "1234567890"

		got, err := Run()
    want := "Version: v1.0.0, Checksum: 1234567890"

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})
}
