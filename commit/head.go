package commit

import (
	"path/filepath"
	"strings"

	"github.com/reverbdotcom/sbx/cli"
)

var cmdFn = cli.Cmd

func HeadSHA() (string, error) {
	out, err := cmdFn("git", "rev-parse", "HEAD")

	if err != nil {
		return out, err
	}

	path := strings.TrimSpace(out)
	sha := filepath.Base(path)

	return sha, nil
}
