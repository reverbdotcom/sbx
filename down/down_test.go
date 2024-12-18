package down

import (
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("it runs", func(t *testing.T) {
		nameFn = func() (string, error) { return "sandbox-blake-julian-kevin", nil }
		teardownSandboxFn = func(_ string) error { return nil }

		_, err := Run([]string{})

		if err != nil {
			t.Errorf("got %v", err)
		}
	})

	t.Run("it errs on nameFn", func(t *testing.T) {
		teardownSandboxFn = func(_ string) error { return nil }

		want := errors.New("name error")
		nameFn = func() (string, error) { return "", want }

		_, err := Run([]string{})

		if err.Error() != want.Error() {
			t.Errorf("got %v", err)
		}
	})

	t.Run("it errs on teardownSandboxFn", func(t *testing.T) {
		nameFn = func() (string, error) { return "sandbox-blake-julian-kevin", nil }

		want := errors.New("teardown error")
		teardownSandboxFn = func(_ string) error { return want }

		_, err := Run([]string{})

		if err.Error() != want.Error() {
			t.Errorf("got %v", err)
		}
	})
}
