package name

import (
	"io"
	"os"
	"strings"
  "crypto/md5"

	"github.com/reverbdotcom/sbx/cli"
)

func Run() (string, error) {
  return name()
}

func name() (string, error) {
  branch, err := branch()

  if err != nil {
    return "", err
  }

  name, err := hash(branch)

  if err != nil {
    return "", err
  }

  return prefix(name), nil
}

var branch = currentBranch // need to mock for testing
func currentBranch() (string, error) {
  out, err := cli.Cmd("git", "branch", "--show-current")

  if err != nil {
    return out, err
  }

  return out, nil
}

func hash(name string) (string, error) {
  md5h := md5.Sum([]byte(name))
  size := len(name)
  words, err := properNames()

  if err != nil {
    return "", err
  }

  idx1 := int(md5h[0]) % size
  idx2 := int(md5h[1]) % size
  idx3 := int(md5h[2]) % size

  names := []string{
    words[idx1],
    words[idx2],
    words[idx3],
  }

  newName := strings.Join(names, "-")

  return strings.ToLower(newName), nil
}

func prefix(name string) (string) {
  return "sandbox-" + name
}

func properNames() ([]string, error) {
	file, err := os.Open("/usr/share/dict/propernames")
	if err != nil {
    return []string{}, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
    return []string{}, err
	}

  // TODO only return > 3 characters
	return strings.Split(string(bytes), "\n"), nil
}
