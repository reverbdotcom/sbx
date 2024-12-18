package help

import (
	"strings"
	"testing"

	"github.com/reverbdotcom/sbx/up"
)

func TestRun(t *testing.T) {
	t.Run("it returns generic help text with no args", func(t *testing.T) {
		output, err := Run([]string{})

		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		if output != Help {
			t.Errorf("did not return generic help text, got %v", output)
		}
	})

	t.Run("it fails when there's too many args", func(t *testing.T) {
		output, err := Run([]string{"up", "foo"})
		if err == nil {
			t.Errorf("help should have errored, but got %v", output)
		}
	})

	t.Run("it returns help text for command", func(t *testing.T) {
		output, err := Run([]string{"up"})
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		text, _ := up.Help()
		if output != text {
			t.Errorf("did not return help text for up, got %v", output)
		}
	})

	t.Run("it reports when command has no help info", func(t *testing.T) {
		output, err := Run([]string{"foo"})
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		if !strings.Contains(output, "no help info") {
			t.Errorf("expected to report no help info for command, but got %v", output)
		}
	})
}
