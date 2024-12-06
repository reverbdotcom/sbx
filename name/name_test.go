package name

import (
	"testing"
)

func TestRun(t *testing.T) {
	dictionary = func() ([]string, error) {
		return []string{"blake", "julian", "kevin"}, nil
	}

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
	dictionary = func() ([]string, error) {
		return []string{"blake", "julian", "kevin", "a", "super-long-name-that-does-not-fit"}, nil
	}

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

func TestName(t *testing.T) {
	dictionary = func() ([]string, error) {
		return []string{"blake", "julian", "kevin"}, nil
	}

	t.Run("it generates a sandbox name", func(t *testing.T) {
		branch = func() (string, error) {
			return "nn-sbx-1234", nil
		}

		got, err := Name()
		want := "sandbox-blake-julian-kevin"

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("skips name hash if branch starts with sandbox", func(t *testing.T) {
		branch = func() (string, error) {
			return "sandbox-already-named", nil
		}

		got, err := Name()
		want := "sandbox-already-named"

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it errs if name does not start with sandbox", func(t *testing.T) {
		name = func() (string, error) {
			return "not-sandbox", nil
		}

		_, err := Name()

		if err == nil {
			t.Errorf("got nil, want error")
		}

		want := "name does not start with sandbox-"
		if want != err.Error() {
			t.Errorf("got %v, want %v", err, want)
		}
	})
}
