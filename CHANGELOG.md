# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial unified CLI implementation with `sl` command
- Interactive TUI for project bootstrap
- Non-interactive bootstrap mode with flags
- Specification dependency management commands
- GitHub Releases distribution
- Cross-platform binary builds (Linux, macOS, Windows)
- Installation scripts for all platforms
- Debug logging system
- CLI configuration system
- Error handling with actionable suggestions

### Features

- Single CLI tool for bootstrap and dependency management
- Backward compatibility with `specledger` alias
- Multiple distribution channels (releases, source, package managers)

## [1.0.0] - 2026-01-31

### Added

- CLI unification from bash `sl` script and Go `specledger` CLI
- Interactive TUI for project bootstrap using Bubble Tea
- Non-interactive bootstrap with flags (--project-name, --short-code, etc.)
- Specification dependency management commands
- GitHub Releases distribution
- Cross-platform binary builds (Linux, macOS, Windows)
- Installation scripts for all platforms
- Debug logging system
- CLI configuration system
- Error handling with actionable suggestions

### Security

- Removed hardcoded credentials
- No data is persisted locally
- All configuration is optional

### Documentation

- README with installation instructions
- Migration guide from legacy scripts
- API documentation

### Technical

- Cobra CLI framework
- Bubble Tea TUI library
- Dependency registry with fallback
- Cross-platform build support
- GoReleaser integration

## [Unreleased]

## [1.0.0] - 2026-01-31

### Added

- Initial release of unified CLI
- Interactive TUI for project bootstrap
- Non-interactive bootstrap mode
- Specification dependency management
- GitHub Releases distribution
- Cross-platform support
