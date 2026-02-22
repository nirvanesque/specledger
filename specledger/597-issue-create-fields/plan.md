# Implementation Plan: Issue Create Fields Enhancement

**Branch**: `597-issue-create-fields` | **Date**: 2026-02-22 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specledger/597-issue-create-fields/spec.md`

## Summary

Enhance `sl issue create` and `sl issue update` commands to support structured fields (`--acceptance-criteria`, `--dod`, `--design`, `--notes`) that already exist in the JSONL model. Update the `specledger.tasks` and `specledger.implement` prompts to utilize these fields for better task generation and implementation tracking. Improve task blocking relationship logic for proper dependency trees.

## Technical Context

**Language/Version**: Go 1.24.2
**Primary Dependencies**: Cobra (CLI), YAML v3 (config), JSONL (storage)
**Storage**: File-based JSONL at `specledger/<spec>/issues.jsonl`
**Testing**: Go standard testing (`go test`), coverage via `make test-coverage`
**Target Platform**: Cross-platform CLI (Linux, macOS, Windows)
**Project Type**: Single project (CLI tool)
**Performance Goals**: CLI response time <100ms for all operations
**Constraints**: Backward compatible with existing issue JSONL format
**Scale/Scope**: Single-user CLI, no concurrent access requirements

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **Specification-First**: Spec.md complete with 5 prioritized user stories and clarifications
- [x] **Test-First**: Test strategy defined (unit tests for CLI flags, integration tests for JSONL persistence)
- [x] **Code Quality**: `gofmt`, `go vet`, golangci-lint configured in CI
- [x] **UX Consistency**: User flows documented in acceptance scenarios per user story
- [x] **Performance**: CLI operations <100ms, no performance-sensitive code paths
- [x] **Observability**: Errors returned via CLI stderr, standard Go error handling
- [x] **Issue Tracking**: Epic to be created via `/specledger.tasks`

**Complexity Violations**: None identified

## Project Structure

### Documentation (this feature)

```text
specledger/597-issue-create-fields/
├── spec.md              # Feature specification (complete)
├── plan.md              # This file
├── research.md          # Phase 0 output (minimal - existing code analyzed)
├── data-model.md        # Phase 1 output (Issue model already exists)
├── quickstart.md        # Phase 1 output (test scenarios)
├── contracts/           # Phase 1 output (CLI command contracts)
└── tasks.md             # Phase 2 output (/specledger.tasks command)
```

### Source Code (repository root)

```text
pkg/
├── cli/
│   └── commands/
│       └── issue.go           # CLI flags and command handlers (MODIFY)
├── issues/
│   ├── issue.go               # Issue model (existing - no changes needed)
│   └── store.go               # JSONL persistence (existing - no changes needed)
└── embedded/
    ├── skills/
    │   └── commands/
    │       ├── specledger.tasks.md      # Task generation prompt (MODIFY)
    │       └── specledger.implement.md  # Implementation prompt (MODIFY)
    └── templates/
        └── specledger/
            └── .claude/
                └── commands/
                    ├── specledger.tasks.md      # Template copy (MODIFY)
                    └── specledger.implement.md  # Template copy (MODIFY)

.claude/
└── commands/
    ├── specledger.tasks.md      # Local command (MODIFY)
    └── specledger.implement.md  # Local command (MODIFY)

tests/
└── issues/
    └── issue_test.go           # Existing tests (ADD new test cases)
```

**Structure Decision**: Single project structure. Changes are isolated to:
1. `pkg/cli/commands/issue.go` - CLI flag additions and handlers
2. `.claude/commands/` and `pkg/embedded/` - Prompt template updates

## Complexity Tracking

No violations to justify.

## Implementation Phases

### Phase 1: CLI Flag Additions (US1, US2)

**Files to modify**: `pkg/cli/commands/issue.go`

**New flags for `issueCreateCmd`**:
```go
issueAcceptFlag     string  // --acceptance-criteria
issueDoDFlag        []string // --dod (StringArray for repeated flags)
issueDesignFlag     string  // --design
issueNotesFlag      string  // --notes (already exists for update, add to create)
```

**New flags for `issueUpdateCmd`**:
```go
issueDoDFlag        []string // --dod (replace entire DoD)
issueCheckDoDFlag   string   // --check-dod
issueUncheckDoDFlag string   // --uncheck-dod
```

**Handler changes**:
- `runIssueCreate`: Populate new fields from flags before saving
- `runIssueUpdate`: Handle DoD replacement and check/uncheck operations
- `runIssueShow`: Display acceptance_criteria and design in dedicated sections

### Phase 2: Update Issue Show Display (FR-010)

**Current state**: Shows definition_of_done but not acceptance_criteria or design

**New output format**:
```
Issue: SL-xxxxxx
  Title: Task title
  Type: task
  Status: open
  ...

Description:
  [description text]

Acceptance Criteria:
  [acceptance_criteria text]

Design:
  [design text]

Definition of Done:
  [x] Item 1
  [ ] Item 2

Notes:
  [notes text]
```

### Phase 3: Prompt Template Updates (US4, US5)

**Files to modify**:
- `.claude/commands/specledger.tasks.md`
- `.claude/commands/specledger.implement.md`
- `pkg/embedded/skills/commands/specledger.tasks.md`
- `pkg/embedded/templates/specledger/.claude/commands/specledger.tasks.md`
- `pkg/embedded/templates/specledger/.claude/commands/specledger.implement.md`

**Tasks prompt changes**:
1. Update CLI example section to use `--acceptance-criteria`, `--dod`, `--design` flags
2. Add instruction to populate design field from plan.md for feature issues
3. Update error handling examples for new flag patterns
4. Improve blocking relationship instructions (US3)

**Implement prompt changes**:
1. Add instruction to read design field at task start
2. Add instruction to read acceptance_criteria at task start
3. Add instruction to verify against acceptance_criteria before completion
4. Add instruction to use `--check-dod` when subtasks complete

### Phase 4: Tests

**Unit tests** in `tests/issues/issue_test.go`:
- Test issue create with all 4 new flags
- Test issue create with repeated --dod flags
- Test issue update with --dod replacement
- Test issue update with --check-dod success and error cases
- Test issue update with --uncheck-dod success and error cases
- Test exact text matching for DoD operations

## Dependencies

### Existing (no changes needed)
- `pkg/issues/issue.go`: Issue model already has all required fields
- `pkg/issues/store.go`: JSONL persistence already handles the fields
- `IssueUpdate` struct: Already has DefinitionOfDone, CheckDoDItem, UncheckDoDItem

### New
- Cobra StringArray for `--dod` repeated flag pattern

## Risk Assessment

| Risk | Impact | Mitigation |
|------|--------|------------|
| Breaking existing issue JSONL format | Low | Fields already exist in model, just exposing via CLI |
| Prompt template drift | Medium | Update both .claude/commands and pkg/embedded/ copies |
| DoD text matching edge cases | Low | Clarified: exact match, documented in spec |

## Acceptance Criteria Summary

From spec.md:
- [ ] FR-001 through FR-009: CLI flag support
- [ ] FR-010: Issue show display
- [ ] FR-011 through FR-015: Prompt updates
- [ ] FR-016 through FR-018: Blocking relationship improvements
- [ ] SC-001 through SC-009: All success criteria met
