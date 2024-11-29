package check

import (
	"github.com/reverbdotcom/sbx/cli"
	"os"
	"strings"
)

const yml = "/.github/workflows/conductor-on-orchestra.yml"

func HasGithubToken() bool {
	token := os.Getenv("GITHUB_TOKEN")

	if token != "" {
		return true
	}

	return false
}

func OnOrchestra() (bool, error) {
	dir, err := gitDir()

	if err != nil {
		return false, err
	}

	path := dir + yml

	if fileExists(path) {
		return true, nil
	}

	return false, nil
}

func gitDir() (string, error) {
	out, err := cli.Cmd("git", "rev-parse", "--show-toplevel")

	if err != nil {
		return out, err
	}

	dir := strings.TrimSpace(out)

	return dir, nil
}

var fileExists = _fileExists

func _fileExists(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	}

	return false
}
