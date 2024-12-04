package down

import (
	"fmt"

	"github.com/reverbdotcom/sbx/check"
	"github.com/reverbdotcom/sbx/github"
	"github.com/reverbdotcom/sbx/name"
)

var nameFn = name.Name
var ensureOrchestraFn = check.EnsureOrchestra
var teardownSandboxFn = github.TeardownSandbox

func Run() (string, error) {
	err := ensureOrchestraFn()

	if err != nil {
		return "", err
	}

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
