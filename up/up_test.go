package up

import (
	"errors"
	"testing"

	"github.com/reverbdotcom/sbx/cli"
)

func TestRun(t *testing.T) {
	summaryFn = func(_ string) error { return nil }
	htmlUrlFn = func() (string, error) {
		return "some.url", nil
	}

	t.Run("it errs on main", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git branch --show-current",
				Out:     "main",
				Err:     nil,
			},
		}

		cmdFn = cli.MockCmd(t, wants)
		_, err := Run()

		if err.Error() != "cannot deploy from main branch" {
			t.Errorf("got %v", err.Error())
		}
	})

	t.Run("it errs on nameFn", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git branch --show-current",
				Out:     "sandbox-blake-julian-kevin",
				Err:     nil,
			},
		}

		cmdFn = cli.MockCmd(t, wants)
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

	t.Run("it push a new remote", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git branch --show-current",
				Out:     "nn-my-branch",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "nn-my-branch",
				Err:     nil,
			},
			{
				Command: "git branch -f sandbox-blake-julian-kevin HEAD",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git push -f origin sandbox-blake-julian-kevin",
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

	t.Run("it errs on make local", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git branch --show-current",
				Out:     "nn-my-branch",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "nn-my-branch",
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
				Command: "git branch --show-current",
				Out:     "nn-my-branch",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "nn-my-branch",
				Err:     nil,
			},
			{
				Command: "git branch -f sandbox-blake-julian-kevin HEAD",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git push -f origin sandbox-blake-julian-kevin",
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

	t.Run("it creates a noop commit on push remote when remote is up to date", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git branch --show-current",
				Out:     "nn-my-branch",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "nn-my-branch",
				Err:     nil,
			},
			{
				Command: "git branch -f sandbox-blake-julian-kevin HEAD",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git push -f origin sandbox-blake-julian-kevin",
				Out:     "Everything up-to-date",
				Err:     nil,
			},
			{
				Command: "git commit --allow-empty -m 'sandbox is up-to-date, noop to trigger'",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "nn-my-branch",
				Err:     nil,
			},
			{
				Command: "git branch -f sandbox-blake-julian-kevin HEAD",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git push -f origin sandbox-blake-julian-kevin",
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

	t.Run("it pushes to remote with new changes", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git branch --show-current",
				Out:     "nn-my-branch",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "nn-my-branch",
				Err:     nil,
			},
			{
				Command: "git branch -f sandbox-blake-julian-kevin HEAD",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git push -f origin sandbox-blake-julian-kevin",
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

	t.Run("does not create a local branch is already on the sandbox branch", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git branch --show-current",
				Out:     "sandbox-blake-julian-kevin",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "sandbox-blake-julian-kevin",
				Err:     nil,
			},
			{
				Command: "git push -f origin sandbox-blake-julian-kevin",
				Out:     "",
				Err:     nil,
			},
		}

		nameFn = func() (string, error) {
			return "sandbox-blake-julian-kevin", nil
		}

		cmdFn = cli.MockCmd(t, wants)

		_, err := Run()

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}
	})
}
