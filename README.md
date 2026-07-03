<div align="center">

# git-share

**Instantly share a local Git repository over HTTP — zero config, zero setup.**

[![Go Version](https://img.shields.io/badge/go-1.22+-00ADD8?logo=go)](go.mod)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue)](LICENSE)
[![Version](https://img.shields.io/badge/version-0.1.0--dev-orange)](.)
[![Go Report Card](https://goreportcard.com/badge/github.com/marcuwynu23/git-share)](https://goreportcard.com/report/github.com/marcuwynu23/git-share)

</div>

---

## What is git-share?

`git-share` is a cross-platform CLI tool that turns any local Git repository into an instant HTTP server — allowing other machines on the same network to clone, fetch, and push without configuring **SSH, nginx, Docker, or any Git server**.

> Think `python -m http.server`, but for Git.

Run `git share on` inside any Git repo and you're done. A web dashboard shows clone URLs, branch info, and live status. Behind the scenes, it uses Git's built-in Smart HTTP protocol via `git http-backend` CGI — no protocol reimplementation, just a clean bridge.

---

## Features

- **Zero configuration** — run a single command, share instantly
- **Smart HTTP protocol** — full Git clone, fetch, pull, push support
- **Cross-platform** — Windows, macOS, Linux (single binary, no dependencies)
- **Web dashboard** — dark-themed UI with repo info and live status
- **LAN discovery** — automatically detects your network IP and hostname
- **Read-only / read-write modes** — control push access per session
- **Graceful shutdown** — Ctrl+C cleanly stops the server
- **Config persistence** — port, hostname, timeout, and theme stored in `~/.config/git-share/config.yaml`

---

## vs `git daemon`

| | git-share | `git daemon` |
|---|---|---|
| **Protocol** | HTTP (Smart HTTP) | Git-native (port 9418) |
| **Auth ready** | Can add Basic Auth / Bearer Token | No auth built-in |
| **Web dashboard** | Built-in browser UI | None |
| **Firewall friendly** | HTTP/HTTPS ports (usually open) | Requires port 9418 open |
| **Push** | Yes (configurable read-only/read-write) | Yes (with `--enable=receive-pack`) |
| **Browser clone** | `http://lan-ip:8080/repo.git` | Not available over HTTP |
| **Setup** | Single command, no config | Requires `git daemon --base-path=...` arguments |
| **Platform** | Cross-platform binary | Unix-like only (no native Windows) |

Use **git-share** when you need a quick, browser-friendly share over HTTP with minimal friction. Use **git daemon** for a lightweight Git-native solution on Unix networks.

---

## Use Cases

- **Pair programming** — share your WIP branch with a teammate in seconds
- **CI/CD testing** — let build servers pull from your local machine
- **Code reviews** — give reviewers HTTP access to feature branches
- **Local network deploys** — push to a staging machine without a central server
- **Teaching / workshops** — distribute repositories to students without GitHub/GitLab
- **Air-gapped environments** — share repos in isolated networks with no internet
- **Quick demos** — let anyone clone your repo over the conference WiFi

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

---

## Usage

```bash
# Inside any Git repository:
git share on

# Listening:
#   http://localhost:8080
# LAN:
#   http://192.168.1.42:8080
# Clone:
#   git clone http://192.168.1.42:8080/repo.git
# Press Ctrl+C to stop.
```

### Commands

| Command | Description |
|---------|-------------|
| `git share on` | Start sharing the current repository |
| `git share off` | Stop sharing (requires background process support) |
| `git share status` | Show repository and server status |
| `git share version` | Print version information |
| `git share config` | View or modify persistent configuration |

### Options

| Flag | Description |
|------|-------------|
| `--port <n>` | Port to listen on (default: `8080`) |
| `--readonly` | Start in read-only mode (clone only, no push) |
| `--readwrite` | Start in read-write mode (allow push) |
| `--hostname <addr>` | Bind address (default: all interfaces) |

### Examples

```bash
# Start on a custom port
git share on --port 9000

# Allow push access
git share on --readwrite

# Start in read-only mode
git share on --readonly

# Bind to a specific interface
git share on --hostname 192.168.1.42

# Configure persistent defaults
git share config port 9000
git share config readonly false
```

### Pushing to a shared repo

Git refuses to push to the currently checked-out branch of a non-bare repo (it would corrupt the working tree). Push to a different branch instead:

```bash
# From the client — push to a feature branch
git push origin main:wip/my-changes

# Or create a local branch and push it
git checkout -b my-feature
git push -u origin my-feature
```

Then on the machine serving the repo, merge it into `main`:

```bash
# On the server machine
git merge my-feature
git branch -d my-feature
```

If you need to push directly to `main`, either:

- Serve a bare repository (`git init --bare`), or
- Set `git config receive.denyCurrentBranch ignore` on the served repo (not recommended — requires manual `git reset --hard` to sync the working tree).

---

## Configuration

Persistent configuration is stored in `~/.config/git-share/config.yaml` (Linux/macOS) or `%APPDATA%/git-share/config.yaml` (Windows).

```yaml
port: 8080
readonly: true
hostname: ""
timeout: 0
theme: dark
```

Set values with `git share config <key> <value>`.

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
├── docs/                # Documentation (coming soon)
├── tests/               # Integration tests (coming soon)
├── Makefile             # Build, test, lint, install automation
├── go.mod / go.sum      # Go module dependencies
└── README.md            # This file
```

---

## Roadmap

- [x] Core Smart HTTP protocol support (clone, fetch)
- [x] Read-only / read-write mode
- [x] Web dashboard
- [x] LAN IP and hostname discovery
- [ ] `off` command with background process management
- [ ] Full `status` command with server state
- [ ] Client connection tracking
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

---

## License

Apache 2.0 — see [`LICENSE`](LICENSE).

---

## Acknowledgements

- Git's [`http-backend`](https://git-scm.com/docs/git-http-backend) for the Smart HTTP protocol
- The Go standard library for `net/http`, `os/exec`, and `flag`
- [yaml.v3](https://gopkg.in/yaml.v3) for configuration parsing

---

Happy coding!
