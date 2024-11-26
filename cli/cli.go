package cli

import (
	"fmt"
	"os"
	"os/exec"
)

type CmdFn = func() error

func Cmd(name string, arg ...string) {
	// Define the git command
	cmd := exec.Command(name, arg...)

	// Set the command's standard input/output to the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}
