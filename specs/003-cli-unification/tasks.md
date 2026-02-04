# Implementation Tasks: CLI Unification

**Feature**: 003-cli-unification
**Epic**: sl-6ak
**Date**: 2026-01-30
**Source**:
- Spec: [spec.md](./spec.md)
- Plan: [plan.md](./plan.md)
- Research: [research.md](./research.md)
- Data Model: [data-model.md](./data-model.md)
- Contracts: [contracts/CLI-INTERFACE.md](./contracts/CLI-INTERFACE.md)
- Quickstart: [quickstart.md](./quickstart.md)

## User Stories

| ID | Title | Priority | Independent Test |
|----|-------|----------|------------------|
| US1 | Single Unified CLI Tool | P1 | Install unified CLI and verify it supports both bootstrap and dependency management commands |
| US2 | GitHub Releases | P1 | Download CLI from GitHub releases and verify it executes correctly |
| US3 | Self-Built Binaries | P2 | Build CLI from source and verify it works correctly |
| US4 | Self-Hosted / Local Binaries | P2 | Run CLI from non-PATH location and verify it works |
| US5 | UVX Style | P2 | Execute standalone CLI and verify help output works |
| US6 | Package Manager Integration | P3 | Install via package manager and verify CLI is accessible |

## Task Breakdown by Phase

### Phase 1: Setup (Foundation)

Tasks in this phase provide the shared infrastructure needed by all user stories.

```bash
# Filter tasks
bd list --label "spec:003-cli-unification" --label "phase:setup" -n 5
bd ready --label "spec:003-cli-unification" --label "phase:setup" -n 5
```

**T001**: Setup TUI integration with Bubble Tea
- Implement terminal detection utility
- Add interactive prompt components for bootstrap
- Configure fallback to plain CLI when TUI unavailable
- **Acceptance**: Terminal detection works, prompts display correctly in TUI mode

**T002**: Implement dependency registry with fallback
- Create DependencyRegistry struct with local, mise, and interactive fallback
- Implement gum client integration
- Implement mise client integration
- Add clear error messages with installation instructions
- **Acceptance**: Missing dependencies prompt user instead of failing

**T003**: Create CLI configuration system
- Implement config file loading from ~/.config/specledger/config.yaml
- Add validation for configuration values
- Store user preferences (shell, theme, language, etc.)
- **Acceptance**: Configuration loads and validates correctly

**T004**: Setup debug logging system
- Implement logger with debug-level output to stderr
- Add log formatting and context
- **Acceptance**: Debug logs appear on stderr during execution

**T005**: Update main.go with unified command structure
- Rename `specledger` to `sl` as primary command
- Add `specledger` alias for backward compatibility
- Implement `--help` and `--version` flags
- **Acceptance**: Both `sl` and `specledger` work, version and help show correctly

### Phase 2: User Story 1 - Single Unified CLI Tool (P1)

Core unification goal: single CLI for both bootstrap and dependency management.

```bash
# Filter tasks
bd list --label "spec:003-cli-unification" --label "story:US1" -n 5
bd ready --label "spec:003-cli-unification" --label "story:US1" -n 5
```

**T006**: Implement `sl bootstrap` / `sl new` command (TUI) ✅ COMPLETED
- Create bootstrap command with Bubble Tea TUI
- Add prompts for project name, short code, playbook, agent shell
- Integrate dependency registry for gum/mise checks
- Display success/failure messages
- **Acceptance**: Interactive TUI prompts work correctly in terminal
- **Status**: COMPLETED - Full Bubble Tea TUI implementation with:
  - 5-step wizard (Project Name → Short Code → Playbook → Shell → Confirm)
  - Real-time text input with validation
  - Arrow key navigation for selections
  - Error handling and graceful exit

**T007**: Implement `sl bootstrap` / `sl new` command (non-interactive) ✅ COMPLETED
- Add flags: --project-name, --short-code, --playbook, --shell, --ci
- Validate inputs before execution
- Handle existing project directory detection
- **Acceptance**: Non-interactive bootstrap works with flags
- **Status**: COMPLETED - Full non-interactive implementation with:
  - All required flags (--project-name, --short-code, --playbook, --shell, --ci, --project-dir)
  - Input validation
  - Error handling for existing directories

