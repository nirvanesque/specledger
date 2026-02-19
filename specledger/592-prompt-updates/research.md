# Research: Improve SpecLedger Command Prompts

**Feature**: 592-prompt-updates
**Date**: 2026-02-20

## Prior Work

### Related Features

| Feature | Description | Relevance |
|---------|-------------|-----------|
| 591-issue-tracking-upgrade | Built-in issue tracking with `sl issue` commands | Provides `sl issue create`, `sl issue link`, `sl issue show` commands used by updated prompts |
| 008-fix-sl-deps | SpecLedger dependencies management | Provides `sl deps list`, `sl deps add` commands for dependency resolution |

### Key Findings from Prior Work

1. **Issue Tracking System (591)**:
   - Issues stored as JSONL in `specledger/<spec>/issues.jsonl`
   - Each issue has: id, title, description, design, acceptance, definition_of_done, status, type, priority, labels
   - CLI commands: `sl issue create`, `sl issue link`, `sl issue show`, `sl issue update`, `sl issue close`

2. **Dependencies System (008)**:
   - Dependencies cached in `~/.specledger/cache/<alias>/`
   - Metadata in `specledger.yaml` under `dependencies` key
   - CLI commands: `sl deps add`, `sl deps list`, `sl deps update`

## Decisions

### D1: Dependency Detection Syntax

**Decision**: Use explicit syntax only (`deps:alias` or `@alias`)

**Rationale**:
- Predictable behavior - no false positives from natural language
- Clear user intent - explicit references are unambiguous
- Simple implementation - regex pattern matching sufficient
- Backward compatible - no breaking changes to existing workflow

**Alternatives Considered**:
| Approach | Rejected Because |
|----------|-----------------|
| LLM-based semantic detection | Unpredictable, may miss references or false positive |
| Pattern matching "integrate with X" | Ambiguous, may match non-dependency references |
| Required configuration file | Too much friction for users |

### D2: DoD Population Strategy

**Decision**: Derive DoD items from acceptance criteria in spec.md

**Rationale**:
- Acceptance criteria already defined in spec - no duplicate work
- Traceability from DoD back to requirements
- Consistent quality across all generated issues

**Implementation**:
1. Parse acceptance scenarios from spec.md
2. Convert each "Then" clause to a DoD checklist item
3. Include in `definition_of_done` field when creating issue

### D3: DoD Verification Approach

**Decision**: Automated where possible, interactive fallback

**Rationale**:
- Automation reduces manual verification burden
- Interactive fallback handles subjective criteria
- Clear feedback when verification fails

**Automated Verification Patterns**:
| Pattern | Check |
|---------|-------|
| `file exists: <path>` | `test -f <path>` |
| `directory exists: <path>` | `test -d <path>` |
| `command succeeds: <cmd>` | Execute command, check exit code |
| `tests pass` | Run `go test ./...` or equivalent |

### D4: Error Handling Strategy

**Decision**: Auto-fix and retry with corrected parameters

**Rationale**:
- Reduces friction - most errors are fixable (special chars, etc.)
- Clear feedback when unfixable errors occur
- Maintains workflow continuity

**Common Error Fixes**:
| Error Type | Fix |
|------------|-----|
| Special characters in description | Escape quotes, newlines |
| Missing required field | Use default or prompt |
| Invalid label format | Sanitize to valid format |

## No NEEDS CLARIFICATION Items

All technical decisions resolved through this research. No clarifications required before implementation.

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Prompt changes break existing workflow | Low | High | Manual testing of all three commands |
| Embedded prompts drift from dev prompts | Medium | Medium | Add diff check to CI |
| DoD verification too aggressive | Low | Medium | Allow --force override |
