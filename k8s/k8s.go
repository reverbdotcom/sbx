package k8s

import (
	"fmt"

	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/open"
)

const template = "https://app.datadoghq.com/kubernetes?query=kube_namespace:%s%%20kube_cluster_name:preprod-v6"

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
