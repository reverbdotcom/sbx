[![Vet](https://github.com/reverbdotcom/sbx/actions/workflows/vet.yaml/badge.svg)](https://github.com/reverbdotcom/sbx/actions/workflows/vet.yaml)
[![Release](https://github.com/reverbdotcom/sbx/actions/workflows/release.yml/badge.svg)](https://github.com/reverbdotcom/sbx/actions/workflows/release.yml)

# sbx
Orchestra CLI tool: `sbx up`


## Install / Upgrade


#### golang

```bash
export GOPRIVATE=github.com/reverbdotcom && go install github.com/reverbdotcom/sbx@latest
```


#### brew

```bash
brew tap reverbdotcom/homebrew-reverb
brew update
brew install sbx
```

#### bash

```bash
VERSION=v1.1.3 \
    curl \
        -s\
        -L \
        -o "/tmp/sbx-darwin-arm64.tar.gz" \
        "https://github.com/reverbdotcom/sbx/releases/download/${VERSION}/sbx-darwin-arm64.tar.gz" \
    && tar -xzf /tmp/sbx-darwin-arm64.tar.gz -C /tmp \
    && sudo mv /tmp/sbx /usr/local/bin/sbx
```

## Development

`sbx.go` is the main entry point for the CLI tool.
Every command should be a go package. Commands are
configured in `commands/commands.go`.


`make <command>.run` will build and run the command.
This runs live.

#### Test
`make test`
