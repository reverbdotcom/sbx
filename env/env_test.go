package env

import (
	"testing"
)

func TestDuration(t *testing.T) {
	t.Run("it returns duration if set", func(t *testing.T) {
		Getenv = func(key string) string {
			return "5h"
		}

		got, err := Duration()
		want := "5h"

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it errs if duration is invalid", func(t *testing.T) {
		Getenv = func(key string) string {
			return "invalid"
		}

		_, err := Duration()
		if err == nil {
			t.Errorf("got nil, want error")
		}
	})

	t.Run("warns when duration is too long", func(t *testing.T) {
		Getenv = func(key string) string {
			return "10h"
		}

		warned := false
		warning = func(message string) {
			warned = true
		}

		_, _ = Duration()
		if !warned {
			t.Error("got no warning, want warning about duration being too long")
		}
	})
}
