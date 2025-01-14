package env

import (
	"os"
	"fmt"
	"time"
	"github.com/joho/godotenv"
	"github.com/reverbdotcom/sbx/errr"
)

const format = "%-20s%-20s%-120s\n"
const durationTooLong = 8 * time.Hour

const DURATION = "DURATION"
var allowlist = map[string]string{
	DURATION: "how long sandboxes live for",
}

var warning = errr.Warning

func Run() (string, error) {
	output := fmt.Sprintf(format, "KEY", "VALUE", "DESCRIPTION")
	output += fmt.Sprint("\n")

	for key, description := range allowlist {
		val := Getenv(key)
		if val == "" {
			val = "n/a"
		}

		output += fmt.Sprintf(format, key, val, description)
	}

	return output, nil
}

func Verify() (error) {
	fmtErr := func(msg string) error {
		return fmt.Errorf("invalid env, please fix the following issue: %s", msg)
	}

	dur := Getenv(DURATION)
	if dur != "" {
		parsed, err := time.ParseDuration(dur)
		if err != nil {
			return fmtErr(fmt.Sprintf("%s is an invalid duration", dur))
		}

		if parsed > durationTooLong {
			warning(fmt.Sprintf("%s is a long duration! this is okay on occassion, but consider lowering it", dur))
		}
	}

	return nil
}

var Getenv = _getenv
func _getenv(key string) string {
	return readenv()[key]
}

func readenv() map[string]string {
	filepath := fmt.Sprintf("%s/.sbx", os.Getenv("HOME"))

	var err error
	env, err := godotenv.Read(filepath)
	if err != nil {
		env = make(map[string]string)
	}

	return env
}
