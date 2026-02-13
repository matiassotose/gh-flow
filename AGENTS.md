# AGENTS.md

Guidelines for agentic coding agents working on the gh-flow repository.

## Build Commands

```bash
# Build the binary
go build -o gh-flow .
# Or using make:
make build

# Install globally
go install .
# Or:
make install

# Clean build artifacts
make clean

# Update dependencies
make tidy

# Format code
make fmt
```

## Test Commands

```bash
# Run all tests
go test ./...
# Or:
make test

# Run a single test
go test -v -run TestFunctionName ./path/to/package

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test -v ./internal/git

# Run with race detector
go test -race ./...
```

## Lint Commands

```bash
# Run linter (requires golangci-lint)
make lint

# Or use go vet as fallback
go vet ./...

# Check formatting
gofmt -l .
```

## Manual Testing

```bash
# Create test repositories
make setup-test

# Test start command
cd test-repos && ../gh-flow start

# Test finish command
make test-finish

# Clean test repositories
make clean-test
```

## Code Style Guidelines

### Imports
- Use `goimports` format: standard library first, then third-party, then local packages
- Group imports with blank lines between groups
- Use aliases for long import paths when needed
- Example:
  ```go
  import (
      "fmt"
      "os"

      "github.com/spf13/cobra"
      tea "github.com/charmbracelet/bubbletea"

      "gh-flow/internal/git"
      "gh-flow/internal/prompt"
  )
  ```

### Formatting
- Run `gofmt` on all code
- Use tabs for indentation (Go standard)
- Keep lines under 100 characters when possible
- No trailing whitespace
- End files with a newline

### Naming Conventions
- Use `camelCase` for unexported identifiers
- Use `PascalCase` for exported identifiers
- Use `snake_case` for file names
- Use descriptive names; avoid single-letter names except for loop indices
- Interface names should end in `-er` (e.g., `Reader`, `Cmd`)
- Constructor functions start with `New` (e.g., `NewRepoSelectionModel`)

### Types
- Prefer concrete types over interfaces for implementations
- Define interfaces where they are used, not where they are implemented
- Use struct tags for JSON/YAML marshaling
- Document exported types with comments

### Error Handling
- Always check errors explicitly
- Return errors rather than logging them (let caller decide)
- Wrap errors with context using `fmt.Errorf("...: %w", err)`
- Create sentinel errors for specific conditions
- Example:
  ```go
  if err != nil {
      return fmt.Errorf("failed to detect repositories: %w", err)
  }
  ```

### Functions
- Keep functions small and focused on a single task
- Limit to 50 lines when possible
- Return early to reduce nesting
- Use named return values sparingly
- Document exported functions with comments

### Comments
- All exported identifiers must have a comment
- Comments should start with the identifier name
- Use complete sentences with proper punctuation
- Example: `// DetectRepositories finds all git repositories in subdirectories`

### Project Structure
```
gh-flow/
├── cmd/              # Cobra commands
│   ├── root.go
│   ├── start.go
│   └── finish.go
├── internal/         # Internal packages
│   ├── git/         # Git operations
│   ├── prompt/      # TUI prompts
│   └── config/      # Configuration
├── main.go          # Entry point
├── go.mod           # Dependencies
├── Makefile         # Build automation
└── README.md        # Documentation
```

### Testing
- Place tests in `_test.go` files alongside source
- Use table-driven tests
- Mock external dependencies using interfaces
- Name tests clearly: `Test<FunctionName>` or `Test<Scenario>`
- Use `t.Parallel()` for independent tests

### Dependencies
- Minimize external dependencies
- Use go-git for git operations (pure Go)
- Use cobra for CLI framework
- Use bubbletea for TUI
- Pin dependency versions in go.mod

### Git Workflow
- Create feature branches: `feat/description`
- Create hotfix branches: `hotfix/description`
- Make atomic commits with descriptive messages
- Follow conventional commits format
