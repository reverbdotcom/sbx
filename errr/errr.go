package errr

import (
	"errors"
	"fmt"
)

func New(message string) error {
	return errors.New(fmt.Sprintf("🚫 %s", message))
}
