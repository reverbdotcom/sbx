package name

import (
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("it generates a sandbox name", func(t *testing.T) {
		branch = func() (string, error) {
			return "nn-sbx-1234", nil
		}

		got, err := Run()
		want := "sandbox-blake-julian-kevin"

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestProperNames(t *testing.T) {
	t.Run("it should be longer than 2 and less than 13", func(t *testing.T) {
		words, err := properNames()

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		for _, word := range words {
			if len(word) < 3 {
				t.Errorf("got %v, want longer than 2", word)
			}

			if len(word) > 12 {
				t.Errorf("got %v, want less than 13", word)
			}
		}
	})
}
