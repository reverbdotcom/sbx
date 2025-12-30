package beta

import (
	"errors"
	"testing"

	"github.com/reverbdotcom/sbx/cli"
)

func TestRun(t *testing.T) {
	t.Run("it errs on getting git dir", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     "dev/reverb",
				Err:     errors.New("current dir error"),
			},
			{
				Command: "basename dev/reverb",
				Out:     "reverb",
				Err:     errors.New("current dir error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)
		_, err := Run()

		if err == nil || err.Error() != "current dir error" {
			t.Errorf("got %v, want current dir error", err)
		}
	})

	t.Run("it errs on basename", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     "dev/reverb",
				Err:     nil,
			},
			{
				Command: "basename dev/reverb",
				Out:     "reverb",
				Err:     errors.New("basename error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)
		_, err := Run()

		if err == nil || err.Error() != "basename error" {
			t.Errorf("got %v, want basename error", err)
		}
	})

	t.Run("it errs on get branch", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     "dev/reverb",
				Err:     nil,
			},
			{
				Command: "basename dev/reverb",
				Out:     "reverb",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "",
				Err:     errors.New("fetch error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)
		_, err := Run()

		if err == nil || err.Error() != "fetch error" {
			t.Errorf("got %v, want fetch error", err)
		}
	})

	t.Run("it errs when changing directory", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     "dev/reverb",
				Err:     nil,
			},
			{
				Command: "basename dev/reverb",
				Out:     "reverb",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "feature-branch",
				Err:     nil,
			},
		}

		changeDir = func(dir string) error {
			return errors.New("cd error")
		}

		cmdFn = cli.MockCmd(t, wants)
		_, err := Run()

		if err == nil || err.Error() != "failed to change dir to reverb: cd error" {
			t.Errorf("got %v, want failed to change dir to reverb: cd error", err)
		}
	})

	changeDir = func(dir string) error {
		return nil
	}

	t.Run("it errs on fetch", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     "dev/reverb",
				Err:     nil,
			},
			{
				Command: "basename dev/reverb",
				Out:     "reverb",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "feature-branch",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     errors.New("fetch error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)
		_, err := Run()

		if err == nil || err.Error() != "failed to fetch repo reverb: fetch error" {
			t.Errorf("got %v, want failed to fetch repo reverb: fetch error", err)
		}
	})

	t.Run("it errs on checkout branch", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     "dev/reverb",
				Err:     nil,
			},
			{
				Command: "basename dev/reverb",
				Out:     "reverb",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "feature-branch",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "",
				Err:     errors.New("checkout error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)
		_, err := Run()

		if err == nil || err.Error() != "failed to checkout branch feature-branch in repo reverb: checkout error" {
			t.Errorf("got %v, want failed to checkout branch feature-branch in repo reverb: checkout error", err)
		}
	})

	t.Run("it errs on creating a new branch", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     "dev/reverb",
				Err:     nil,
			},
			{
				Command: "basename dev/reverb",
				Out:     "reverb",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "feature-branch",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "error: pathspec 'feature-branch' did not match any file(s) known to git",
				Err:     errors.New("checkout error"),
			},
			{
				Command: "git checkout -b feature-branch origin/main-sandbox",
				Out:     "Already on 'feature-branch'",
				Err:     errors.New("create branch error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)
		_, err := Run()

		if err == nil || err.Error() != "failed to checkout branch feature-branch in repo reverb: create branch error" {
			t.Errorf("got %v, want failed to checkout branch feature-branch in repo reverb: create branch error", err)
		}
	})

	t.Run("it errs on sbx up", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     "dev/reverb",
				Err:     nil,
			},
			{
				Command: "basename dev/reverb",
				Out:     "reverb",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "feature-branch",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "Already on 'feature-branch'",
				Err:     nil,
			},
		}
		upFn = func() (string, error) {
			return "", errors.New("deploy error")
		}

		cmdFn = cli.MockCmd(t, wants)
		_, err := Run()

		if err == nil || err.Error() != "deploy failed for reverb: deploy error" {
			t.Errorf("got %v, want deploy failed for reverb: deploy error", err)
		}
	})

	upFn = func() (string, error) {
		return "", nil
	}

	sleep = func(_ int) {}

	t.Run("it successfully deploys four repos", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     "dev/reverb",
				Err:     nil,
			},
			{
				Command: "basename dev/reverb",
				Out:     "reverb",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "feature-branch",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "Already on 'feature-branch'",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "Already on 'feature-branch'",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "Already on 'feature-branch'",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "Already on 'feature-branch'",
				Err:     nil,
			},
		}

		cmdFn = cli.MockCmd(t, wants)
		out, err := Run()

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		if out != "beta sandbox deploy triggered; monitor each deploy for completion" {
			t.Errorf("got %v, want beta sandbox deploy triggered; monitor each deploy for completion", out)
		}
	})

	t.Run("it successfully deploys five repos", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse --show-toplevel",
				Out:     "dev/notreverb",
				Err:     nil,
			},
			{
				Command: "basename dev/notreverb",
				Out:     "notreverb",
				Err:     nil,
			},
			{
				Command: "git branch --show-current",
				Out:     "feature-branch",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "Already on 'feature-branch'",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "Already on 'feature-branch'",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "Already on 'feature-branch'",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "Already on 'feature-branch'",
				Err:     nil,
			},
			{
				Command: "git fetch",
				Out:     "",
				Err:     nil,
			},
			{
				Command: "git checkout feature-branch",
				Out:     "Already on 'feature-branch'",
				Err:     nil,
			},
		}

		cmdFn = cli.MockCmd(t, wants)
		out, err := Run()

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		if out != "beta sandbox deploy triggered; monitor each deploy for completion" {
			t.Errorf("got %v, want beta sandbox deploy triggered; monitor each deploy for completion", out)
		}
	})
}
