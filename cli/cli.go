package cli

import (
	"bytes"
	"os/exec"

	"strings"
	"testing"
)

type CmdFn func(name string, args ...string) (string, error)

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

type MockCall struct {
	Command string
	Out     string
	Err     error
}

func MockCmd(t *testing.T, wants []MockCall) CmdFn {
	step := 0

	return func(name string, args ...string) (string, error) {
		if step >= len(wants) {
			t.Errorf("more calls than mocks, at %v, mock size %v", step, len(wants))
		}

		c := wants[step]

		command := name + " " + strings.Join(args, " ")

		if command != c.Command {
			t.Errorf("got %v, want %v", command, c.Command)
		}

		step++

		return c.Out, c.Err
	}
}
