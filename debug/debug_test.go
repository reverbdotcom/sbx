package debug

import (
	"os"
	"testing"
)

func TestOn(t *testing.T) {
	original := os.Getenv("DEBUG")

	t.Run("when DEBUG is true", func(t *testing.T) {
		os.Setenv("DEBUG", "true")
		if !On() {
			t.Errorf("expected On() to be true")
		}
		os.Setenv("DEBUG", original)
	})

	t.Run("when DEBUG is false", func(t *testing.T) {
		os.Setenv("DEBUG", "false")
		if On() {
			t.Errorf("expected On() to be false")
		}
		os.Setenv("DEBUG", original)
	})
}
