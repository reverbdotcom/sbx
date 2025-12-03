# GitHub Copilot Instructions for sbx

## Project Overview

`sbx` (short for "sandbox") is an Orchestra CLI tool written in Go for managing development sandboxes. It provides commands to spin up, tear down, and interact with Orchestra sandboxes.

**Important**: This is a **public repository**. Never commit secrets or sensitive information.

## Architecture

### Main Entry Point
- `sbx.go` - Main entry point that:
  - Verifies environment variables
  - Parses command-line arguments
  - Executes the appropriate command
  - Handles errors uniformly

### Command Structure
- Each command is implemented as a separate Go package in its own directory
- Commands are registered in `commands/commands.go`
- All commands follow a consistent pattern with a `Run()` function

### Key Packages
- `check/` - Environment and repository validation
- `env/` - Environment variable verification
- `parser/` - Command-line argument parsing
- `github/` - GitHub API interactions
- `up/` - Sandbox deployment logic
- `down/` - Sandbox teardown logic
- `name/` - Sandbox naming conventions
- `web/`, `dash/`, `logs/`, `graphiql/` - Browser opening utilities
- `version/` - Version management
- `errr/` - Error formatting
- `debug/` - Debug mode utilities
- `retries/` - Retry logic with exponential backoff

## Code Conventions

### Testing
- All packages should have corresponding `*_test.go` files
- Use table-driven tests with subtests
- Mock external dependencies using function variables
- Test file naming: `<package>_test.go`
- Run tests with: `make test` or `make <package>.test`

### Code Style
- Follow standard Go conventions
- Use `gofmt` for formatting (run with `make fmt`)
- Run `go vet` before committing (run with `make vet`)
- Keep functions small and focused
- Use descriptive variable names

### Error Handling
- Always check and handle errors
- Use `errr.New()` for error formatting
- Return errors up the call stack
- Main function handles final error display

### Environment Variables
- `GITHUB_TOKEN` - Required for GitHub API access
- `DEBUG` - Optional, enables debug output
- Duration can be passed for sandbox lifetime

## Development Workflow

### Building
```bash
make build  # Build the binary
```

### Testing
```bash
make test              # Run all tests
make <package>.test    # Run tests for specific package
make <test_file>       # Run specific test file
```

### Running Commands
```bash
make <command>.run     # Build and run a specific command
```

### Prerequisites
- Go 1.22.0 or later
- `GITHUB_TOKEN` environment variable set
- Orchestra-enabled repository for most commands

## Adding New Commands

1. Create a new directory for the command
2. Implement `Run() (string, error)` function
3. Add tests in `<command>_test.go`
4. Register command in `commands/commands.go`
5. Update README.md with command documentation

## Testing Strategy

### Unit Tests
- Mock external dependencies (git, GitHub API, file system)
- Test error paths
- Test successful execution paths
- Use table-driven tests for multiple scenarios

### Test Helpers
- Use function variables for dependency injection
- Example: `cmdFn` for command execution, `openURL` for browser opening

## CI/CD

### Workflows
- `vet.yaml` - Runs on every push, executes `make vet` and `make test`
- `release.yml` - Triggered on release publication, builds binaries for macOS

### Release Process
- Publish a new tag following semver
- Workflow automatically builds darwin-amd64 and darwin-arm64 binaries
- Updates `version/SBX_VERSION` file
- Uploads release assets

## Common Patterns

### Command Implementation
```go
func Run() (string, error) {
    // 1. Validate prerequisites
    // 2. Perform main logic
    // 3. Return success message or error
}
```

### URL Opening Pattern
```go
func Run() (string, error) {
    url := generateURL()
    err := openURL(url)
    return "Opening " + url, err
}
```

### GitHub API Usage
- Use `github.com/google/go-github/v67` package
- Handle API errors gracefully
- Implement retries for transient failures

## Important Notes

1. **Public Repository**: Never commit secrets, API keys, or sensitive data
2. **Orchestra Requirement**: Most commands require running in an Orchestra-enabled repository
3. **Branch Naming**: Sandbox branches should start with "sandbox"
4. **Minimal Dependencies**: Keep external dependencies minimal
5. **Cross-platform**: Currently supports macOS (amd64 and arm64)

## When Making Changes

1. Ensure tests pass: `make test`
2. Ensure code is formatted: `make fmt`
3. Ensure code passes vet: `make vet`
4. Update documentation if adding features
5. Follow existing patterns and conventions
6. Keep changes minimal and focused
