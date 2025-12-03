package k8s

import (
	"fmt"
	"os"

	"github.com/reverbdotcom/sbx/containers"
	"github.com/reverbdotcom/sbx/crons"
	"github.com/reverbdotcom/sbx/deployments"
	"github.com/reverbdotcom/sbx/ingresses"
	"github.com/reverbdotcom/sbx/jobs"
	"github.com/reverbdotcom/sbx/login"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/open"
	"github.com/reverbdotcom/sbx/pods"
	"github.com/reverbdotcom/sbx/processes"
	"github.com/reverbdotcom/sbx/services"
	"github.com/reverbdotcom/sbx/ssh"
)

const template = "https://app.datadoghq.com/kubernetes?query=kube_namespace:%s%%20kube_cluster_name:preprod-v6"

var openURL = open.Open

const subcommandHelp = `USAGE:
  sbx k8s [subcommand]

SUBCOMMANDS:
  login         authenticates with AWS and switches to preprod kubernetes context
  ssh           drops into a kubernetes pod container shell
  pods          opens the kubernetes pod view
  deployments   opens the kubernetes deployment view
  jobs          opens the kubernetes job view
  crons         opens the kubernetes cron job view
  services      opens the kubernetes service view
  ingresses     opens the kubernetes ingress view
  processes     opens the datadog process view
  containers    opens the datadog containers view

If no subcommand is provided, opens the kubernetes cluster view.
`

type SubcommandFn func() (string, error)

var subcommands = map[string]SubcommandFn{
	"login":       login.Run,
	"ssh":         ssh.Run,
	"pods":        pods.Run,
	"deployments": deployments.Run,
	"jobs":        jobs.Run,
	"crons":       crons.Run,
	"services":    services.Run,
	"ingresses":   ingresses.Run,
	"processes":   processes.Run,
	"containers":  containers.Run,
}

var getArgs = func() []string {
	return os.Args
}

func Run() (string, error) {
	args := getArgs()

	// Check if there's a subcommand
	if len(args) > 2 {
		subcommand := args[2]

		if subcommand == "help" || subcommand == "-h" || subcommand == "--help" {
			return subcommandHelp, nil
		}

		if fn, ok := subcommands[subcommand]; ok {
			return fn()
		}

		return "", fmt.Errorf("unknown subcommand: %s\n\n%s", subcommand, subcommandHelp)
	}

	// Default behavior: open cluster view
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
