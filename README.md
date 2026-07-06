<div align="center">

# git-share

**Instantly share a local Git repository over HTTP — zero config, zero setup.**

[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)](go.mod)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)
[![Version](https://img.shields.io/badge/version-0.1.0--dev-orange)](.)
[![Go Report Card](https://goreportcard.com/badge/github.com/marcuwynu23/git-share)](https://goreportcard.com/report/github.com/marcuwynu23/git-share)
[![CI](https://github.com/marcuwynu23/git-share/actions/workflows/test.yml/badge.svg)](https://github.com/marcuwynu23/git-share/actions/workflows/test.yml)

➡️ **[Read the full user guide →](USER-GUIDE.md)**

</div>

---

## Table of Contents

- [What Is git-share?](#what-is-git-share)
- [Use Cases](#use-cases)
- [Benefits for Developers](#benefits-for-developers)
- [Advantages Over Other Tools](#advantages-over-other-tools)
- [User Guide](USER-GUIDE.md)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [CLI Reference](#cli-reference)
- [Configuration](#configuration)
- [Example Output](#example-output)
- [CI/CD Integration](#cicd-integration)
- [Development](#development)
- [Architecture](#architecture)
- [Roadmap](#roadmap)
- [Contributing](CONTRIBUTING.md)

---

## What Is git-share?

**git-share** is a cross-platform CLI tool that turns any local Git repository into an instant HTTP server — allowing other machines on the same network to clone, fetch, and push without configuring **SSH, nginx, Docker, or any Git server**.

> Think `python -m http.server`, but for Git.

Run `git share on` inside any Git repo and you're done. A web dashboard shows clone URLs, branch info, and live status. Behind the scenes, it uses Git's built-in Smart HTTP protocol via `git http-backend` CGI — no protocol reimplementation, just a clean bridge.

### What It Does

- **Share** — expose any local Git repository over HTTP with a single command
- **Serve** — full clone, fetch, pull, and push via Git's Smart HTTP protocol
- **Discover** — automatically detects your LAN IP and hostname for network access
- **Dashboard** — built-in dark-themed web UI with repo info and live status
- **Configure** — persistent settings for port, hostname, timeout, and theme
- **Daemonize** — run in the background and stop with `git share off`

### Why Use It?

| Problem | How git-share Solves It |
|---|---|
| Need to share a WIP branch with a teammate | **Instant HTTP server** — share in seconds, no infrastructure |
| Firewall blocks Git port 9418 | **HTTP on standard ports** — works through most firewalls |
| No SSH access to your machine | **No auth required** — just a URL |
| Setting up Git server is too heavy | **Zero config** — no nginx, no Docker, no SSH setup |
| Want a clean UI for repo status | **Web dashboard** — clone URLs, branch info, live status |
| Cross-platform sharing | **Single binary** — Windows, macOS, Linux |

### The Philosophy

1. **Minimal setup, maximum value.** One command, no config files required.
2. **Your process stays yours.** No lock-in, no accounts, no cloud dependency.
3. **Leverage existing tools.** Uses `git http-backend` — the same protocol GitHub uses.

---

## Use Cases

| Scenario | How git-share Helps |
|---|---|
| **Pair programming** | Share your WIP branch with a teammate in seconds |
| **CI/CD testing** | Let build servers pull from your local machine |
| **Code reviews** | Give reviewers HTTP access to feature branches |
| **Local network deploys** | Push to a staging machine without a central server |
| **Teaching / workshops** | Distribute repos to students without GitHub/GitLab |
| **Air-gapped environments** | Share repos in isolated networks with no internet |
| **Quick demos** | Let anyone clone your repo over conference WiFi |

---

## Benefits for Developers

- **One command** — `git share on` and you're live
- **Zero dependencies** — single binary, nothing to install
- **Cross-platform** — Windows, macOS, Linux
- **Smart HTTP** — full Git protocol support (clone, fetch, push)
- **Web dashboard** — browser UI with real-time status
- **LAN discovery** — auto-detect network IP and hostname
- **Read-only / read-write** — control push access per session
- **Config persistence** — settings survive restarts
- **Graceful shutdown** — clean stop on Ctrl+C or `git share off`
- **Daemon mode** — run in background, manage with simple commands

---

## Advantages Over Other Tools

| Aspect | git-share | `git daemon` | `git-http-server` | `simple-http-server` | Manual (nginx + fcgiwrap) |
|---|---|---|---|---|---|
| **Setup time** | ~5 seconds | ~30 seconds | ~1 minute | ~1 minute | ~30 minutes |
| **Protocol** | HTTP (Smart HTTP) | Git-native (9418) | HTTP (dumb only) | HTTP (static files) | HTTP (Smart HTTP) |
| **Push support** | Yes (configurable) | Yes (`receive-pack`) | No | No | Yes |
| **Web dashboard** | Built-in | None | None | None | Custom |
| **Auth ready** | Can add Basic Auth | No built-in | No | No | Yes (nginx) |
| **Firewall friendly** | HTTP/HTTPS ports | Port 9418 required | HTTP ports | HTTP ports | HTTP ports |
| **Browser clone** | `http://ip:9720/repo.git` | Not over HTTP | Dumb HTTP only | Not Git-aware | `http://ip/repo.git` |
| **Cross-platform** | Windows, macOS, Linux | Unix-like only | Cross-platform | Cross-platform | Unix-like only |
| **Daemon support** | Yes (background mode) | Yes | No | No | Yes |
| **Configuration** | YAML file + CLI flags | CLI flags only | CLI flags only | CLI flags only | Multiple files |
| **Binary size** | ~8 MB | Part of Git (~50 MB) | ~5 MB | ~5 MB | N/A |
| **License** | Apache 2.0 | GPL-2.0 | MIT | MIT | Mixed |

---

## Installation

### From source

```bash
git clone https://github.com/marcuwynu23/git-share.git
cd git-share
make build
sudo make install
```

### With Go installed

```bash
go install github.com/marcuwynu23/git-share/cmd/git-share@latest
```

### Binary releases

Download the latest binary from the [Releases page](https://github.com/marcuwynu23/git-share/releases) and place it in your `PATH`.

### Verify

```bash
git share version
```

---

## Quick Start

```bash
# Inside any Git repository:
cd my-project
git share on

# You'll see:
#   Listening:  http://localhost:9720
#   LAN:        http://192.168.1.42:9720
#   Clone:      git clone http://192.168.1.42:9720/
#   Dashboard:  http://localhost:9720/
#   Press Ctrl+C to stop.

# On another machine on the same network:
git clone http://192.168.1.42:9720/
```

---

## CLI Reference

### `git share on`

Start sharing the current repository.

```bash
git share on [flags]
```

| Flag | Default | Description |
|---|---|---|
| `--port <n>` | `9720` | Port to listen on |
| `--readonly` | `false` | Disable push (clone and fetch only) |
| `--readwrite` | `false` | Allow push access |
| `--hostname <addr>` | all interfaces | Bind address (IP or hostname) |
| `--daemon` | `false` | Run in background |

#### Examples

```bash
# Default share
git share on

# Custom port
git share on --port 9000

# Read-only mode
git share on --readonly

# Allow push access
git share on --readwrite

# Bind to a specific interface
git share on --hostname 192.168.1.42

# Run in background
git share on --daemon
```

### `git share off`

Stop a running daemon instance.

```bash
git share off
```

### `git share status`

Show repository and server information.

```bash
git share status
```

### `git share version`

Print version information.

```bash
git share version
```

### `git share config`

View or modify persistent configuration.

```bash
git share config             # Print current config
git share config <key> <value>  # Set a value
```

---

## Configuration

Persistent configuration is stored in `~/.config/git-share/config.yaml` (Linux/macOS) or `%APPDATA%/git-share/config.yaml` (Windows).

```yaml
port: 9720
readonly: false
hostname: ""
timeout: 0
theme: dark
```

### Options

| Key | Type | Default | Description |
|---|---|---|---|
| `port` | int | `9720` | Default port for `git share on` |
| `readonly` | bool | `false` | Default access mode |
| `hostname` | string | `""` | Default bind address (empty = all interfaces) |
| `timeout` | int | `0` | HTTP timeout in seconds (0 = no timeout) |
| `theme` | string | `"dark"` | Dashboard theme |

### Precedence

CLI flags override config file values. Config file values override defaults.

```bash
# Set persistent defaults
git share config port 9000
git share config readonly true
git share config hostname 0.0.0.0
git share config timeout 60
git share config theme light
```

---

## Example Output

### Starting a share

```
$ git share on

  ╔══════════════════════════════════════════╗
  ║            git-share                     ║
  ║    Repository sharing is active          ║
  ╚══════════════════════════════════════════╝

  Repository:  /home/user/my-project
  Branch:      main
  Mode:        Read / Write

  Listening:   http://localhost:9720
  LAN:         http://192.168.1.42:9720
  Clone:       git clone http://192.168.1.42:9720/
  Dashboard:   http://localhost:9720/

  Press Ctrl+C to stop.
```

### Web dashboard

Open `http://localhost:9720/` in a browser to see a dark-themed dashboard showing:

- Repository path and current branch
- Read/Write or Read Only mode badge
- Clone URL with `git clone` command
- Live status refreshed every 5 seconds

### Health endpoint

```bash
curl http://localhost:9720/health
# {"status":"ok"}
```

### Info endpoint

```bash
curl http://localhost:9720/info
# {"repository":"/home/user/my-project","branch":"main","bare":false,"port":9720,"readonly":false,"clone_url":"http://localhost:9720/","lan_url":"http://192.168.1.42:9720/"}
```

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Share Repository
on: [push]

jobs:
  share:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install git-share
        run: |
          go install github.com/marcuwynu23/git-share/cmd/git-share@latest
      - name: Start share
        run: |
          git share on --daemon --port 9720
          sleep 2
          git share status
```

### GitLab CI

```yaml
share:
  script:
    - go install github.com/marcuwynu23/git-share/cmd/git-share@latest
    - git share on --daemon --port 9720
    - sleep 2
    - git clone http://localhost:9720/ /tmp/test-clone
```

---

## Development

### Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.22+ | Compiler and toolchain |
| Git | 2.x | Repository operations (runtime dependency) |
| golangci-lint | latest | Code linting |

### Commands

```bash
make build    # Build the binary
make test     # Run tests with race detector
make vet      # Run go vet
make lint     # Run golangci-lint
make clean    # Remove build artifacts
make run      # Build and run
make install  # Build and copy to install directory
```

### Project Structure

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
└── README.md            # This file
```

---

## Architecture

- **Single binary, zero runtime dependencies** — compiled in Go with no external server required
- **Git Smart HTTP protocol** — delegates to `git http-backend` via CGI over stdout; no protocol reimplementation
- **Smart request routing** — the same root URL `/` serves both the web UI (to browsers) and Git protocol (to `git` CLI), determined by User-Agent and path heuristics
- **PID-based daemon management** — background process managed via a simple PID file in the OS temp directory
- **Dashboard auto-refresh** — the web UI polls `/info` every 5 seconds for live status updates
- **Flag precedence** — CLI flags override config file values, which override defaults

---

## Roadmap

- [x] Core Smart HTTP protocol support (clone, fetch)
- [x] Read-only / read-write mode
- [x] Web dashboard
- [x] LAN IP and hostname discovery
- [ ] Full `off` command with background process management
- [ ] Connection tracking and live client list
- [ ] Authentication (Basic Auth, Bearer Token)
- [ ] mDNS / Zeroconf discovery
- [ ] QR code for mobile cloning
- [ ] End-to-end integration tests

---

## Community Guidelines

This project is governed by the following documents:

- [`CONTRIBUTING.md`](CONTRIBUTING.md) — how to contribute, branch naming, commit messages, PR process
- [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md) — expected behavior and reporting guidelines
- [`SECURITY.md`](SECURITY.md) — how to report security vulnerabilities
- [`USER-GUIDE.md`](USER-GUIDE.md) — comprehensive documentation

---

## License

Apache 2.0 — see [`LICENSE`](LICENSE).

---

## Acknowledgements

- Git's [`http-backend`](https://git-scm.com/docs/git-http-backend) for the Smart HTTP protocol
- The Go standard library for `net/http`, `os/exec`, and `flag`
- [yaml.v3](https://gopkg.in/yaml.v3) for configuration parsing
