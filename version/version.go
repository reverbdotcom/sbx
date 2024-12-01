package version

import (
	"errors"
	"fmt"

	"runtime/debug"
)

var buildInfoFn = debug.ReadBuildInfo

func Run() (string, error) {
	info, ok := buildInfoFn()
	if !ok {
		return "", errors.New("failed to retrieve build info")
	}

	return fmt.Sprintf("Module: %s, Version: %s\n", info.Main.Path, info.Main.Version), nil
}
