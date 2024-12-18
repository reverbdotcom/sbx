package version

import (
	_ "embed"
	"fmt"
	"strings"
)

//go:embed SBX_VERSION
var version string

func Run(_ []string) (string, error) {
	return fmt.Sprintf("Version: %s\n", strings.TrimSpace(version)), nil
}
