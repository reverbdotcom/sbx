package k8s

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestUrl(t *testing.T) {
	name := "sandbox-foo-bar-baz"
	nameFn = func() (string, error) { return name, nil }

	t.Run("it returns url", func(t *testing.T) {
		got := Url()
		want := fmt.Sprintf(template, name)

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestRun(t *testing.T) {
	name := "sandbox-foo-bar-baz"
	nameFn = func() (string, error) { return name, nil }

	t.Run("it opens url", func(t *testing.T) {
		// Mock getArgs to return no subcommand
		getArgs = func() []string { return []string{"sbx", "k8s"} }

		want := fmt.Sprintf(template, name)

		openURL = func(got string) error {
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}

			return nil
		}

		_, err := Run()

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("it returns err", func(t *testing.T) {
		// Mock getArgs to return no subcommand
		getArgs = func() []string { return []string{"sbx", "k8s"} }

		openURL = func(_ string) error {
			return errors.New("open error")
		}

		_, err := Run()

		want := "open error"
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})

	t.Run("it returns help text", func(t *testing.T) {
		// Mock getArgs to return help subcommand
		getArgs = func() []string { return []string{"sbx", "k8s", "help"} }

		result, err := Run()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if result != subcommandHelp {
			t.Errorf("expected help text, got %v", result)
		}
	})

	t.Run("it returns error for unknown subcommand", func(t *testing.T) {
		// Mock getArgs to return unknown subcommand
		getArgs = func() []string { return []string{"sbx", "k8s", "unknown"} }

		_, err := Run()

		if err == nil {
			t.Errorf("expected error for unknown subcommand")
		}

		expectedErr := "unknown subcommand: unknown"
		if !strings.Contains(err.Error(), expectedErr) {
			t.Errorf("expected error to contain '%s', got %v", expectedErr, err.Error())
		}
	})

	t.Run("it calls login subcommand", func(t *testing.T) {
		// Mock getArgs to return login subcommand
		getArgs = func() []string { return []string{"sbx", "k8s", "login"} }

		// We need to ensure this calls the login.Run function
		// Since login.Run requires actual commands, this would fail in tests
		// but we can at least verify the subcommand is registered
		loginCalled := false
		originalLogin := subcommands["login"]
		subcommands["login"] = func() (string, error) {
			loginCalled = true
			return "login called", nil
		}
		defer func() {
			subcommands["login"] = originalLogin
		}()

		result, err := Run()

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		if !loginCalled {
			t.Errorf("expected login subcommand to be called")
		}

		if result != "login called" {
			t.Errorf("expected 'login called', got %v", result)
		}
	})
}
