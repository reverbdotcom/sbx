package parser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("it returns up command", func(t *testing.T) {
		args := []string{"sbx", "up"}
		cmdfn, err := Parse(args)

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		if cmdfn == nil {
			t.Errorf("got nil, want up.Run")
		}
	})

	t.Run("it errs when no command is provided", func(t *testing.T) {
		args := []string{"sbx"}
		cmdfn, err := Parse(args)

		if err == nil {
			t.Errorf("got nil, want error")
		}

		if strings.Contains("command required", err.Error()) {
			t.Errorf("got %v, want 'command required'", err)
		}

		if cmdfn != nil {
			t.Errorf("got up.Run, want nil")
		}
	})

	t.Run("it errs when command is not found", func(t *testing.T) {
		args := []string{"sbx", "down"}
		cmdfn, err := Parse(args)

		if err == nil {
			t.Errorf("got nil, want error")
		}

		if strings.Contains("command not found", err.Error()) {
			t.Errorf("got %v, want 'command not found'", err)
		}

		if cmdfn != nil {
			t.Errorf("got up.Run, want nil")
		}
	})
}
