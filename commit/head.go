package commit

import (
	"path/filepath"
	"strings"

	"github.com/reverbdotcom/sbx/cli"
)

func HeadSHA() (string, error) {
	out, err := cli.Cmd("git", "rev-parse", "HEAD")

	if err != nil {
		return out, err
	}

	path := strings.TrimSpace(out)
	sha := filepath.Base(path)

	return sha, nil
}
