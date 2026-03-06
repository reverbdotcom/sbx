package version

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed SBX_VERSION
var version string

func Get() string {
	return strings.TrimSpace(version)
}

func Run() (string, error) {
	return fmt.Sprintf("Version: %s\n", Get()), nil
}