**T008**: Integrate existing deps commands into unified CLI ✅ COMPLETED
- Map specledger deps subcommands to sl deps subcommands
- Ensure all commands work with sl binary
- Add specledger alias for backward compatibility
- **Acceptance**: All deps commands work via sl
- **Status**: COMPLETED - All dependency management commands integrated:
  - `sl deps add`, `list`, `resolve`, `remove`, `update`
  - `sl refs validate`, `list`
  - `sl vendor` subcommands
  - `sl conflict check`, `detect`

**T009**: Implement CLI error handling
- Standardize error messages with actionable suggestions
- Handle CI/CD non-interactive environments
- Implement exit code 0 for success, 1 for any failure
- **Acceptance**: All error scenarios show clear, helpful messages

**T010**: Test unified CLI end-to-end
- Test interactive bootstrap with TUI
- Test non-interactive bootstrap with flags
- Test all deps commands
- Verify backward compatibility with specledger alias
- **Acceptance**: All acceptance scenarios from spec pass

### Phase 3: User Story 2 - GitHub Releases (P1)

Enable users to install CLI directly from GitHub releases.

```bash
# Filter tasks
bd list --label "spec:003-cli-unification" --label "story:US2" -n 5
bd ready --label "spec:003-cli-unification" --label "story:US2" -n 5
```

**T011**: Configure GoReleaser for cross-platform builds
- Create .goreleaser.yaml configuration
- Add builds for Linux, macOS, Windows (amd64, arm64)
- Configure archive formats (tar.gz, zip)
- Configure package manager integrations (Homebrew, Chocolatey)
- **Acceptance**: GoReleaser config validates and produces expected archives

**T012**: Create GitHub Actions release workflow
- Add .github/workflows/release.yml
- Configure goreleaser-action integration
- Set up artifact uploads to GitHub Releases
- **Acceptance**: Workflow runs successfully and creates releases

**T013**: Implement installation scripts
- Create macOS installation script (install.sh)
- Create Linux installation script (install.sh)
- Create Windows installation script (install.ps1)
- Verify scripts work on target platforms
- **Acceptance**: Installation scripts download and install CLI correctly

**T014**: Test GitHub Releases distribution
- Create a release with goreleaser
- Download and test binary on macOS, Linux, Windows
- Verify CLI executes and works correctly
- **Acceptance**: Users can install CLI from GitHub releases

### Phase 4: User Story 3 - Self-Built Binaries (P2)

Enable users to build and install CLI from source.

```bash
# Filter tasks
bd list --label "spec:003-cli-unification" --label "story:US3" -n 5
bd ready --label "spec:003-cli-unification" --label "story:US3" -n 5
```

**T015**: Update Makefile for CLI build
- Add `make build` target
- Add cross-platform build targets (linux, darwin, windows)
- Configure output to bin/sl
- **Acceptance**: `make build` produces bin/sl binary

**T016**: Build and verify CLI from source
- Test `make build` on target platforms
- Verify built binary executes correctly
- **Acceptance**: Source build produces working binary

### Phase 5: User Story 4 - Self-Hosted / Local Binaries (P2)

Enable users to run CLI from non-PATH locations.

```bash
# Filter tasks
bd list --label "spec:003-cli-unification" --label "story:US4" -n 5
bd ready --label "spec:003-cli-unification" --label "story:US4" -n 5
```

**T017**: Verify self-hosted binary execution
- Test binary execution from non-PATH locations
- Test binary in project directory
- Test binary with `--help` and `--version` flags
- **Acceptance**: Binary works from any location

### Phase 6: User Story 5 - UVX Style (P2)

Enable users to execute CLI without setup.

```bash
# Filter tasks
bd list --label "spec:003-cli-unification" --label "story:US5" -n 5
bd ready --label "spec:003-cli-unification" --label "story:US5" -n 5
```

