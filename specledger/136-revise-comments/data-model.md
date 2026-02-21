# Data Model: 136-revise-comments

**Date**: 2026-02-20

## Entities

### ReviewComment (remote — `review_comments` table)

Artifact-level review feedback attached to a specific file within a spec change.

| Field | Type | Nullable | Notes |
|-------|------|----------|-------|
| `id` | UUID | No | Primary key, auto-generated |
| `change_id` | UUID | No | FK → `changes.id` |
| `file_path` | text | No | Artifact path (e.g., `specledger/006-xxx/spec.md`) |
| `content` | text | No | Reviewer's feedback text |
| `selected_text` | text | Yes | Text passage the reviewer highlighted |
| `line` | integer | Yes | End line number (often NULL) |
| `start_line` | integer | Yes | Start line number (often NULL) |
| `is_resolved` | boolean | No | Default `false`. Set to `true` to resolve |
| `author_id` | UUID | No | FK → `auth.users.id` |
| `author_name` | text | Yes | Display name |
| `author_email` | text | Yes | Email |
| `parent_comment_id` | UUID | Yes | FK → self (threaded replies). NULL for top-level |
| `created_at` | timestamptz | Yes | Default `now()` |
| `updated_at` | timestamptz | Yes | Default `now()` |

**State transitions**: `is_resolved: false` → `is_resolved: true` (one-way, via PATCH)

**Relationships**:
- `review_comments` → `changes` (via `change_id`)
- `changes` → `specs` (via `spec_id`)
- `specs` → `projects` (via `project_id`)
- `review_comments` → `review_comments` (self-referential via `parent_comment_id`)

### Change (remote — `changes` table)

A branch-based changeset for a spec. Acts as the bridge between comments and specs.

| Field | Type | Notes |
|-------|------|-------|
| `id` | UUID | Primary key |
| `spec_id` | UUID | FK → `specs.id` |
| `head_branch` | text | Feature branch name (matches spec_key) |
| `base_branch` | text | Usually `main` |
| `state` | text | `open`, `merged`, etc. |

### Spec (remote — `specs` table)

A specification tracked by SpecLedger.

| Field | Type | Notes |
|-------|------|-------|
| `id` | UUID | Primary key |
| `project_id` | UUID | FK → `projects.id` |
| `spec_key` | text | Branch/folder name (e.g., `136-revise-comments`) |

### Project (remote — `projects` table)

A GitHub repository tracked by SpecLedger.

| Field | Type | Notes |
|-------|------|-------|
| `id` | UUID | Primary key |
| `repo_owner` | text | GitHub org/user (e.g., `specledger`) |
| `repo_name` | text | Repository name (e.g., `specledger`) |

## Local Data Structures (Go)

### ProcessedComment

In-memory struct representing a comment the user chose to "process" (not skip).

```go
type ProcessedComment struct {
    Comment  ReviewComment  // The fetched comment
    Guidance string         // Optional user-provided guidance text
    Index    int            // Display index (1-based)
}
```

### RevisionContext

Template rendering context for the combined prompt.

```go
type RevisionContext struct {
    SpecKey     string              // e.g., "136-revise-comments"
    Comments    []PromptComment     // Processed comments with context
}

type PromptComment struct {
    Index       int     // 1-based display index
    ID          string  // Comment UUID (internal, for resolution)
    FilePath    string  // Artifact file path
    Target      string  // selected_text or "Line N" or "General"
    Feedback    string  // Comment content
    Guidance    string  // User guidance (optional)
}
```

### AutoFixture

Fixture file structure for non-interactive automation mode.

```go
type AutoFixture struct {
    Branch   string           `json:"branch"`    // Target branch name
    Comments []FixtureComment `json:"comments"`  // Comments to process
}

type FixtureComment struct {
    FilePath     string `json:"file_path"`      // Artifact file path (for matching)
    SelectedText string `json:"selected_text"`  // Text passage (for matching)
    Guidance     string `json:"guidance"`        // Optional guidance for LLM
}
```

Comments are matched against fetched review comments by `file_path` + `selected_text` (not by UUID, since IDs are internal).

## Query Patterns

### Fetch unresolved comments for a spec

```
project(repo_owner, repo_name) → spec(project_id, spec_key) → change(spec_id) → review_comments(change_id, is_resolved=false)
```

4 sequential HTTP calls via PostgREST.

### Resolve a comment

```
PATCH /rest/v1/review_comments?id=eq.{uuid}
Body: {"is_resolved": true}
```

### List specs with unresolved comments (for branch picker)

Fetch all unresolved comments for project, aggregate client-side by `spec_key`.
