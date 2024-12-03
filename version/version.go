package version

import (
  _ "embed"
	"fmt"
	"strings"
)


//go:embed SBX_VERSION
var version string


func Run() (string, error) {
	return fmt.Sprintf("Version: %s\n", strings.TrimSpace(version)), nil
}
