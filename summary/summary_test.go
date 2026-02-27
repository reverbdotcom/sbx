package summary

import (
	"errors"
	"testing"
)

func TestPrint(t *testing.T) {
	nameFn = func() (string, error) { return "sandbox-lorraine-francisco-jonathan", nil }

	t.Run("it prints the summary", func(t *testing.T) {
		name := "sandbox-blake-julian-kevin"

		headSHA = func() (string, error) { return "head.sha", nil }
		got, err := printSummary(name)

		if err != nil {
			t.Errorf("got %v, want nil", err)
		}

		want := `»»»
Name:           sandbox-blake-julian-kevin
SHA:            head.sha
Host:           https://sandbox-lorraine-francisco-jonathan.int.orchestra.rvb.ai/

»»»
Site:           sbx w | sbx web
Console:        sbx c | sbx console
Graphiql:       sbx g | sbx graphiql
Redis:          sbx db r | sbx db redis
Postgres:       sbx db p | sbx db postgres
Elasticsearch:  sbx db e | sbx db elasticsearch

Dash:           sbx d | sbx dash
Logs:           sbx l | sbx logs
Deployment:     sbx p | sbx progress
`

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("it errs on headSHA", func(t *testing.T) {
		name := "sandbox-blake-julian-kevin"

		headSHA = func() (string, error) { return "", errors.New("headSHA error") }

		_, err := printSummary(name)

		if err.Error() != "headSHA error" {
			t.Errorf("got %v, want headSHA error", err.Error())
		}
	})
}
