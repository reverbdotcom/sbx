package graphiql

import (
	"fmt"

	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/open"
)

const template = "https://graphiql-%s.int.orchestra.rvb.ai/graphql"

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
