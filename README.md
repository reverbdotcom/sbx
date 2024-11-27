[![Vet](https://github.com/reverbdotcom/sbx/actions/workflows/vet.yaml/badge.svg)](https://github.com/reverbdotcom/sbx/actions/workflows/vet.yaml)

# sbx
Orchestra CLI tool: `sbx up`


#### Install / Upgrade

```bash
export GOPRIVATE=github.com/reverbdotcom && go install github.com/reverbdotcom/sbx@latest
```

#### Development

`sbx.go` is the main entry point for the CLI tool.
Every command should be a go package. Commands are
configured in `commands/commands.go`.


`make <command>.run` will build and run the command.
This runs live.

#### Test
`make test`
