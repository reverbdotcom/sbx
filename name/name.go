package name

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/reverbdotcom/sbx/cli"
)

const maxStep = 2

func Run() (string, error) {
	return name()
}

func name() (string, error) {
	branch, err := branch()

	if err != nil {
		return "", err
	}

	name, err := names(branch)

	if err != nil {
		return "", err
	}

	return prefix(name), nil
}

var branch = _branch

func _branch() (string, error) {
	out, err := cli.Cmd("git", "branch", "--show-current")

	if err != nil {
		return out, err
	}

	return out, nil
}

func names(name string) (string, error) {
	hash1, err1 := hash(name, 0)
	hash2, err2 := hash(name, 1)
	hash3, err3 := hash(name, 2)

	if err1 != nil {
		return "", err1
	}

	if err2 != nil {
		return "", err2
	}

	if err3 != nil {
		return "", err3
	}

	return fmt.Sprintf("%s-%s-%s", hash1, hash2, hash3), nil
}

func hash(name string, step int) (string, error) {
	md5h := md5.Sum([]byte(name))
	offset := step * maxStep
	offsetmd5h := md5h[offset : offset+maxStep]

	words, err := properNames()
	size := len(words)

	if err != nil {
		return "", err
	}

	hex := fmt.Sprintf("%x", offsetmd5h)
	hexInt, err := strconv.ParseInt(hex, 16, 64)

	if err != nil {
		return "", err
	}

	index := int(hexInt) % size

	return strings.ToLower(words[index]), nil
}

func prefix(name string) string {
	return "sandbox-" + name
}

func properNames() ([]string, error) {
	dict, err := dictionary()

	if err != nil {
		return []string{}, err
	}

	names := []string{}
	for _, word := range dict {
		if len(word) > 2 && len(word) < 13 {
			names = append(names, word)
		}
	}

	return names, nil
}

var dictionary = _dictionary

func _dictionary() ([]string, error) {
	file, err := os.Open("/usr/share/dict/propernames")
	if err != nil {
		return []string{}, err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(bytes), "\n"), nil
}
