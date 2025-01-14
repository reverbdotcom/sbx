package web

import (
	"fmt"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/open"
	"github.com/reverbdotcom/sbx/run"
)

const template = "https://%s.int.orchestra.rvb.ai/"
const headlampTemplate = "https://headlamp.preprod.reverb.tools/c/main/deployments?namespace=%s"

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

func OpenHeadlamp() (string, error) {
	name, _ := nameFn()
	url := fmt.Sprintf(headlampTemplate, name)
	err := openURL(url)

	if err != nil {
		return "", err
	}

	return "", nil
}
