# Implementation Plan: Project Template & Coding Agent Selection

**Branch**: `593-init-project-templates` | **Date**: 2026-02-20 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `specledger/593-init-project-templates/spec.md`

**Note**: This template is filled in by the `/specledger.plan` command. See `.specledger/templates/commands/plan.md` for the execution workflow.

## Summary

Enable developers to select from 7 business-defined project templates (General Purpose, Full-Stack, Batch Data Processing, Real-Time Workflow, ML Image Processing, Real-Time Data Pipeline, AI Chatbot) and choose their preferred coding agent (Claude Code, OpenCode, or None) during the interactive `sl new` command. Each project receives a unique UUID v4 for session tracking in Supabase. The system extends the existing Bubble Tea TUI framework and playbook system to support multiple embedded template types with template-specific directory structures and agent configuration files.

## Technical Context

**Language/Version**: Go 1.24.2
**Primary Dependencies**: github.com/charmbracelet/bubbletea v1.3.10 (TUI framework), github.com/charmbracelet/bubbles v0.21.1 (TUI components), github.com/charmbracelet/lipgloss v1.1.0 (styling), github.com/google/uuid v1.6.0 (UUID generation), github.com/spf13/cobra v1.10.2 (CLI framework), gopkg.in/yaml.v3 (metadata serialization)
**Storage**: File-based (embedded templates via //go:embed, project metadata in specledger.yaml)
**Testing**: Go testing package with table-driven tests for template generation, TUI state transitions, and metadata migration
**Target Platform**: Linux, macOS, Windows (CLI tool compiled for all platforms via GoReleaser)
**Project Type**: Single CLI application extending existing cmd/sl/main.go with pkg/ library structure
**Performance Goals**: Template selection UI responds in <100ms, project creation completes in <5 seconds for default template, UUID generation in <1ms
**Constraints**: Binary size increase limited to 15MB for embedded templates, backward compatibility with existing v1.0.0 projects required, zero breaking changes to current `sl new` default behavior
**Scale/Scope**: 7 embedded templates, 3 agent configurations, supports projects from 1 file to 1000+ files per template

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with principles from `.specledger/memory/constitution.md`:

- [x] **Specification-First**: Spec.md complete with 7 prioritized user stories (P1 and P2), all with acceptance scenarios and test criteria
- [x] **Test-First**: Test strategy defined - table-driven tests for template generation, TUI state machines, metadata migration, and UUID generation. Contract tests for template structure validation
- [x] **Code Quality**: Go 1.24.2 with standard formatting (gofmt), existing .golangci.yml for linting, existing CI/CD pipeline via GitHub Actions
- [x] **UX Consistency**: Interactive TUI flows documented with step-by-step acceptance scenarios, backward compatibility with current default behavior ensures zero breaking changes
- [x] **Performance**: Metrics defined - UI response <100ms, project creation <5 seconds, UUID generation <1ms, binary size increase limited to 15MB
- [x] **Observability**: Structured logging strategy - log all template operations (selection, file copies, errors) to stdout using existing pkg/cli/ui package
- [ ] **Issue Tracking**: Issue tracking system not yet available in codebase (feature 591 documents built-in system but `sl issue` command not found)

**Complexity Violations** (if any, justify in Complexity Tracking table below):
- None identified - feature extends existing TUI framework, playbook system, and metadata structures without introducing new architectural patterns

## Project Structure

### Documentation (this feature)

```text
specledger/[###-feature]/
├── plan.md              # This file (/specledger.plan command output)
├── research.md          # Phase 0 output (/specledger.plan command)
├── data-model.md        # Phase 1 output (/specledger.plan command)
├── quickstart.md        # Phase 1 output (/specledger.plan command)
├── contracts/           # Phase 1 output (/specledger.plan command)
└── tasks.md             # Phase 2 output (/specledger.tasks command - NOT created by /specledger.plan)
```

### Source Code (repository root)

```text
cmd/
└── sl/
    └── main.go                           # CLI entry point

pkg/
├── cli/
│   ├── tui/
│   │   ├── sl_new.go                     # Modified: Add template selection step
│   │   └── sl_init.go                    # Unmodified
│   ├── playbooks/
│   │   ├── templates.go                  # Modified: Support multiple templates
│   │   ├── embedded.go                   # Modified: Load from new template dirs
│   │   ├── manifest.go                   # Modified: Template metadata
│   │   └── copy.go                       # Unmodified
│   ├── commands/
│   │   ├── new.go                        # Modified: Add --template, --agent flags
│   │   ├── init.go                       # Unmodified
│   │   └── bootstrap_helpers.go          # Modified: Template-specific setup
│   ├── metadata/
│   │   └── metadata.go                   # Modified: Add UUID, template, agent fields
│   ├── launcher/
│   │   └── launcher.go                   # Unmodified (existing agent launcher)
│   └── ui/
│       └── colors.go                     # Unmodified (existing styles)
├── embedded/
│   └── templates/
│       ├── specledger/                   # Existing: General Purpose template
│       ├── full-stack/                   # New: Full-Stack Application template
│       ├── batch-processing/             # New: Batch Data Processing template
│       ├── realtime-workflow/            # New: Real-Time Workflow template
│       ├── ml-image/                     # New: ML Image Processing template
│       ├── realtime-pipeline/            # New: Real-Time Data Pipeline template
│       └── ai-chatbot/                   # New: AI Chatbot template
└── version/
    └── version.go                        # Unmodified

tests/
└── integration/
    ├── template_selection_test.go        # New: TUI template selection tests
    ├── agent_selection_test.go           # New: TUI agent selection tests
    ├── uuid_generation_test.go           # New: UUID collision tests
    ├── metadata_migration_test.go        # New: v1.0.0 → v1.1.0 migration tests
    └── template_structure_test.go        # New: Template validation tests
```

**Structure Decision**: Single CLI project extending existing pkg/ structure. Feature adds:
1. Six new embedded template directories under pkg/embedded/templates/
2. Template selection step in existing TUI flow (pkg/cli/tui/sl_new.go)
3. Template and agent metadata in pkg/cli/metadata/metadata.go
4. Integration tests for new functionality in tests/integration/

Follows existing patterns:
- Embedded templates use //go:embed directive (established in feature 005)
- TUI uses Bubble Tea step-based state machine (established in feature 011)
- Metadata stored in specledger.yaml (established in feature 004)
- Agent launching uses existing launcher package (established in feature 011)

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
