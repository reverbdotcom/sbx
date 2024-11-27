package parser

import (
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

		want := "command required"
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
		}

		if cmdfn != nil {
			t.Errorf("got up.Run, want nil")
		}
	})

	t.Run("it errs when command is not found", func(t *testing.T) {
		args := []string{"sbx", "does-not-exist"}
		cmdfn, err := Parse(args)

		if err == nil {
			t.Errorf("got nil, want error")
		}

		want := "command not found"
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
		}

		if cmdfn != nil {
			t.Errorf("got up.Run, want nil")
		}
	})
}
