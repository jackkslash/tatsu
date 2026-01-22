# Tatsu

An iterative AI coding assistant that runs tasks in a validation loop until tests pass.

## What is Tatsu?

Tatsu automates the code-test-fix cycle by:
1. Running your AI agent with a task description
2. Validating the result (tests, lint, etc.)
3. Automatically retrying if validation fails
4. Stopping when validation passes

Think of it as a persistent coding assistant that doesn't give up until your tests pass.

## Installation

### From Source

```bash
git clone https://github.com/jack/tatsu.git
cd tatsu
go install
```

### Using Go Install

```bash
go install github.com/jack/tatsu@latest
```

Make sure `$GOPATH/bin` is in your `PATH`.

## Quick Start

1. **Create `tatsu.yaml` in your project:**

```yaml
agent:
  command: 'opencode run "%s"'
validate:
  command: 'go test ./...'
```

2. **Run a task:**

```bash
tatsu run "add unit tests to the parser"
```

Tatsu will:
- Run OpenCode with your task
- Run your validation command
- Retry up to 15 times if validation fails
- Stop when tests pass ✅

## Configuration

Create a `tatsu.yaml` file in your project root:

```yaml
agent:
  # Command to run your AI agent
  # %s will be replaced with the task description
  command: 'opencode run "%s"'

validate:
  # Command to check if the task is complete
  # Should exit with code 0 on success
  command: 'go test ./...'
```

### Configuration Options

- **`agent.command`** - Command to execute your AI coding assistant. Use `%s` as a placeholder for the task description.
- **`validate.command`** - Command that validates the work. Should exit with code 0 on success, non-zero on failure.

### Example Configurations

**Go project:**
```yaml
agent:
  command: 'opencode run "%s"'
validate:
  command: 'go test ./...'
```

**Node.js project:**
```yaml
agent:
  command: 'opencode run "%s"'
validate:
  command: 'npm test'
```

**Python project:**
```yaml
agent:
  command: 'opencode run "%s"'
validate:
  command: 'pytest'
```

**Multiple validations:**
```yaml
agent:
  command: 'opencode run "%s"'
validate:
  command: 'go test ./... && go vet ./... && golangci-lint run'
```

## Usage

### Basic Usage

```bash
tatsu run "task description"
```

### Examples

```bash
# Add a feature
tatsu run "implement user authentication"

# Fix a bug
tatsu run "fix the memory leak in the cache"

# Add tests
tatsu run "add tests for the API endpoints"

# Refactor code
tatsu run "refactor the parser to use a visitor pattern"
```

### Version

```bash
tatsu version
```

## How It Works

```
1. Run agent with task description
2. Run validation command
3. If validation passes → done ✅
4. If validation fails → retry (max 15 times)
5. Show clear feedback at each step
```

### Iteration Loop

Tatsu will:
- Show iteration count (1/15, 2/15, etc.)
- Display agent output in real-time
- Show validation output on failure
- Stop immediately when validation succeeds
- Stop after 15 iterations if still failing

## Requirements

- **AI Coding Assistant** - Currently supports [OpenCode](https://github.com/EmbeddedLLM/opencode)
- **Validation Command** - Any command that exits 0 on success (tests, linters, etc.)
- **Go 1.21+** - If building from source

## Development

### Building

```bash
make build
# or
go build
```

### Testing

```bash
make test
# or
go test ./...
```

### Running Tests with Coverage

```bash
make coverage
# or
go test ./... -cover
```

### Linting

```bash
make lint
# or
go vet ./...
```

### All-in-One

```bash
make all  # Runs fmt, lint, test, and build
```

See `Makefile` for all available commands.

## Project Structure

```
tatsu/
├── main.go              # CLI entry point
├── config/              # Configuration loading
│   └── config.go
├── harness/             # AI harness interfaces
│   ├── harness.go
│   └── opencode.go
├── runner/              # Task execution engine
│   └── runner.go
└── .github/workflows/   # CI/CD
    └── ci.yml
```

## Limitations

- Currently supports OpenCode only (more harnesses coming soon)
- Maximum 15 iterations per task
- Single task at a time
- Manual configuration required (no auto-detection)

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT

## Acknowledgments

Built with Go and designed for developers who want AI assistance that actually works.
