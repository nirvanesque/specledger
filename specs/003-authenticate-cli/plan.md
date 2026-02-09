# Implementation Plan: Authenticate CLI

**Branch**: `feat-003/authenticate-cli` | **Date**: 2026-02-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/003-authenticate-cli/spec.md`
**Status**: Implemented

## Summary

Implement browser-based OAuth-style authentication for the SpecLedger CLI, enabling users to sign in via their browser and automatically capture credentials via a local callback server. Support CI/headless environments with direct token authentication flags.

## Technical Context

**Language/Version**: Go 1.24.2
**Primary Dependencies**:
- `github.com/spf13/cobra v1.10.2` - CLI framework
- `net/http` - Callback server
- `encoding/json` - Credential serialization
- `os/exec` - Browser launching

**Storage**: File-based credentials at `~/.specledger/credentials.json`
**Testing**: Go standard testing (`go test`)
**Target Platform**: Cross-platform (macOS, Linux, Windows)
**Project Type**: Single CLI application
**Performance Goals**: Authentication completes in <60 seconds, status check <1 second
**Constraints**: Credentials stored with 0600 permissions, 5-minute authentication timeout
**Scale/Scope**: Single-user local authentication

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with principles from `.specify/memory/constitution.md`:

- [x] **Specification-First**: Spec.md complete with prioritized user stories (5 stories documented)
- [x] **Test-First**: Test strategy defined (unit tests for credentials, integration tests for auth flow)
- [x] **Code Quality**: Go standard formatting with `gofmt`
- [x] **UX Consistency**: User flows documented in spec.md acceptance scenarios
- [x] **Performance**: Metrics defined (<60s auth, <1s status check, <10s refresh)
- [x] **Observability**: CLI provides clear feedback at each authentication step
- [x] **Issue Tracking**: Feature tracked on branch feat-003/authenticate-cli

**Complexity Violations**: None identified

## Project Structure

### Documentation (this feature)

```text
specs/003-authenticate-cli/
├── spec.md              # Feature specification
├── plan.md              # This file
├── research.md          # Technical decisions
├── data-model.md        # Data structures
├── quickstart.md        # Developer guide
└── checklists/
    └── requirements.md  # Validation checklist
```

### Source Code (repository root)

```text
pkg/cli/
├── auth/
│   ├── browser.go       # Cross-platform browser opening
│   ├── client.go        # Token refresh HTTP client
│   ├── credentials.go   # Credential storage/loading
│   └── server.go        # Local OAuth callback server
├── commands/
│   └── auth.go          # sl auth command implementation
├── config/
├── dependencies/
├── logger/
└── tui/

cmd/
└── main.go              # CLI entry point
```

**Structure Decision**: Single project structure following Go conventions. Auth functionality isolated in `pkg/cli/auth/` package for clean separation of concerns.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |

## Implementation Phases

### Phase 0: Research (Complete)
- OAuth callback flow patterns researched
- Cross-platform browser opening methods identified
- Secure credential storage patterns established

### Phase 1: Design (Complete)
- Data models defined (Credentials, CallbackResult, CallbackServer)
- API contracts established (callback endpoint, refresh endpoint)
- File structure determined

### Phase 2: Implementation (Complete)
- All 5 user stories implemented and functional
- Commands: login, logout, status, refresh
- Token flags: --token, --refresh for CI/headless
