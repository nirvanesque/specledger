# Implementation Plan: Improve SpecLedger Command Prompts

**Branch**: `592-prompt-updates` | **Date**: 2026-02-20 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specledger/592-prompt-updates/spec.md`

## Summary

Update three core SpecLedger command prompts (specify, tasks, implement) to improve dependency handling, issue quality, and Definition of Done verification. Changes apply to both `.claude/commands/` (development) and `pkg/embedded/skills/commands/` (embedded templates).

## Technical Context

**Language/Version**: Markdown (prompt files), Go 1.24+ (embedding system)
**Primary Dependencies**: Existing `sl deps` CLI, `sl issue` CLI commands
**Storage**: N/A (documentation updates only)
**Testing**: Manual verification of prompt behavior, existing integration tests
**Target Platform**: Claude Code CLI (via skill/command system)
**Project Type**: Documentation update (no new source code)
**Performance Goals**: N/A (prompt text changes)
**Constraints**: Must maintain backward compatibility with existing workflow
**Scale/Scope**: 6 prompt files (3 commands × 2 locations)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Verify compliance with principles from `.specledger/memory/constitution.md`:

- [X] **Specification-First**: Spec.md complete with prioritized user stories
- [X] **Test-First**: Test strategy defined (manual verification + existing tests)
- [X] **Code Quality**: N/A (documentation only)
- [X] **UX Consistency**: User flows documented in spec.md acceptance scenarios
- [X] **Performance**: N/A (prompt text changes)
- [X] **Observability**: N/A (prompt text changes)
- [X] **Issue Tracking**: Will create epic during /specledger.tasks

**Complexity Violations**: None identified

## Project Structure

### Documentation (this feature)

```text
specledger/592-prompt-updates/
├── spec.md              # Feature specification (complete)
├── plan.md              # This file
├── research.md          # Phase 0 output
├── quickstart.md        # Phase 1 output
└── tasks.md             # Phase 2 output (/specledger.tasks)
```

### Files to Modify

```text
# Development prompts (local)
.claude/commands/
├── specledger.specify.md    # Add dependency detection section
├── specledger.tasks.md      # Add DoD population, error handling, DoD summary
└── specledger.implement.md  # Add automated DoD verification section

# Embedded prompts (shipped with binary)
pkg/embedded/skills/commands/
├── specledger.specify.md    # Mirror of development version
├── specledger.tasks.md      # Mirror of development version
└── specledger.implement.md  # Mirror of development version
```

**Structure Decision**: Update existing prompt files in both locations. Both sets must remain identical.

## Implementation Approach

### Phase 1: specledger.specify.md Updates

**Dependency Detection Section** (insert after step 4 in Outline):
```markdown
4a. **Dependency Detection** (explicit syntax):
    - Scan user description for patterns: `deps:alias-name` or `@alias-name`
    - For each match, run `sl deps list` to check if dependency exists
    - If exists: Load content from cache (~/.specledger/cache/<alias>/) and include in spec context
    - If not found: Display "Dependency '<alias>' not found. Use 'sl deps add --alias <alias> <source>' to add it."
    - Continue with spec generation after resolving all dependencies
```

### Phase 2: specledger.tasks.md Updates

**Issue Content Structure** (update step 5):
```markdown
Tasks must include:
- title: Concise summary (under 80 chars)
- description: Problem statement (WHY this matters)
- design: Implementation details (HOW/WHERE to build)
- acceptance: Success criteria (WHAT done looks like)
- definition_of_done: Checklist items derived from acceptance criteria
```

**DoD Summary Section** (add to tasks.md template):
```markdown
## Definition of Done Summary

| Issue ID | DoD Items |
|----------|-----------|
| SL-xxxxx | - Item 1\n- Item 2 |
```

**Error Handling** (add after Example CLI Calls):
```markdown
## Error Handling

When `sl issue create` or `sl issue link` fails:
1. Sanitize special characters in description (escape quotes, newlines)
2. Retry with sanitized parameters
3. If still failing: Display clear error with specific issue and remediation
```

### Phase 3: specledger.implement.md Updates

**DoD Verification Section** (add before Completion validation):
```markdown
10a. **Definition of Done Verification** (before closing issues):
     - Read issue's definition_of_done field via `sl issue show <id>`
     - For each DoD item, attempt automated verification:
       - "file exists: <path>" → Check if file exists
       - "test passes: <test>" → Run test command
       - "syntax valid: <file>" → Run linter/syntax check
     - For items that cannot be automated: Prompt user "Is '<item>' complete? (y/n)"
     - If any verification fails: Display failed items with reasons, require --force to proceed
     - Log verification results for audit trail
```

## Complexity Tracking

No complexity violations identified. This is a documentation-only update.

## Acceptance Criteria Mapping

| FR | User Story | Verification Method |
|----|------------|---------------------|
| FR-001 | US1 | Manual: Test `deps:alias` syntax detection |
| FR-002 | US1 | Manual: Verify dep content loads in spec |
| FR-003 | US1 | Manual: Test missing dep prompt |
| FR-004 | US2 | Manual: Verify issue structure |
| FR-005 | US2 | Manual: Check DoD items populated |
| FR-006 | US2 | Manual: Verify DoD summary in tasks.md |
| FR-007 | US2 | Manual: Test error handling |
| FR-008 | US3 | Manual: Test automated DoD verification |
| FR-009 | US3 | Manual: Test interactive fallback |
| FR-010 | US3 | Manual: Verify failed item display |
| FR-011 | All | Diff: Compare .claude vs embedded prompts |
