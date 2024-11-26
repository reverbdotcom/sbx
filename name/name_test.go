package name

import (
	"testing"
)

func TestName(t *testing.T) {
  t.Run("it generates a sandbox name", func(t *testing.T) {
    branch = func() (string, error) {
      return "nn-sbx-123", nil
    }

    got, err := Run()
    want := "sandbox-foo-bar-baz"

    if err != nil {
      t.Errorf("got %v, want nil", err)
    }

    if got != want {
      t.Errorf("got %v, want %v", got, want)
    }
  })
}
