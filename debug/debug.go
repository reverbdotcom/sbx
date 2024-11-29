package debug

import (
	"os"
)

func On() bool {
	return os.Getenv("DEBUG") == "true"
}
