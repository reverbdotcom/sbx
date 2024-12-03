package version

import (
  _ "embed"
	"fmt"
	"strings"
)


//go:embed SBX_CHECKSUM
var checksum string

//go:embed SBX_VERSION
var version string


func Run() (string, error) {
	return fmt.Sprintf("Version: %s, Checksum: %s", strings.TrimSpace(version), checksum), nil
}
