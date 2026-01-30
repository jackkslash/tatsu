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
git clone https://github.com/jackkslash/tatsu.git
cd tatsu
go install
```

### Using Go Install

```bash
go install github.com/jackkslash/tatsu@latest
```

Make sure `$GOPATH/bin` is in your `PATH`.

## Quick Start

**Option 1: TUI** (run with no arguments):

```bash
tatsu
```

Launches the terminal UI. Choose Task or PRD mode, enter your task or PRD path, press Enter to run. Live output and validation results appear in the TUI; when done, press `r` or Enter to run again, or `q` to quit.

**Option 2: CLI** (config auto-generates on first run):

```bash
tatsu run "add unit tests to the parser"
```

Tatsu will:
- Auto-generate `tatsu.yaml` if missing
- Run your AI agent with the task
- Validate with your test command
- Retry up to 15 times until tests pass ✅

## Configuration

`tatsu.yaml` is auto-generated on first run. Edit it to customize:

```yaml
agent:
  command: 'opencode run "%s"'  # %s = task description
validate:
  command: 'go test ./...'      # Must exit 0 on success
```

**Examples:**
- Go: `go test ./...`
- Node: `npm test`
- Python: `pytest`
- Multiple: `go test ./... && go vet ./... && golangci-lint run`

**Regenerate config:**
```bash
tatsu generate --force
```

## Usage

### TUI (Interactive)

Run `tatsu` with no arguments to start the TUI:

```bash
tatsu
```

- **Tab** or **←/→** – Switch between **Task** mode and **PRD** mode
- **Task mode** – Type a task description and press **Enter** to run
- **PRD mode** – Input defaults to `prd.md` (editable); press **Enter** to run
- During a run – Live iteration count, agent output, and validation results
- After a run – Scroll with **↑/↓** or **j/k**; **r** or **Enter** to run again; **q** or **Ctrl+C** to quit

### Single Task (CLI)

```bash
tatsu run "task description"
tatsu run -max-iterations 5 "task description"  # Custom retry limit
```

**Examples:**
```bash
tatsu run "implement user authentication"
tatsu run "fix memory leak in cache"
tatsu run -max-iterations 10 "add tests for API endpoints"
```

### PRD (Multiple Tasks)

Execute tasks from a markdown file:

```bash
tatsu prd PRD.example.md
tatsu prd -max-iterations 10 PRD.example.md  # Custom retry limit
```

**PRD Format:**
```markdown
## Tasks
- [ ] create authentication system
- [ ] add user dashboard
- [x] setup database (already done)
```

- `- [ ]` = incomplete (executed)
- `- [x]` = completed (skipped)

**Behavior:**
- Executes incomplete tasks sequentially
- Stops on first failure (after max iterations)
- Each task title becomes the AI agent prompt

### Flags

- `-max-iterations N` - Maximum retry iterations per task (default: 15)
  - Applies to both `run` and `prd` commands
  - Example: `tatsu run -max-iterations 5 "task"`

### Other Commands

```bash
tatsu generate [--force]  # Generate/regenerate config
tatsu version             # Show version
```

## How It Works

```
1. Run AI agent with task → 2. Validate → 3. Pass? Done ✅
                                              ↓ No
                                        4. Retry (default: max 15x)
```

**Each iteration:**
- Shows progress (1/15, 2/15, etc.)
- Displays agent output in real-time
- Shows validation errors on failure
- Stops immediately on success
- Customizable max iterations via `-max-iterations` flag

## Requirements

- [OpenCode](https://github.com/EmbeddedLLM/opencode) installed and in PATH
- Validation command that exits 0 on success
- Go 1.21+ (for building from source)

When Tatsu runs the agent, it sets `OPENCODE_CONFIG_CONTENT` (permission allow), `CI=true`, and `TERM=dumb` so OpenCode runs non-interactively without approval prompts.

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
├── main.go              # CLI entry point (TUI when no args, else CLI)
├── config/              # Configuration management
├── harness/             # AI harness (OpenCode, allow-env for non-interactive)
├── runner/              # Task execution & retry loop (CLI)
├── prd/                 # PRD parsing & execution
├── tui/                 # Terminal UI (Bubbletea)
└── .github/workflows/   # CI/CD
```

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
