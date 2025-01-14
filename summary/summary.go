package summary

import (
	"fmt"
	"github.com/reverbdotcom/sbx/commit"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/run"
)

const template = "https://%s.int.orchestra.rvb.ai/"
const summary = `»»»
Name:           %s
SHA:            %s
Host:           %s

»»»
Site:           sbx w | sbx web
Graphiql:       sbx g | sbx graphiql

Headlamp:       sbx h | sbx headlamp
Dash:           sbx d | sbx dash
Logs:           sbx l | sbx logs
Deployment:     sbx p | sbx progress
`

var nameFn = name.Name

func Run() (string, error) {
	name, err := nameFn()

	if err != nil {
		return "", err
	}

	return printSummary(name)
}

func Print(name string) error {
	summary, err := printSummary(name)

	if err != nil {
		return err
	}

	fmt.Println(summary)

	return nil
}

var htmlUrl = run.HtmlUrl
var headSHA = commit.HeadSHA

func printSummary(name string) (string, error) {
	sha, err := headSHA()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		summary,
		name,
		sha,
		url(),
	), nil
}

func url() string {
	name, _ := nameFn()
	return fmt.Sprintf(template, name)
}
