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

var allowlist = map[string]string{
	"DURATION": "how long sandboxes live for",
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

var Getenv = _getenv
func _getenv(key string) string {
	return readenv()[key]
}

func Duration() (string, error) {
	dur := Getenv("DURATION")
	if dur == "" {
		return dur, nil
	}

	parsed, err := time.ParseDuration(dur)
	if err != nil {
		return "", fmt.Errorf("invalid duration: %s", dur)
	}

	if parsed > durationTooLong {
		warning(fmt.Sprintf("%s is a long duration! this is okay on occassion, but consider lowering it", dur))
	}

	return dur, nil
}

func readenv() map[string]string {
	filepath := fmt.Sprintf("%s/.sbx", os.Getenv("HOME"))

	var err error
	env, err := godotenv.Read(filepath)
	if err != nil {
		errr.Warning(fmt.Sprintf("could not open %s, defaulting to blank env", filepath))
		env = make(map[string]string)
	}

	return env
}