**T018**: Prepare standalone executable distribution
- Verify binary is self-contained (no runtime dependencies)
- Create distribution manifest
- **Acceptance**: Standalone executable works without setup

### Phase 7: User Story 6 - Package Manager Integration (P3)

Enable users to install via package managers.

```bash
# Filter tasks
bd list --label "spec:003-cli-unification" --label "story:US6" -n 5
bd ready --label "spec:003-cli-unification" --label "story:US6" -n 5
```

**T019**: Create Homebrew formula
- Create formula in Homebrew tap repository
- Configure formula to download from GitHub releases
- Verify installation works
- **Acceptance**: `brew install specledger` works

**T020**: Create npm/npx package
- Create npm package manifest
- Configure bin field to point to CLI binary
- Push to npm registry
- **Acceptance**: `npx @specledger/cli` works

### Phase 8: Polish & Cross-Cutting Concerns

Tasks that improve quality and completeness.

```bash
# Filter tasks
bd list --label "spec:003-cli-unification" --label "phase:polish" -n 5
bd ready --label "spec:003-cli-unification" --label "phase:polish" -n 5
```

**T021**: Write comprehensive README
- Document installation methods
- Document all commands and flags
- Add examples for common use cases
- Add troubleshooting section
- **Acceptance**: README covers all installation methods and commands

**T022**: Create migration guide
- Document transition from old sl script to unified CLI
- List breaking changes (if any)
- Provide migration scripts or commands
- **Acceptance**: Migration guide helps users transition smoothly

**T023**: Run full integration test suite
- Test all commands end-to-end
- Test all distribution channels
- Verify success criteria from spec are met
- **Acceptance**: All tests pass, success criteria validated

## Dependency Graph

```
Phase 1 (Setup) → T001, T002, T003, T004, T005
                 ↓
          ┌───────┴───────┐
          ↓               ↓
     Phase 2 (US1)    Phase 3 (US2)
     (Unified CLI)   (GitHub Releases)
          ↓               ↓
     Phase 4 (US3)    Phase 5 (US4)
     (Self-Built)    (Self-Hosted)
          ↓               ↓
     Phase 6 (US5)    Phase 7 (US6)
     (UVX Style)     (Package Managers)
          ↓               ↓
     Phase 8 (Polish)
```

**Note**: Phases 3-7 are largely independent and can be developed in parallel with Phase 2. The core unified CLI (US1) must complete before distribution features can be validated.

## Parallel Execution Opportunities

**Within Phase 1**:
- T001 (TUI), T002 (dependency registry), T003 (config), T004 (logging), T005 (main.go) can mostly run in parallel, with minor dependencies

**Across Phases 2-7**:
- US3, US4, US5, US6 can be developed in parallel with US1 (once T001-T005 are complete)
- Each distribution channel has minimal dependencies on other stories

## Independent Test Criteria

### US1: Single Unified CLI Tool
- Install unified CLI
- Run `sl new` in new directory → TUI prompts appear
- Run `sl new --ci --project-name test --short-code t` → Project created
- Run `sl deps list` in existing project → Dependencies listed
- Run `sl --version` → Version message shown
- Run `specledger new` → Works as alias

### US2: GitHub Releases
- Download binary from GitHub release
- Execute binary → Help text shows
- Execute `sl new --ci` → Works correctly

### US3: Self-Built Binaries
- Run `make build` → Binary created
- Execute `bin/sl --version` → Version shown
- Execute `bin/sl new --ci` → Works correctly

### US4: Self-Hosted / Local Binaries
- Place binary in project directory
- Run `./sl --help` → Help shown
- Run `./sl new --ci` → Works correctly

### US5: UVX Style
- Execute standalone binary URL
- Help output appears
- `sl new --ci` works

### US6: Package Manager Integration
- Run `brew install specledger` → CLI accessible
- Run `npx @specledger/cli --version` → Version shown

## MVP Scope

**MVP = User Story 1 (P1)** only

