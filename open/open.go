package open

import (
	"github.com/reverbdotcom/sbx/cli"
)

func Open(url string) error {
	_, err := cli.Cmd("open", url)

	if err != nil {
		return err
	}

	return nil
}
