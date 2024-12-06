package parser

import (
	"errors"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("it returns up command", func(t *testing.T) {
		ensureOrchestra = func() error { return nil }
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
		ensureOrchestra = func() error { return nil }
		args := []string{"sbx"}
		cmdfn, err := Parse(args)

		if err == nil {
			t.Errorf("got nil, want error")
		}

		want := "command required"
		if strings.Contains(want, err.Error()) {
			t.Errorf("got %v, want %v", err, want)
		}

		if cmdfn != nil {
			t.Errorf("got up.Run, want nil")
		}
	})

	t.Run("it errs when command is not found", func(t *testing.T) {
		ensureOrchestra = func() error { return nil }
		args := []string{"sbx", "does-not-exist"}
		cmdfn, err := Parse(args)

		if err == nil {
			t.Errorf("got nil, want error")
		}

		want := "command not found"
		if strings.Contains(want, err.Error()) {
			t.Errorf("got %v, want %v", err, want)
		}

		if cmdfn != nil {
			t.Errorf("got up.Run, want nil")
		}
	})

	t.Run("it errs on non orchestra for non-general commands", func(t *testing.T) {
		ensureOrchestra = func() error { return errors.New("This project is not on Orchestra.") }

		args := []string{"sbx", "up"}
		_, err := Parse(args)

		want := "This project is not on Orchestra."
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})

	t.Run("it does not err on non orchestra for genereal command", func(t *testing.T) {
    general := []string{
      "help",
      "h",
      "version",
      "v",
    }

    for _, cmd := range general {
      ensureOrchestra = func() error { return errors.New("This project is not on Orchestra.") }

      args := []string{"sbx", cmd}
      _, err := Parse(args)

      if err != nil {
        t.Errorf("got %v, want nil, for cmd %s", err, cmd)
      }
    }
	})
}
