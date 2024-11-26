package cli

import (
	"bytes"
	"os/exec"
)

type CmdFn = func() (string, error)

func Cmd(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)

	buffer := bytes.Buffer{}
	cmd.Stdout = &buffer
	cmd.Stderr = &buffer

	err := cmd.Run()
	if err != nil {
		return buffer.String(), err
	}

	return buffer.String(), nil
}
