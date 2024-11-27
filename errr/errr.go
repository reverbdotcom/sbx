package errr

import (
	"errors"
	"fmt"

	"github.com/reverbdotcom/sbx/commands"
)

func New(message string) error {
	return errors.New(fmt.Sprintf("🚫 %s\n\n\n%s", message, commands.Help))
}
