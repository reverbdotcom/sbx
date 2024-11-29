package commit

import (
	"errors"
	"testing"

	"github.com/reverbdotcom/sbx/cli"
)

func TestHeadSHA(t *testing.T) {
	t.Run("it gets the head sha", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse HEAD",
				Out:     "123456",
				Err:     nil,
			},
		}

		cmdFn = cli.MockCmd(t, wants)
		sha, err := HeadSHA()

		if err != nil {
			t.Errorf("got %v", err)
		}

		if sha != "123456" {
			t.Errorf("got %v", sha)
		}
	})

	t.Run("it errs on cmdFn", func(t *testing.T) {
		wants := []cli.MockCall{
			{
				Command: "git rev-parse HEAD",
				Out:     "123456",
				Err:     errors.New("cmd error"),
			},
		}

		cmdFn = cli.MockCmd(t, wants)
		_, err := HeadSHA()

		if err.Error() != "cmd error" {
			t.Errorf("got %v", err)
		}
	})
}
