package redis

import (
	"errors"
	"fmt"
	"testing"
)

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

		_, err := Open()

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("it returns err", func(t *testing.T) {
		openURL = func(_ string) error {
			return errors.New("open error")
		}

		_, err := Open()

		want := "open error"
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}
