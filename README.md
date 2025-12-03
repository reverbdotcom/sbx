[![Vet](https://github.com/reverbdotcom/sbx/actions/workflows/vet.yaml/badge.svg)](https://github.com/reverbdotcom/sbx/actions/workflows/vet.yaml)
[![Release](https://github.com/reverbdotcom/sbx/actions/workflows/release.yml/badge.svg?event=release)](https://github.com/reverbdotcom/sbx/actions/workflows/release.yml)

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
      progress
      k8s

DESCRIPTION

  command     shorthand     description

  help                      shows the help message.
  up          u             spins up an orchestra sandbox.
  down                      tears down an orchestra sandbox.
  name        n             shows the sandbox name.
  dash        d             opens the dashboard in a browser.
  logs        l             opens the logs in a browser.
  web         w             opens the site in a browser.
  graphiql    g             opens graphql user interface in a browser.
  version     v             shows the version of the sbx cli.
  info        i             shows the summary of the sandbox.
  progress    p             opens deployment progress in a browser.
  k8s                       kubernetes resources explorer. Use 'sbx k8s help' for subcommands.

USAGE:
  sbx <command> [flags]
  sbx up
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

> [!IMPORTANT]
> This is public repo. Do not commit any secrets or sensitive information.
> Simpler for brew install. Though we can work with an internal repo with `HOMEBREW_GITHUB_API_TOKEN` set.

`sbx.go` is the main entry point for the CLI tool.
Every command should be a go package. Commands are
configured in `commands/commands.go`.


`make <command>.run` will build and run the command.
This runs live.

#### Test
`make test`
`make <package>.test`
`make <test file>`

#### Test with another repo

Run, in any orchestra enabled repo.

```bash
export GOPRIVATE=github.com/reverbdotcom && go install github.com/reverbdotcom/sbx@your-test-branch
```

`sbx` now points to your branch version.
