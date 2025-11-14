package env

import (
	"testing"
)

func TestVerify(t *testing.T) {
	t.Run("it errs if duration is invalid", func(t *testing.T) {
		Getenv = func(key string) string {
			if key == DURATION {
				return "invalid"
			}

			return ""
		}

		err := Verify()
		if err == nil {
			t.Errorf("got nil, want error")
		}
	})

	t.Run("it warns when duration is too long", func(t *testing.T) {
		Getenv = func(key string) string {
			if key == DURATION {
				return "10h"
			}

			return ""
		}

		warned := false
		warning = func(message string) {
			warned = true
		}

		_ = Verify()
		if !warned {
			t.Error("got no warning, want warning about duration being too long")
		}
	})
}
