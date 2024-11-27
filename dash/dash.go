package dash

import (
  "fmt"
	"github.com/reverbdotcom/sbx/cli"
  "github.com/reverbdotcom/sbx/name"
)

func Run() (string, error) {
  name, err := name.Name()

  if err != nil {
    return "", err
  }

  url := fmt.Sprintf("https://app.datadoghq.com/dashboard/9rm-fjs-8tx/orchestra?tpl_var_sandbox[]=%s", name)

  out, err := cli.Cmd("open", url)

  if err != nil {
    return out, err
  }

  return "", nil
}

