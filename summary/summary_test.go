package summary

import (
	"errors"
	"testing"
)

func TestPrint(t *testing.T) {
	t.Run("it prints the summary", func(t *testing.T) {
		name := "sandbox-blake-julian-kevin"

		htmlUrl = func() (string, error) { return "deploy.url", nil }
		headSHA = func() (string, error) { return "head.sha", nil }
		dashUrl = func() string { return "dash.url" }
		logsUrl = func() string { return "logs.url" }
		webUrl = func() string { return "web.url" }
		graphiqlUrl = func() string { return "graphiql.url" }

		got, err := printSummary(name)

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		want := `»»»
Name:       sandbox-blake-julian-kevin
SHA:        head.sha

Deploy:     deploy.url
Dash:       dash.url
Logs:       logs.url

Host:       web.url
Graphiql:   graphiql.url
`

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it errs on htmlUrl", func(t *testing.T) {
		name := "sandbox-blake-julian-kevin"

		htmlUrl = func() (string, error) { return "", errors.New("htmlUrl error") }

		_, err := printSummary(name)

		if err.Error() != "htmlUrl error" {
			t.Errorf("got %v, want htmlUrl error", err.Error())
		}
	})

	t.Run("it errs on headSHA", func(t *testing.T) {
		name := "sandbox-blake-julian-kevin"

		htmlUrl = func() (string, error) { return "deploy.url", nil }
		headSHA = func() (string, error) { return "", errors.New("headSHA error") }

		_, err := printSummary(name)

		if err.Error() != "headSHA error" {
			t.Errorf("got %v, want headSHA error", err.Error())
		}
	})
}
