package web

import (
	"fmt"
	"github.com/reverbdotcom/sbx/open"
)

func Run() (string, error) {
	urlFn := func(name string) string {
		return fmt.Sprintf("https://%s.int.orchestra.rvb.ai/", name)
	}

	err := open.Open(urlFn)

	if err != nil {
		return "", err
	}

	return "", nil
}
