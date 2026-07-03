# Contributing

Thanks for your interest in contributing to `git-share`.

We welcome bug fixes, features, documentation, and suggestions.

---

## Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/git-share.git
   cd git-share
   ```
3. Create a branch from `main`:
   ```bash
   git checkout -b feat/your-feature
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

---

## Branching Strategy

We use **trunk-based development** with feature branches off `main`.

| Prefix | Purpose |
|--------|---------|
| `feat/*` | New features |
| `fix/*` | Bug fixes |
| `test/*` | Adding or updating tests |
| `docs/*` | Documentation changes |
| `refactor/*` | Code restructuring |
| `chore/*` | Maintenance tasks |

Examples:

```
feat/add-mdns-discovery
fix/nil-body-panic
test/cgi-env-windows-paths
```

---

## Development Workflow

1. Make focused, atomic commits (one logical change per commit)
2. Follow the project's coding conventions:
   - Use `gofmt` or `go fmt ./...` before committing
   - Run `go vet ./...` to catch issues
   - Keep dependencies minimal (no new imports beyond `gopkg.in/yaml.v3`)
3. Add or update tests for your changes
4. Verify everything passes:
   ```bash
   go build ./...
   go test -count=1 ./...
   ```
5. Push your branch and open a pull request

---

## Commit Messages (Conventional Commits)

We follow the **Conventional Commits** specification.

### Format

```
<type>: <short description>
```

### Types

| Type | When to use |
|------|-------------|
| `feat` | A new feature |
| `fix` | A bug fix |
| `test` | Adding or updating tests |
| `docs` | Documentation changes |
| `refactor` | Code restructuring (no behavior change) |
| `chore` | Maintenance, tooling, CI, config |

### Examples

```
feat: add web dashboard with live client list
fix: guard nil request body in HandleSmartHTTP
test: add integration tests for bare repo smart HTTP
docs: update installation instructions
refactor: simplify CGI environment construction
chore: update Makefile for cross-compilation
```

### Rules

- Use lowercase
- No period at the end
- Keep under 72 characters
- Use the body for additional context when needed

---

## Pull Request Process

1. Ensure your branch is up to date with `main`:
   ```bash
   git fetch origin
   git rebase origin/main
   ```
2. Verify all checks pass:
   ```bash
   go build ./...
   go vet ./...
   go test -count=1 ./...
   ```
3. Open a pull request targeting `main`
4. Describe:
   - **What** changed
   - **Why** it was needed
   - **How** it was tested
5. Use the PR template at [`.github/PULL_REQUEST_TEMPLATE.md`](.github/PULL_REQUEST_TEMPLATE.md)

---

## Testing

Tests live alongside source code (`*_test.go` in the same package). We prefer:

- **Table-driven tests** for multiple input/output cases
- **Integration tests** using `t.TempDir()` and real `git` commands for backend tests
- Avoid `testing.Main` and external test dependencies

Run all tests:

```bash
go test -count=1 ./...
```

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

## Code of Conduct

This project follows the guidelines in [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md).
Be respectful and constructive in all interactions.

---

Thanks again for contributing.
