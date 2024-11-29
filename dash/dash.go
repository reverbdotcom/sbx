package dash

import (
	"fmt"

	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/open"
)

const template = "https://app.datadoghq.com/dashboard/9rm-fjs-8tx/orchestra?tpl_var_sandbox[]=%s"

func Run() (string, error) {
	err := open.Open(Url())

	if err != nil {
		return "", err
	}

	return "", nil
}

func Url() string {
	name, _ := name.Name()
	return fmt.Sprintf(template, name)
}
