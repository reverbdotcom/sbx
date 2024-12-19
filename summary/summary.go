package summary

import (
	"fmt"
	"github.com/reverbdotcom/sbx/commit"
	"github.com/reverbdotcom/sbx/name"
	"github.com/reverbdotcom/sbx/run"
	"github.com/reverbdotcom/sbx/web"
)

const summary = `»»»
Name:           %s
SHA:            %s
Host:           %s

»»»
Site:           sbx w | sbx web
Graphiql:       sbx g | sbx graphiql

Dash:           sbx d | sbx dash
Logs:           sbx l | sbx logs
Deployment:     sbx p | sbx progress
`

func Run() (string, error) {
	name, err := name.Name()

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
var webUrl = web.Url

func printSummary(name string) (string, error) {
	sha, err := headSHA()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		summary,
		name,
		sha,
		webUrl(),
	), nil
}
