# User Guide

Comprehensive reference for **git-share** — instantly share a local Git repository over HTTP.

> 📖 See the [README](README.md) for a quick overview. See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup.

---

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Configuration](#configuration)
- [Web Dashboard](#web-dashboard)
- [API Endpoints](#api-endpoints)
- [Concepts](#concepts)
- [Workflows](#workflows)
- [CI/CD Integration](#cicd-integration)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

---

## Installation

### Prerequisites

| Requirement | Details |
|---|---|
| Git | 2.x (required — git-share runs `git http-backend` under the hood) |
| Operating system | Windows 10+, macOS 12+, Linux (kernel 4.x+) |
| Architecture | amd64, arm64 |

### Method 1: Binary download (recommended)

1. Go to the [Releases page](https://github.com/marcuwynu23/git-share/releases)
2. Download the archive for your OS and architecture
3. Extract the binary and place it in your `PATH`:
   - **Linux/macOS:** `/usr/local/bin/git-share`
   - **Windows:** `C:\Users\<you>\AppData\Local\Programs\git-share\git-share.exe` (ensure this directory is in your `PATH`)

### Method 2: Go install

```bash
go install github.com/marcuwynu23/git-share/cmd/git-share@latest
```

Requires [Go 1.22+](https://go.dev/dl/). The binary is placed in `$(go env GOPATH)/bin`.

### Method 3: Build from source

```bash
git clone https://github.com/marcuwynu23/git-share.git
cd git-share
make build
sudo make install     # Linux/macOS
# or manually copy the binary to a directory in your PATH
```

### Verify

```bash
git share version
# git-share version 0.1.0
```

---

## Quick Start

### Share a repository

```bash
cd /path/to/your/project
git share on
```

You'll see output like:

```
  Repository:  /home/user/my-project
  Branch:      main
  Mode:        Read / Write

  Listening:   http://localhost:9720
  LAN:         http://192.168.1.42:9720
  Clone:       git clone http://192.168.1.42:9720/
  Dashboard:   http://localhost:9720/
```

### Clone from another machine

```bash
git clone http://192.168.1.42:9720/
```

### Stop sharing

Press **Ctrl+C** in the terminal where `git share on` is running. If running in daemon mode:

```bash
git share off
```

---

## Command Reference

### `git share on`

Start sharing the current repository over HTTP.

```bash
git share on [flags]
```

| Flag | Type | Default | Description |
|---|---|---|---|
| `--port` | int | `9720` | TCP port to listen on |
| `--readonly` | bool | `false` | Only allow clone and fetch (no push) |
| `--readwrite` | bool | `false` | Allow push access |
| `--hostname` | string | `""` (all interfaces) | IP address or hostname to bind to |
| `--daemon` | bool | `false` | Run as a background process |

**Flag precedence:** CLI flags > config file > built-in defaults.

#### Examples by Use Case

**Default share (read-write):**

```bash
git share on
```

**Read-only for code review:**

```bash
git share on --readonly
```

**Custom port for firewall compatibility:**

```bash
git share on --port 9000
```

**Bind to a specific network interface:**

```bash
git share on --hostname 192.168.1.42
```

**Background daemon:**

```bash
git share on --daemon
```

---

### `git share off`

Stop a running daemon.

```bash
git share off
```

This reads the PID file from the temp directory, sends a termination signal, and removes the PID file.

---

### `git share status`

Display repository and server information.

```bash
git share status
```

Prints the current repository path, branch, and whether the repo is bare.

---

### `git share version`

Print version information.

```bash
git share version
```

---

### `git share config`

View or modify persistent configuration.

```bash
git share config                  # Print current configuration
git share config <key> <value>    # Set a configuration value
```

#### Examples

```bash
git share config port 9000
git share config readonly true
git share config hostname 0.0.0.0
git share config timeout 60
git share config theme light
```

---

## Configuration

### Config file location

| OS | Path |
|---|---|
| Linux / macOS | `~/.config/git-share/config.yaml` |
| Windows | `%APPDATA%\git-share\config.yaml` |

### Full schema

```yaml
port: 9720
readonly: false
hostname: ""
timeout: 0
theme: dark
```

### Field reference

| Key | Type | Default | Description |
|---|---|---|---|
| `port` | int | `9720` | Default TCP port when `--port` is not provided |
| `readonly` | bool | `false` | Default access mode when `--readonly`/`--readwrite` not provided |
| `hostname` | string | `""` | Default bind address (empty string binds to all interfaces) |
| `timeout` | int | `0` | HTTP read/write timeout in seconds (`0` = no timeout) |
| `theme` | string | `"dark"` | Dashboard theme identifier (affects CSS class) |

### Precedence rules

```
CLI flags (highest)
        ↓
Config file (~/.config/git-share/config.yaml)
        ↓
Built-in defaults (lowest)
```

### Configurable keys

| Key | Validation | Example |
|---|---|---|
| `port` | Must be a positive integer | `git share config port 9000` |
| `readonly` | Accepts `true`, `yes`, `1`; anything else = `false` | `git share config readonly true` |
| `hostname` | Any string | `git share config hostname 0.0.0.0` |
| `timeout` | Must be a non-negative integer (seconds) | `git share config timeout 60` |
| `theme` | Any string | `git share config theme light` |

---

## Web Dashboard

When you run `git share on`, open `http://localhost:9720/` in a browser to see the web dashboard.

### Dashboard features

- **Repository path** — full filesystem path to the shared repo
- **Branch** — current Git branch name
- **Mode badge** — "Read / Write" or "Read Only" with a status dot
- **Clone URL** — pre-built `git clone` command (prefers LAN IP, falls back to localhost)
- **Auto-refresh** — status updates every 5 seconds

### Layout

```text
╔══════════════════════════════════════╗
║           git-share                  ║
║   Repository sharing is active       ║
║                                      ║
║   ● Repository: /path/to/repo       ║
║   ● Branch:     main                 ║
║   ● Mode:       Read / Write        ║
║   ● Clients:    0                    ║
║   ● Uptime:     -                    ║
║                                      ║
║   ┌──────────────────────────────┐   ║
║   │ git clone http://...         │   ║
║   └──────────────────────────────┘   ║
╚══════════════════════════════════════╝
```

---

## API Endpoints

### `GET /health`

Returns the server health status.

```bash
curl http://localhost:9720/health
```

Response:

```json
{"status":"ok"}
```

### `GET /info`

Returns repository metadata and server information.

```bash
curl http://localhost:9720/info
```

Response:

```json
{
  "repository": "/absolute/path/to/repo",
  "branch": "main",
  "bare": false,
  "port": 9720,
  "readonly": false,
  "clone_url": "http://localhost:9720/",
  "lan_url": "http://192.168.1.42:9720/"
}
```

### `GET /git/*`

Git Smart HTTP protocol endpoint. Used internally by `git clone`, `git fetch`, `git push`. Not intended for direct browser access.

---

## Concepts

### Smart HTTP Protocol

git-share uses Git's **Smart HTTP** protocol, the same protocol used by GitHub and GitLab. Unlike the "Dumb HTTP" protocol (which only serves raw objects), Smart HTTP:

- Negotiates the smallest possible data transfer
- Supports push (receive-pack)
- Uses Git's native wire protocol over HTTP

### How it works

```
User runs:  git share on
               │
               ▼
    Starts HTTP server on port 9720
               │
          ┌────┴────┐
          ▼         ▼
    Browser      Git client
    (dashboard)  (clone/push)
          │         │
          ▼         ▼
    GET /info   GET /git/HEAD
    GET /       POST /git/git-receive-pack
```

### Request routing

The server distinguishes between browser and Git clients using:

- **User-Agent header** — Git clients send `git/X.Y.Z`
- **Query string** — Git requests include `?service=git-upload-pack` or `?service=git-receive-pack`
- **Path patterns** — Git requests target paths like `/info/refs`, `/HEAD`, `/git-upload-pack`, `/git-receive-pack`
- **`.git` in URL path** — indicates a Git request

### Read-only vs read-write mode

| Mode | Clone | Fetch | Pull | Push |
|---|---|---|---|---|
| Read-only | ✓ | ✓ | ✓ | ✗ |
| Read-write | ✓ | ✓ | ✓ | ✓ |

In read-write mode, git-share sets `GIT_CONFIG_COUNT=1` with `http.receivepack=true`, enabling push via `git http-backend`.

---

## Workflows

### Ad-hoc sharing with a teammate

```bash
# Machine A (server):
cd ~/projects/my-app
git share on --readwrite

# Machine B (client):
git clone http://192.168.1.42:9720/
# ... make changes ...
git push origin main:wip/my-review
```

### Code review without push access

```bash
# Reviewer starts in read-only mode:
cd ~/projects/my-app
git share on --readonly

# Reviewer shares the URL:
# "Clone my branch at http://192.168.1.42:9720/"
```

### Pushing to a shared repo

Git refuses to push to the currently checked-out branch of a non-bare repo. Push to a different branch instead:

```bash
# Client: push to a feature branch
git push origin main:wip/my-changes

# Or create and push a local branch
git checkout -b my-feature
git push -u origin my-feature
```

On the server machine, merge the changes:

```bash
git merge my-feature
git branch -d my-feature
```

### Using a bare repository

For direct push to `main`, serve a bare repository:

```bash
git init --bare ~/shared-repo.git
cd ~/shared-repo.git
git share on --readwrite

# Client:
git push http://192.168.1.42:9720/ main
```

### Background daemon

```bash
# Start in background
git share on --daemon

# Check on it later
git share status

# Stop it
git share off
```

---

## CI/CD Integration

### GitHub Actions — Test with shared repo

```yaml
name: Integration Test
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - name: Build git-share
        run: go build -o git-share ./cmd/git-share
      - name: Start share server
        run: |
          ./git-share on --daemon --port 9720
          sleep 2
      - name: Clone from server
        run: |
          git clone http://localhost:9720/ /tmp/test-clone
          diff -r . /tmp/test-clone
      - name: Stop server
        run: ./git-share off
```

### GitLab CI

```yaml
integration:
  stage: test
  script:
    - go build -o git-share ./cmd/git-share
    - ./git-share on --daemon --port 9720
    - sleep 2
    - git clone http://localhost:9720/ /tmp/test-clone
    - diff -r . /tmp/test-clone
    - ./git-share off
```

### Jenkins pipeline

```groovy
stage('Share Repository') {
    steps {
        sh 'go build -o git-share ./cmd/git-share'
        sh './git-share on --daemon --port 9720'
        sh 'sleep 2'
        sh 'git clone http://localhost:9720/ /tmp/test-clone'
        sh './git-share off'
    }
}
```

---

## Troubleshooting

| Problem | Cause | Fix |
|---|---|---|
| `git: command not found` when starting server | Git is not installed on the server machine | Install Git 2.x from [git-scm.com](https://git-scm.com) |
| `connection refused` when cloning | Server is not running or wrong host/port | Verify `git share on` is running; check port and hostname |
| `git received HTTP 403` | Server is in read-only mode and you're trying to push | Start with `--readwrite` or configure `readonly: false` |
| `cannot push to checked out branch` | Git protects non-bare repos from branch corruption | Push to a different branch: `git push origin main:wip/feature` |
| `port already in use` | Another process is using the port | Use `--port` to specify a different port, or kill the existing process |
| `no repositories available` | Not inside a Git repository | Run `git share on` inside a directory with a `.git` folder |
| Server starts but LAN URL shows wrong IP | Multiple network interfaces | Use `--hostname` to bind to the correct interface |
| `git share off` does nothing | No daemon is running; the PID file is missing | Press Ctrl+C in the terminal where `git share on` is running |
| `fatal: repository 'http://...' not found` | URL is missing the trailing slash or repo path | Ensure the URL ends with `/` — e.g., `http://192.168.1.42:9720/` |

---

## FAQ

**Q: Do I need to install anything besides git-share?**

A: You need Git itself (the `git` command). git-share delegates to `git http-backend` under the hood.

**Q: Can I use this over the internet, not just LAN?**

A: Technically yes, but it's not recommended without authentication. git-share has no built-in auth yet. For internet use, wrap it behind a reverse proxy with Basic Auth.

**Q: Is my data encrypted?**

A: No. git-share serves plain HTTP. Do not use it over untrusted networks without a TLS proxy (e.g., nginx + Let's Encrypt).

**Q: Can multiple people clone at the same time?**

A: Yes. The Go `net/http` server handles concurrent connections.

**Q: What happens if I push to a non-bare repo?**

A: Git pushes are accepted (in read-write mode), but the working tree is not automatically updated. You need to run `git reset --hard` on the server to sync. For direct push to `main`, use a bare repository.

**Q: How do I change the default port permanently?**

A: `git share config port 9000`

**Q: How do I check if git-share is running?**

A: `git share status` if running in the foreground, or check the web dashboard at `http://localhost:9720/`.

**Q: Can I run multiple instances on different ports?**

A: Yes. Run `git share on --port 9721` in a second terminal.

**Q: Does it work with monorepos?**

A: Yes. git-share serves whatever Git repository you're currently in, monorepo or not.

**Q: How do I uninstall git-share?**

A: Delete the binary and remove `~/.config/git-share/` (Linux/macOS) or `%APPDATA%\git-share\` (Windows).

**Q: Does git-share support shallow clones?**

A: Yes. `git clone --depth 1 http://192.168.1.42:9720/` works.

**Q: Where does the PID file live?**

A: In your OS temp directory (e.g., `/tmp/.pstemp` on Linux, `%TEMP%\.pstemp` on Windows).

---

## Additional Resources

- [README.md](README.md) — project overview and quick start
- [CONTRIBUTING.md](CONTRIBUTING.md) — development setup and contributing guidelines
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) — community standards
- [SECURITY.md](SECURITY.md) — security policies
- [GitHub Issues](https://github.com/marcuwynu23/git-share/issues) — report bugs and request features
