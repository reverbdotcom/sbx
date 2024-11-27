package dash

import (
	"fmt"
	"github.com/reverbdotcom/sbx/open"
)

func Run() (string, error) {
	urlFn := func(name string) string {
		return fmt.Sprintf("https://app.datadoghq.com/dashboard/9rm-fjs-8tx/orchestra?tpl_var_sandbox[]=%s", name)
	}

	err := open.Open(urlFn)

	if err != nil {
		return "", err
	}

	return "", nil
}
