package summary

import (
	"fmt"

	"github.com/reverbdotcom/sbx/commit"
	"github.com/reverbdotcom/sbx/dash"
	"github.com/reverbdotcom/sbx/graphiql"
	"github.com/reverbdotcom/sbx/logs"
	"github.com/reverbdotcom/sbx/run"
	"github.com/reverbdotcom/sbx/web"
)

const summary = `»»»
Name:       %s
SHA:        %s

Deploy:     %s
Dash:       %s
Logs:       %s

Host:       %s
Graphiql:   %s
`

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
var dashUrl = dash.Url
var logsUrl = logs.Url
var webUrl = web.Url
var graphiqlUrl = graphiql.Url

func printSummary(name string) (string, error) {
	deployUrl, err := htmlUrl()

	if err != nil {
		return "", err
	}

	sha, err := headSHA()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		summary,
		name,
		sha,
		deployUrl,
		dashUrl(),
		logsUrl(),
		webUrl(),
		graphiqlUrl(),
	), nil
}
