# Data Model: Issue Create Fields Enhancement

**Feature**: 597-issue-create-fields
**Date**: 2026-02-22

## Overview

This feature exposes existing fields in the Issue model via CLI. No model changes required - only CLI flag additions.

## Entities

### Issue (existing - no changes)

Location: `pkg/issues/issue.go`

```go
type Issue struct {
    // Required fields
    ID          string      `json:"id"`
    Title       string      `json:"title"`
    Description string      `json:"description,omitempty"`
    Status      IssueStatus `json:"status"`
    Priority    int         `json:"priority"`      // 0=highest, 5=lowest
    IssueType   IssueType   `json:"issue_type"`
    SpecContext string      `json:"spec_context"`
    CreatedAt   time.Time   `json:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at"`

    // Optional fields (being exposed via CLI)
    ClosedAt           *time.Time        `json:"closed_at,omitempty"`
    DefinitionOfDone   *DefinitionOfDone `json:"definition_of_done,omitempty"`
    BlockedBy          []string          `json:"blocked_by,omitempty"`  // Issue IDs
    Blocks             []string          `json:"blocks,omitempty"`      // Issue IDs
    Labels             []string          `json:"labels,omitempty"`
    Assignee           string            `json:"assignee,omitempty"`
    Notes              string            `json:"notes,omitempty"`
    Design             string            `json:"design,omitempty"`
    AcceptanceCriteria string            `json:"acceptance_criteria,omitempty"`

    // Migration metadata
    BeadsMigration *BeadsMigration `json:"beads_migration,omitempty"`
}
```

### DefinitionOfDone (existing - no changes)

```go
type DefinitionOfDone struct {
    Items []ChecklistItem `json:"items"`
}

type ChecklistItem struct {
    Item       string     `json:"item"`
    Checked    bool       `json:"checked"`
    VerifiedAt *time.Time `json:"verified_at,omitempty"`
}
```

### IssueUpdate (existing - no changes)

```go
type IssueUpdate struct {
    Title              *string
    Description        *string
    Status             *IssueStatus
    Priority           *int
    IssueType          *IssueType
    Assignee           *string
    Notes              *string
    Design             *string
    AcceptanceCriteria *string
    Labels             *[]string
    AddLabels          []string
    RemoveLabels       []string
    BlockedBy          *[]string
    Blocks             *[]string
    DefinitionOfDone   *DefinitionOfDone  // Replace entire DoD
    CheckDoDItem       string             // Item to mark as checked
    UncheckDoDItem     string             // Item to mark as unchecked
}
```

## Field to Flag Mapping

### Issue Create

| Field | CLI Flag | Type | Notes |
|-------|----------|------|-------|
| AcceptanceCriteria | `--acceptance-criteria` | string | Single value |
| DefinitionOfDone | `--dod` | []string (repeated) | Creates unchecked items |
| Design | `--design` | string | Single value |
| Notes | `--notes` | string | Single value |

### Issue Update

| Field | CLI Flag | Type | Notes |
|-------|----------|------|-------|
| DefinitionOfDone | `--dod` | []string (repeated) | Replaces entire DoD |
| (DoD check) | `--check-dod` | string | Exact match, sets checked=true, verified_at=now |
| (DoD uncheck) | `--uncheck-dod` | string | Exact match, sets checked=false, verified_at=nil |

## Validation Rules

1. **DoD item text matching**: Exact match (case-sensitive, no whitespace normalization)
2. **Error on not found**: `--check-dod` and `--uncheck-dod` return error if item text not found
3. **DoD replacement**: `--dod` on update replaces entire DoD, not additive

## State Transitions

### DefinitionOfDone Item

```
[unchecked] --check-dod--> [checked + verified_at]
[checked] --uncheck-dod--> [unchecked + verified_at=nil]
```

### Issue (no changes to existing transitions)

```
open --(status change)--> in_progress --(status change)--> closed
```

Note: All DoD items being checked does NOT auto-close the issue. Explicit `sl issue close` required.
