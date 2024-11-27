package up

import (
	"errors"
	"testing"

	"github.com/reverbdotcom/sbx/cli"
)

func TestRun(t *testing.T) {
	t.Run("it errs on nameFn", func(t *testing.T) {
		nameFn = func() (string, error) {
			return "sandbox-blake-julian-kevin", errors.New("name error")
		}

		_, err := Run()

		if err.Error() != "name error" {
			t.Errorf("got %v", err.Error())
		}
	})

	nameFn = func() (string, error) {
		return "sandbox-blake-julian-kevin", nil
	}

	t.Run("it resets remote if exists, then push", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git ls-remote --heads origin sandbox-blake-julian-kevin",
				Out:     "refs/heads/sandbox-blake-julian-kevin",
				Err:     nil,
			},
			{
				Command: "git push origin --delete sandbox-blake-julian-kevin",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git branch -f sandbox-blake-julian-kevin HEAD",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git push origin sandbox-blake-julian-kevin",
				Out:     "",
				Err:     nil,
			},
		}

		cmdFn = cli.MockCmd(t, wants)

		_, err := Run()

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("it skip reset remote if not exists, then push", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git ls-remote --heads origin sandbox-blake-julian-kevin",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git branch -f sandbox-blake-julian-kevin HEAD",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git push origin sandbox-blake-julian-kevin",
				Out:     "",
				Err:     nil,
			},
		}

		cmdFn = cli.MockCmd(t, wants)

		_, err := Run()

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})

	t.Run("it errs on resets remote", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git ls-remote --heads origin sandbox-blake-julian-kevin",
				Out:     "refs/heads/sandbox-blake-julian-kevin",
				Err:     errors.New("reset error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)

		_, err := Run()

		if err != wants[0].Err {
			t.Errorf("got %v, want %v", err, wants[0].Err)
		}
	})

	t.Run("it errs on push delete", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git ls-remote --heads origin sandbox-blake-julian-kevin",
				Out:     "refs/heads/sandbox-blake-julian-kevin",
				Err:     nil,
			},
			{
				Command: "git push origin --delete sandbox-blake-julian-kevin",
				Out:     "",
				Err:     errors.New("push delete error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)

		_, err := Run()

		if err != wants[1].Err {
			t.Errorf("got %v, want %v", err, wants[1].Err)
		}
	})

	t.Run("it errs on make local", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git ls-remote --heads origin sandbox-blake-julian-kevin",
				Out:     "refs/heads/sandbox-blake-julian-kevin",
				Err:     nil,
			},
			{
				Command: "git push origin --delete sandbox-blake-julian-kevin",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git branch -f sandbox-blake-julian-kevin HEAD",
				Out:     "",
				Err:     errors.New("make local error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)

		_, err := Run()

		if err.Error() != wants[2].Err.Error() {
			t.Errorf("got %v, want %v", err, wants[2].Err)
		}
	})

	t.Run("it errs on push remote", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git ls-remote --heads origin sandbox-blake-julian-kevin",
				Out:     "refs/heads/sandbox-blake-julian-kevin",
				Err:     nil,
			},
			{
				Command: "git push origin --delete sandbox-blake-julian-kevin",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git branch -f sandbox-blake-julian-kevin HEAD",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git push origin sandbox-blake-julian-kevin",
				Out:     "",
				Err:     errors.New("push remote error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)

		_, err := Run()

		if err != wants[3].Err {
			t.Errorf("got %v, want %v", err, wants[3].Err)
		}
	})

	t.Run("it noops on push remote when there are no changes", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git ls-remote --heads origin sandbox-blake-julian-kevin",
				Out:     "refs/heads/sandbox-blake-julian-kevin",
				Err:     nil,
			},
			{
				Command: "git push origin --delete sandbox-blake-julian-kevin",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git branch -f sandbox-blake-julian-kevin HEAD",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git push origin sandbox-blake-julian-kevin",
				Out:     "Everything up-to-date",
				Err:     nil,
			},
		}

		cmdFn = cli.MockCmd(t, wants)

		_,err := Run()

    want := "sandbox-blake-julian-kevin is up to date, make a new commit"
    if err.Error() != want {
      t.Errorf("got %s, want %s", err, want)
    }
	})
}
