package graphiql

import (
	"fmt"

	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/open"
)

func Run() (string, error) {
	urlFn := func(name string) string {
		return fmt.Sprintf("https://graphiql-%s.int.orchestra.rvb.ai/graphql", name)
	}

	err := open.Open(urlFn)

	if err != nil {
		return "", err
	}

	return "", nil
}

func Url() string {
	name, _ := name.Name()
	return fmt.Sprintf("https://graphiql-%s.int.orchestra.rvb.ai/graphql", name)
}
