package open

import (
	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/name"
)

type urlFn func(name string) string

func Open(fn urlFn) error {
	name, err := name.Name()

	if err != nil {
		return err
	}

	url := fn(name)
	_, err = cli.Cmd("open", url)

	if err != nil {
		return err
	}

	return nil
}
