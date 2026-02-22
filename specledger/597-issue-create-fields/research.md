# Research: Issue Create Fields Enhancement

**Feature**: 597-issue-create-fields
**Date**: 2026-02-22

## Prior Work

### Related Features

| Feature | Description | Relevance |
|---------|-------------|-----------|
| 591-issue-tracking-upgrade | Issue model with all JSONL fields | Model already supports target fields |
| 595-issue-tree-ready | Dependency tree and ready commands | Blocking/dependency logic exists |

### Existing Code Analysis

**`pkg/issues/issue.go`**:
- `Issue` struct already has: `AcceptanceCriteria`, `DefinitionOfDone`, `Design`, `Notes`
- `IssueUpdate` struct already has: `AcceptanceCriteria`, `DefinitionOfDone`, `CheckDoDItem`, `UncheckDoDItem`
- `DefinitionOfDone` struct with `CheckItem()` and `UncheckItem()` methods already implement exact match

**`pkg/cli/commands/issue.go`**:
- `issueCreateCmd` flags: `--title`, `--description`, `--type`, `--priority`, `--labels`, `--spec`, `--force`
- `issueUpdateCmd` flags: `--title`, `--description`, `--status`, `--priority`, `--assignee`, `--notes`, `--design`, `--acceptance-criteria`, `--add-label`, `--remove-label`
- Missing: `--dod`, `--check-dod`, `--uncheck-dod` on update; `--dod`, `--design`, `--acceptance-criteria`, `--notes` on create

## Decisions

### 1. Cobra StringArray for --dod Flag

**Decision**: Use `StringArrayVar` for repeated `--dod` flags

**Rationale**:
- Standard Cobra pattern for repeated string flags
- Allows natural CLI usage: `--dod "Item 1" --dod "Item 2"`
- Avoids comma-splitting complexity (items can contain commas)

**Implementation**:
```go
var issueDoDFlag []string
issueCreateCmd.Flags().StringArrayVar(&issueDoDFlag, "dod", []string{}, "Definition of Done items (can be repeated)")
```

### 2. Exact Text Matching for DoD Operations

**Decision**: Use exact string match (case-sensitive, no whitespace normalization)

**Rationale**:
- Matches existing `CheckItem()` implementation in `pkg/issues/issue.go`
- Predictable behavior - no surprises from normalization
- Clear error message when item not found

**Implementation**: No changes needed to existing `CheckItem()`/`UncheckItem()` methods

### 3. Error Message Format

**Decision**: Return error with format `"DoD item not found: '<text>'"`

**Rationale**:
- Clear identification of what text was searched
- Single quotes distinguish the search text from error message
- Consistent with CLI error patterns

### 4. Prompt Template Strategy

**Decision**: Update both `.claude/commands/` and `pkg/embedded/` copies

**Rationale**:
- `.claude/commands/` - used by this repository
- `pkg/embedded/` - copied to user projects via `sl init`
- Both must stay in sync to ensure consistent behavior

## Alternatives Considered

| Alternative | Rejected Because |
|-------------|------------------|
| Comma-separated `--dod` flag | Items might contain commas; less intuitive UX |
| Case-insensitive DoD matching | Could match wrong item; unpredictable |
| Auto-add DoD item on --check-dod | Hides typos; violates explicit is better than implicit |

## External System Analysis: Beads

**Source**: https://github.com/steveyegge/beads

### Overview

Beads is an AI-assisted issue tracking system designed for small teams. It uses a single `beads.jsonl` file for storage and integrates with AI agents via Model Context Protocol (MCP).

### Supported Fields

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique identifier (format: `B-xxxxxx`) |
| `title` | string | Issue title |
| `type` | string | Issue type: `task`, `epic`, `bug`, `feature`, `story` |
| `status` | string | Status: `open`, `in_progress`, `closed` |
| `priority` | int | Priority level (0=highest) |
| `created` | timestamp | Creation time |
| `modified` | timestamp | Last modification time |
| `external_ref` | string | External reference (e.g., Jira ticket) |
| `reason` | string | Why this issue exists |
| `assignee` | string | Assigned user |
| `parent` | string | Parent issue ID |
| `children` | []string | Child issue IDs |
| `blocks` | []string | Issues this blocks |
| `blocked_by` | []string | Issues blocking this |
| `related` | []string | Related issues |
| `discovered_from` | string | Source issue that discovered this |

### Dependency Types

- **blocks/blocked_by**: Direct blocking relationships
- **parent/children**: Hierarchical task breakdown
- **related**: Loose associations
- **discovered_from**: Traces issue discovery lineage

### Fields NOT Supported (Gap Analysis)

| Field | SpecLedger Has | Beads Has | Impact |
|-------|----------------|-----------|--------|
| `acceptance_criteria` | ✓ | ✗ | Beads users cannot specify acceptance criteria |
| `definition_of_done` | ✓ | ✗ | Beads lacks checklist-style DoD tracking |
| `design` | ✓ | ✗ | Beads has no design notes field |
| `notes` | ✓ | ✗ | Beads has no implementation notes field |
| `labels` | ✓ | ✗ | Beads lacks tagging/labeling system |
| `spec` | ✓ | ✗ | Beads has no spec context linking |

### Implications for This Feature

1. **Competitive Advantage**: Adding `--acceptance-criteria`, `--dod`, `--design`, `--notes` flags gives SpecLedger richer issue metadata than Beads.

2. **No Migration Concerns**: Since Beads doesn't support these fields, there's no need to consider Beads compatibility in our implementation.

3. **Prompt Template Value**: Utilizing these fields in AI prompts provides structured context that Beads users would have to manually embed in descriptions.

### Conclusion

This feature differentiates SpecLedger from Beads by providing structured fields for acceptance criteria, definition of done, design notes, and implementation notes. The blocking relationship improvements also address a gap where both systems could improve—SpecLedger's `sl issue link` command already supports dependencies, but this feature ensures prompts properly instruct agents on creating and maintaining blocking trees.

## No NEEDS CLARIFICATION Items

All technical decisions resolved through code analysis and spec clarifications.
