package web

import (
	"fmt"

	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/open"
	"github.com/reverbdotcom/sbx/run"
)

const template = "https://%s.int.orchestra.rvb.ai/"

var openURL = open.Open

func Open(_ []string) (string, error) {
	err := openURL(Url())

	if err != nil {
		return "", err
	}

	return "", nil
}

var htmlUrlFn = run.HtmlUrl

func OpenProgress(_ []string) (string, error) {
	htmlUrl, err := htmlUrlFn()

	if err != nil {
		return "", err
	}

	if err := openURL(htmlUrl); err != nil {
		return "", err
	}

	return "", nil
}

var nameFn = name.Name

func Url() string {
	name, _ := nameFn()
	return fmt.Sprintf(template, name)
}
