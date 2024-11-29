package logs

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestUrl(t *testing.T) {
	name := "sadbox-foo"
	nameFn = func() (string, error) { return name, nil }
	now = func() time.Time { return time.UnixMilli(1732737340502) }
	sha := "b3a1410a1be2e3d4c4ea9983ef4fa660327683a6"
	headSHA = func() (string, error) { return sha, nil }

	t.Run("it returns url", func(t *testing.T) {
		got := Url()
		want := fmt.Sprintf(template, sha, name, "1732733740502")

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestRun(t *testing.T) {
	name := "sadbox-foo"
	nameFn = func() (string, error) { return name, nil }
	now = func() time.Time { return time.UnixMilli(1732737340502) }
	sha := "b3a1410a1be2e3d4c4ea9983ef4fa660327683a6"
	headSHA = func() (string, error) { return sha, nil }

	t.Run("it opens url", func(t *testing.T) {
		want := fmt.Sprintf(template, sha, name, "1732733740502")

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
		openURL = func(_ string) error {
			return errors.New("open error")
		}

		_, err := Run()

		want := "open error"
		if err.Error() != want {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}
