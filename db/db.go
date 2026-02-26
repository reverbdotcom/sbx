package db

import (
	"fmt"
	"os"

	"github.com/reverbdotcom/sbx/elasticsearch"
	"github.com/reverbdotcom/sbx/postgres"
	"github.com/reverbdotcom/sbx/redis"
)

const subcommandHelp = `USAGE:
  sbx db [subcommand]

SUBCOMMANDS:
  redis / r             opens the redis UI in a browser
  postgres / p / psql   opens the postgres UI in a browser
  elasticsearch / e / cerebro   opens the elasticsearch UI (cerebro) in a browser

If no subcommand is provided, shows this help.
`

type SubcommandFn func() (string, error)

var subcommands = map[string]SubcommandFn{
	"redis":         redis.Open,
	"r":             redis.Open,
	"postgres":      postgres.Open,
	"p":             postgres.Open,
	"psql":          postgres.Open,
	"elasticsearch": elasticsearch.Open,
	"e":             elasticsearch.Open,
	"cerebro":       elasticsearch.Open,
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

	// Default behavior: show help
	return subcommandHelp, nil
}
