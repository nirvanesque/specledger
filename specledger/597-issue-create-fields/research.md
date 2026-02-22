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

## No NEEDS CLARIFICATION Items

All technical decisions resolved through code analysis and spec clarifications.
