# Contributing

Thanks for your interest in contributing to `git-share`.

We welcome bug fixes, features, documentation, and suggestions. Please review the guidelines below before contributing.

---

## Code of Conduct

This project follows the guidelines in [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md).
Be respectful and constructive in all interactions.

---

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.22+ | Compiler and toolchain |
| Git | 2.x | Repository operations (also a runtime dependency) |
| golangci-lint | latest | Code linting (optional, for `make lint`) |

---

## Project Structure

```
.
├── cmd/git-share/        # Binary entry point
├── internal/
│   ├── cli/             # CLI argument parsing and command dispatch
│   ├── config/          # YAML configuration management
│   ├── git/             # Repository detection and Smart HTTP backend
│   ├── server/          # HTTP server lifecycle and route handlers
│   ├── middleware/       # Logger, CORS, Recoverer middleware
│   ├── discovery/       # LAN IP and hostname detection
│   ├── ui/              # Web dashboard
│   └── util/            # Signal handling and helpers
├── docs/                # Documentation
├── tests/               # Integration tests
├── Makefile             # Build, test, lint, install automation
├── go.mod / go.sum      # Go module dependencies
├── README.md            # Project overview and quick start
└── USER-GUIDE.md        # Comprehensive user documentation
```

---

## Makefile Reference

| Target | Command | Description |
|---|---|---|
| `build` | `make build` | Compile the binary (`git-share` or `git-share.exe`) |
| `clean` | `make clean` | Remove build artifacts |
| `test` | `make test` | Run all tests with race detector (`go test -v -race -count=1 ./...`) |
| `vet` | `make vet` | Run `go vet ./...` |
| `lint` | `make lint` | Run `golangci-lint run ./...` |
| `run` | `make run` | Build and run the binary |
| `install` | `make install` | Build and copy to install directory |

```bash
# Build only
make build

# Full verification before submitting a PR
make vet && make test
```

---

## Development Workflow

1. **Fork** the repository on GitHub
2. **Clone** your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/git-share.git
   cd git-share
   ```
3. **Create a feature branch** from `main`:
   ```bash
   git checkout -b feat/your-feature
   ```
4. **Install dependencies**:
   ```bash
   go mod download
   ```
5. **Make focused, atomic commits** — one logical change per commit
6. **Run verification** before pushing:
   ```bash
   go build ./...
   go vet ./...
   go test -count=1 ./...
   ```
7. **Push** your branch and **open a pull request**

---

## Coding Standards

- Run `gofmt` or `go fmt ./...` before committing
- Run `go vet ./...` to catch common issues
- Keep dependencies minimal (currently only `gopkg.in/yaml.v3`)
- Use idiomatic Go — prefer standard library over external packages
- Follow Go's naming conventions:
  - `camelCase` for unexported identifiers
  - `PascalCase` for exported identifiers
  - `ALL_CAPS` for environment variables and constants (sparingly)
- Error handling: always check returned errors; use `fmt.Errorf` with `%w` for wrapping
- Import ordering: standard library → third-party → internal packages

---

## Testing

Tests live alongside source code (`*_test.go` in the same package).

### Expectations

- **Table-driven tests** for multiple input/output cases
- **Integration tests** using `t.TempDir()` and real `git` commands for backend tests
- Avoid `testing.Main` and external test dependencies
- Aim for meaningful coverage — test behavior, not implementation details

### Running tests

```bash
go test -count=1 ./...          # all tests
go test -v -race -count=1 ./... # with race detector and verbose output
go test ./internal/server/...   # specific package
```

### Example

```go
func TestHealthEndpoint(t *testing.T) {
    // Setup
    ts := newTestServer(t)
    defer ts.Close()

    // Execute
    resp, err := http.Get(ts.URL + "/health")

    // Assert
    if err != nil {
        t.Fatal(err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected 200, got %d", resp.StatusCode)
    }
}
```

---

## Commit Conventions

We follow the **Conventional Commits** specification.

### Format

```
<type>(<scope>): <description>
```

### Types

| Type | When to use |
|---|---|
| `feat` | A new feature |
| `fix` | A bug fix |
| `test` | Adding or updating tests |
| `docs` | Documentation changes |
| `refactor` | Code restructuring (no behavior change) |
| `chore` | Maintenance, tooling, CI, config |

### Scope

The scope should be the package or area affected:

| Scope | Area |
|---|---|
| `cli` | Command-line interface |
| `config` | Configuration management |
| `git` | Git backend and repo detection |
| `server` | HTTP server and handlers |
| `ui` | Web dashboard |
| `discovery` | LAN discovery |
| `middleware` | HTTP middleware |

### Examples

```
feat(cli): add --timeout flag for server shutdown
fix(git): handle detached HEAD in status command
test(server): add integration tests for push flow
docs: update installation instructions
refactor(config): simplify merge logic
chore: update Makefile for cross-compilation
```

### Rules

- Use lowercase
- No period at the end
- Keep under 72 characters
- Use the body for additional context when needed

### Breaking changes

Add `BREAKING CHANGE` in the commit body or append `!` after the type/scope:

```
feat(server)!: remove --readonly flag in favor of --mode
```

---

## Pull Request Process

### Checklist

Copy this into your PR description and check off each item:

```markdown
- [ ] My code follows the project style guidelines
- [ ] I have performed a self-review of my code
- [ ] I have added tests that prove my fix or feature works
- [ ] New and existing tests pass locally (`make test`)
- [ ] I have run `go vet ./...` with no new warnings
- [ ] I have made corresponding changes to documentation
- [ ] My commits follow the Conventional Commits format
- [ ] My branch is up to date with `main`
```

### Before submitting

1. Ensure your branch is up to date with `main`:
   ```bash
   git fetch origin
   git rebase origin/main
   ```
2. Verify all checks pass:
   ```bash
   make vet && make test
   ```
3. Open a pull request targeting `main`
4. Describe **what** changed, **why** it was needed, and **how** it was tested
5. Use the PR template at [`.github/PULL_REQUEST_TEMPLATE.md`](.github/PULL_REQUEST_TEMPLATE.md)

### Review criteria

Changes are more likely to be merged when they:
- Have clear, focused commits
- Include tests
- Don't break existing functionality
- Follow the coding conventions above

Changes are unlikely to be merged if they:
- Introduce unnecessary dependencies
- Make large, unfocused changes
- Lack tests for new functionality

---

## Release Process

1. **Tag** the release commit:
   ```bash
   git tag -a v1.0.0 -m "v1.0.0"
   git push origin v1.0.0
   ```
2. **CI** builds binaries for all platforms via the [release workflow](.github/workflows/release.yml)
3. **GitHub Release** is created automatically with cross-platform binaries attached
4. **Update** [CHANGELOG.md](CHANGELOG.md) with the new release notes

---

## Reporting Issues

Use the [bug report template](.github/ISSUE_TEMPLATE/bug_report.md). Include:

- Steps to reproduce
- Expected vs actual behavior
- `git version` output
- OS and architecture

### Feature Requests

Use the [feature request template](.github/ISSUE_TEMPLATE/feature_request.md). Include:

- The problem you're solving
- Your proposed solution
- Alternatives you've considered

---

## Questions

- Open a [GitHub Discussion](https://github.com/marcuwynu23/git-share/discussions)
- File an [issue](https://github.com/marcuwynu23/git-share/issues)

---

Thanks again for contributing!
