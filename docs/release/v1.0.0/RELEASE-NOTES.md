# Release Notes — v1.0.0

**Release date:** 2026-07-03

Initial release of git-share — share your local git repository over HTTP with a single command.

---

## [1.0.0] - 2026-07-03

### Added
- CLI with on/off/status/version/config commands
- Smart HTTP support for git clone and push operations over HTTP
- Read-only and read-write sharing modes
- LAN discovery for local network sharing
- Built-in web dashboard UI
- Daemon mode for background operation
- Config file support (`~/.config/git-share/config.yaml`)
- Cross-platform support (Windows, Linux, macOS)
- CI pipeline via GitHub Actions

### Fixed
- Server clone URL displayed port 0 when using random port assignment
- Bare repository push test now branch-agnostic across git versions

### Changed
- Default port changed from 8080 to 9720

---

## Contributors

Thanks to everyone who contributed to this release.

- @marcuwynu23
