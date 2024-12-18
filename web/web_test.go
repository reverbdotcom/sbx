package web

import (
	"errors"
	"fmt"
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

func TestOpen(t *testing.T) {
	name := "sandbox-foo-bar-baz"
	nameFn = func() (string, error) { return name, nil }

	t.Run("it opens url", func(t *testing.T) {
		want := fmt.Sprintf(template, name)

		openURL = func(got string) error {
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}

			return nil
		}

		_, err := Open([]string{})

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("it returns err", func(t *testing.T) {
		openURL = func(_ string) error {
			return errors.New("open error")
		}

		_, err := Open([]string{})

		want := "open error"
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}

func TestOpenProgress(t *testing.T) {
	name := "sandbox-foo-bar-baz"
	nameFn = func() (string, error) { return name, nil }

	t.Run("it opens url", func(t *testing.T) {
		want := "progress.html"

		openURL = func(got string) error {
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}

			return nil
		}

		htmlUrlFn = func() (string, error) {
			return "progress.html", nil
		}

		_, err := OpenProgress([]string{})

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("errs on htmlUrlFn", func(t *testing.T) {
		htmlUrlFn = func() (string, error) {
			return "", errors.New("htmlUrlFn error")
		}

		_, err := OpenProgress([]string{})

		want := "htmlUrlFn error"
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})

	t.Run("errs on openURL", func(t *testing.T) {
		openURL = func(_ string) error {
			return errors.New("openURL error")
		}

		htmlUrlFn = func() (string, error) {
			return "progress.html", nil
		}

		_, err := OpenProgress([]string{})

		want := "openURL error"
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}
