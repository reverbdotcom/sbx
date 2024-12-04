[![Vet](https://github.com/reverbdotcom/sbx/actions/workflows/vet.yaml/badge.svg)](https://github.com/reverbdotcom/sbx/actions/workflows/vet.yaml)
[![Release](https://github.com/reverbdotcom/sbx/actions/workflows/release.yml/badge.svg)](https://github.com/reverbdotcom/sbx/actions/workflows/release.yml)

# sbx
Short for _sandbox_, Orchestra CLI tool: `sbx up`

```bash
➜  sbx help

NAME
  sbx - orchestra cli

COMMANDS

  sbx help
      up
      down
      name
      dash
      logs
      web
      graphiql
      version
      info

DESCRIPTION

  command     shorthand     description

  help        h             show the help message.
  up          u             spin up an orchestra sandbox.
  down                      tear down an orchestra sandbox.
  name        n             show the sandbox name.
  dash        d             open the dashboard in a browser.
  logs        l             open the logs in a browser.
  web         w             open the site in a browser.
  graphiql    g             open graphql user interface in a browser.
  version     v             show the version of the sbx cli.
  info        i             show the summary of the sandbox.

USAGE:

  sbx up
  sbx name
```


## Install / Upgrade

Requires `GITHUB_TOKEN` to be set in the environment.

#### brew

```bash
brew tap reverbdotcom/sbx git@github.com:reverbdotcom/sbx.git
brew update
brew install sbx
```

#### golang

```bash
export GOPRIVATE=github.com/reverbdotcom && go install github.com/reverbdotcom/sbx@main
```

#### bash

```bash
VERSION=<grab the latest tag> \
    curl \
        -s\
        -L \
        -o "/tmp/sbx-darwin-arm64.tar.gz" \
        "https://github.com/reverbdotcom/sbx/releases/download/${VERSION}/sbx-darwin-arm64.tar.gz" \
    && tar -xzf /tmp/sbx-darwin-arm64.tar.gz -C /tmp \
    && sudo mv /tmp/sbx /usr/local/bin/sbx
```
## Release

Release is done for `bash` and `brew` installations. We support only darwin-arm64 ( macos m1 ) for now.
To cut a new release, [publish a new tag](https://github.com/reverbdotcom/sbx/releases) following semver.

## Development

`sbx.go` is the main entry point for the CLI tool.
Every command should be a go package. Commands are
configured in `commands/commands.go`.


`make <command>.run` will build and run the command.
This runs live.

#### Test
`make test`

#### Test with another repo

Run, in any orchestra enabled repo.

```bash
export GOPRIVATE=github.com/reverbdotcom && go install github.com/reverbdotcom/sbx@your-test-branch
```

`sbx` now points to your branch version.
