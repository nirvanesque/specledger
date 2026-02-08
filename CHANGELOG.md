# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- (no changes yet)

## [1.0.1] - 2026-02-09

### Fixed
- Fixed golangci-lint v1.64.8 compatibility issues (removed deprecated `version` field and `check-shadowing` setting)
- Removed unused code identified by linter (15 items across 8 files)
- Added gosec `#nosec` annotations for appropriate cases (subprocess execution, file permissions)
- Fixed errcheck issues (ignored return values appropriately)

### Changed
- Updated Makefile test targets to exclude integration tests from CI pipeline
- Added `make lint` target for running golangci-lint
- Added `make test-integration` target for running integration tests separately
- Integration tests now skip when prerequisites (mise, beads, perles) are not installed

## [1.0.0] - 2026-01-31

## [1.0.0] - 2026-01-31

### Added
- CLI unification from bash `sl` script and Go `specledger` CLI
- Interactive TUI for project bootstrap using Bubble Tea
- Non-interactive bootstrap with flags (--project-name, --short-code, etc.)
- Specification dependency management commands (add, list, remove, resolve, update)
- GitHub Releases distribution with GoReleaser
- Cross-platform binary builds (Linux, macOS, Windows) for multiple architectures
- Installation scripts for all platforms (bash, PowerShell)
- Debug logging system
- CLI configuration system (~/.config/specledger/config.yaml)
- Error handling with actionable suggestions
- Local dependency caching at ~/.specledger/cache/
- Dependency manifest (specledger.mod) and lockfile (specledger.sum)
- SHA-256 content hashing for cryptographic verification

### Security
- No hardcoded credentials
- Cryptographic verification of cached dependencies
- All configuration is optional
- No data is transmitted externally

### Documentation
- README with installation instructions
- Migration guide from legacy scripts
- Embedded templates for new projects
- AGENTS.md for SpecLedger workflow documentation

### Technical
- Cobra CLI framework for command structure
- Bubble Tea for terminal UI
- Go 1.21+ with go:embed for template embedding
- Cross-platform build support via GoReleaser
- Dependency management via Go modules
