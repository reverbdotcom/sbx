package up

import (
  "strings"

	"github.com/reverbdotcom/sbx/cli"
	"github.com/reverbdotcom/sbx/name"
)

const noChanges = "up-to-date"

func Run() (string, error) {
  exists, err := remoteExists()

  if err != nil {
    return "", err
  }

  if exists {
    out, err := resetRemote()

    if err != nil {
      return out, err
    }
  }

  out, err := makeLocal()

  if err != nil {
    return out, err
  }

  out, err = pushRemote()

  if err != nil {
    return out, err
  }

  return out, nil
}

func remoteExists() (bool, error) {
  name, err := name.Name()
  out, err := cli.Cmd("git", "ls-remote", "--heads", "origin", name)

  if err != nil {
    return false, err
  }

  if strings.Contains(out, name) {
    return true, nil
  }

  return false, nil
}

func resetRemote() (string, error) {
  name, err := name.Name()
  out, err := cli.Cmd("git", "push", "origin", "--delete", name)

  if err != nil {
    return out, err
  }

  return out, nil
}

func makeLocal() (string, error) {
  name, err := name.Name()
  out, err := cli.Cmd("git", "branch", "-f", name, "HEAD")


  if err != nil {
    return out, err
  }

  return out, nil
}

func pushRemote() (string, error) {
  name, err := name.Name()
  out, err := cli.Cmd("git", "push", "origin", name)

  if err != nil {
    return out, err
  }

  return out, nil
}
