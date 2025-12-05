package check

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/reverbdotcom/sbx/cli"
)

const yml = "/.github/workflows/conductor-on-orchestra.yml"

func HasGithubToken() bool {
	token := os.Getenv("GITHUB_TOKEN")

	if token != "" {
		return true
	}

	return false
}

func EnsureOrchestra() error {
	has, err := OnOrchestra()
	if err != nil {
		if strings.Contains(err.Error(), "not a git repository") {
			return fmt.Errorf("This project is not on Orchestra: %v", err)
		}
		return err
	}

	if !has {
		return errors.New("This project is not on Orchestra.")
	}

	return nil
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

var gitDir = _gitDir

func _gitDir() (string, error) {
	out, err := cli.Cmd("git", "rev-parse", "--show-toplevel")

	if err != nil {
		return out, fmt.Errorf("%v: %v", out, err)
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
