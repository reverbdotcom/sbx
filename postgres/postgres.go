package postgres

import (
	"fmt"

	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/open"
	"github.com/reverbdotcom/sbx/run"
)

const template = "https://postgres-%s.int.orchestra.rvb.ai/"

var openURL = open.Open

func Open() (string, error) {
	name, _ := nameFn()
	url := fmt.Sprintf(template, name)
	err := openURL(url)

	if err != nil {
		return "", err
	}

	return "", nil
}

var htmlUrlFn = run.HtmlUrl

func OpenProgress() (string, error) {
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