This provides the core unification:
- Single CLI binary
- Interactive TUI bootstrap
- Non-interactive bootstrap
- All deps commands working
- Backward compatibility with specledger alias

Once US1 is complete, the following distribution features can be added incrementally:

**Post-MVP (P1) = User Story 2 (GitHub Releases)**
- Users can install CLI from GitHub releases
- Cross-platform binaries available

**Post-MVP (P2) = User Stories 3-5**
- Self-built binaries for development
- Self-hosted binaries for portability
- UVX-style execution for quick testing

**Post-MVP (P3) = User Story 6**
- Package manager integrations
- Broader accessibility

## Implementation Strategy

1. **Start with MVP (US1)**: Implement unified CLI with TUI and all deps commands
2. **Validate core functionality**: Test US1 independently, verify all acceptance scenarios
3. **Add GitHub Releases (US2)**: Set up GoReleaser and CI/CD for distribution
4. **Incremental distribution features**: Add self-built, self-hosted, UVX, and package manager support
5. **Polish and document**: Complete README, migration guide, and final testing

## Beads Task Management

```bash
# View all tasks
bd list --label "spec:003-cli-unification" -n 30

# View tasks by phase
bd list --label "spec:003-cli-unification" --label "phase:setup" -n 10
bd list --label "spec:003-cli-unification" --label "phase:us1" -n 10
bd list --label "spec:003-cli-unification" --label "phase:us2" -n 10

# View tasks ready to work on
bd ready --label "spec:003-cli-unification" -n 5

# Update task status
bd update <task-id> --status in_progress
bd update <task-id> --status completed
```

## Success Criteria Mapping

| SC | Description | Task Coverage |
|----|-------------|---------------|
| SC-001 | Bootstrap in < 3 min | T006, T007 |
| SC-002 | 95% can install from GitHub releases | T011-T014 |
| SC-003 | All deps commands work | T008 |
| SC-004 | 90% bootstrap on first attempt | T006-T010 |
| SC-005 | Run from PATH and non-PATH | T015, T017 |
| SC-006 | Works across macOS, Linux, Windows | T011, T015 |
| SC-007 | CI/CD bootstrap < 1 min | T007 |
| SC-008 | Clear error messages on failure | T009 |

## Notes

- **Tests**: Constitution requires test-first approach, but spec doesn't explicitly request tests. Tests should be added where contract tests and integration tests are defined in contracts/CLI-INTERFACE.md
- **Beads Integration**: Use `bd` commands to track task execution and progress
- **User Stories**: Each story is independently testable, enabling parallel work across stories (after foundational setup)

## Remaining Improvements

### Bootstrap Enhancements (Post-MVP)
- ~~**Directory Selection**: Add ability to select custom project directory in TUI~~ ✅ COMPLETED (2026-02-04)
- **Existing Directory Support**: Allow bootstrapping into existing directories (currently prompts for overwrite)
- **Template Files**: Copy actual SpecLedger project templates during bootstrap
- **Git Initialization**: Initialize git repo and make initial commit
- **Beads Configuration**: Configure `.beads/config.yaml` with project prefix
- **Tool Installation**: Install mise tools and configure shell

### Recent Updates (2026-02-04)
- ✅ Added directory selection step to TUI (6-step flow)
- ✅ TUI now prompts for parent project directory with default from config
- ✅ Confirmation screen shows full project path
- ✅ Non-interactive mode supports `--project-dir` flag

### Command Stubs (TODO)
- **`sl graph`**: Implement dependency graph visualization
- **`sl update`**: Implement CLI self-update functionality
- **`sl vendor`**: Full vendor implementation (partially complete)

### Code Cleanup (Completed)
- Removed old `sl` bash script (replaced by Go binary)
- Removed `init/` directory (old prompt files)
- Removed scratch files (`rehash.md`, old reports)
- Removed unused tool directories (`.perles/`, `.gemini/`, `.conductor/`, `.specify/`)
- Fixed duplicate `update` command issue (now separate: `sl update` vs `sl deps update`)
- Clarified stub commands with TODO messages
