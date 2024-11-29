package graphiql

import (
	"fmt"

	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/open"
)

const template = "https://graphiql-%s.int.orchestra.rvb.ai/graphql"

var openURL = open.Open

func Run() (string, error) {
	err := openURL(Url())

	if err != nil {
		return "", err
	}

	return "", nil
}

var nameFn = name.Name

func Url() string {
	name, _ := nameFn()
	return fmt.Sprintf(template, name)
}
