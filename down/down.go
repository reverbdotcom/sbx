package down

import (
	"fmt"

	"github.com/reverbdotcom/sbx/github"
	"github.com/reverbdotcom/sbx/name"
)

var nameFn = name.Name
var teardownSandboxFn = github.TeardownSandbox

func Run(_ []string) (string, error) {
	name, err := nameFn()

	if err != nil {
		return "", err
	}

	err = teardownSandboxFn(name)

	if err != nil {
		return "", err
	}

	fmt.Println("teardown sent")

	return "", nil
}
