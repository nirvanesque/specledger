# Implementation Plan: Revise Comments CLI Command

**Branch**: `136-revise-comments` | **Date**: 2026-02-20 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `specledger/136-revise-comments/spec.md`

## Summary

Add `sl revise` as a new Cobra subcommand in the Go CLI binary. The command fetches unresolved artifact comments from Supabase (via PostgREST), presents an interactive TUI for selecting artifacts and processing comments, generates a combined revision prompt from an embedded Go template, launches the user's editor for prompt refinement, spawns the configured coding agent (Claude Code) with the prompt, and offers to commit/push changes and resolve comments afterward.

## Technical Context

**Language/Version**: Go 1.24.2
**Primary Dependencies**: Cobra (CLI), Bubble Tea + Bubbles + Lipgloss (TUI), go-git/v5, `net/http` (PostgREST), `text/template` (prompt rendering), `os/exec` (editor + agent launch)
**Storage**: Supabase PostgreSQL via PostgREST REST API (remote); no local persistence
**Testing**: `go test` with table-driven tests for API client, template rendering, token estimation; integration tests for the full command flow
**Target Platform**: macOS, Linux (CLI binary distributed via GoReleaser + Homebrew)
**Project Type**: Single Go module (existing CLI binary)
**Performance Goals**: Comment fetch <5s (SC-001), full workflow <10min for 5 comments (SC-002)
**Constraints**: Must preserve TTY when spawning editor and agent; no new external dependencies beyond what's already in go.mod
**Scale/Scope**: Typically 1-20 comments per spec, 1-5 artifacts. Single user CLI.

## Constitution Check

*GATE: Constitution template not filled in for this project. Checking against implied principles.*

- [x] **Specification-First**: Spec.md complete with 7 prioritized user stories, 21 FRs, 6 SCs
- [x] **Test-First**: Test strategy defined — unit tests for API client, template rendering, token estimation; integration tests for command flow
- [x] **Code Quality**: `go vet`, `gofmt` (standard Go tooling, already in project)
- [x] **UX Consistency**: All 7 user stories have Given/When/Then acceptance scenarios
- [x] **Performance**: SC-001 (fetch <5s), SC-002 (full flow <10min)
- [x] **Observability**: Logger package exists (`pkg/cli/logger/`), will use for debug logging of API calls
- [ ] **Issue Tracking**: Beads epic to be created during task generation

**Complexity Violations**: None identified.

## Project Structure

### Documentation (this feature)

```text
specledger/136-revise-comments/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/
│   └── postgrest-api.md # PostgREST API contracts
└── tasks.md             # Phase 2 output (via /specledger.tasks)
```

### Source Code (repository root)

```text
cmd/sl/
└── main.go                          # Add: rootCmd.AddCommand(commands.VarReviseCmd)

pkg/cli/
├── commands/
│   └── revise.go                    # NEW: sl revise command handler
├── revise/                          # NEW: revise package
│   ├── client.go                    # PostgREST client for review_comments
│   ├── types.go                     # ReviewComment, ProcessedComment, RevisionContext structs
│   ├── prompt.go                    # Template rendering + token estimation
│   ├── editor.go                    # Editor launch (temp file + $EDITOR)
│   ├── git.go                       # Branch switching, stash, commit/push helpers
│   └── prompt.tmpl                  # Embedded Go template for revision prompt
├── launcher/
│   └── launcher.go                  # MODIFY: Add LaunchWithPrompt(prompt string) method
└── tui/
    └── revise_select.go             # NEW: Multi-select model for artifact/comment selection

tests/
└── unit/
    └── revise/
        ├── client_test.go           # PostgREST client tests (mock HTTP)
        ├── prompt_test.go           # Template rendering + token estimation tests
        └── git_test.go              # Git helper tests
```

**Structure Decision**: Single Go module, new `pkg/cli/revise/` package for domain logic, command handler in `pkg/cli/commands/revise.go`. Follows existing patterns (auth, session).

## Key Design Decisions

### 1. PostgREST Client (`pkg/cli/revise/client.go`)

Follow the `session.MetadataClient` pattern:
- `ReviseClient` struct with `baseURL`, `anonKey`, `*http.Client`
- Methods: `GetProject()`, `GetSpec()`, `GetChange()`, `FetchComments()`, `ResolveComment()`, `ListSpecsWithComments()`
- Each method makes one HTTP call, returns typed structs
- Error handling: parse PostgREST error responses, map to user-friendly messages

### 2. Branch Selection Flow (`commands/revise.go`)

```
sl revise [optional-branch]
    │
    ├─ Has explicit arg? → Use it directly
    │
    ├─ On feature branch (###-*)? → Prefill, confirm with user
    │     └─ User confirms → Use current branch
    │     └─ User wants different → Show branch list
    │
    └─ On main/other? → Show branch list
          └─ Query: all specs with unresolved comments for project
          └─ Present as selectable list with comment counts
          └─ User selects → Check for uncommitted changes → Stash/checkout
```

### 3. TUI Components

- **Branch selector**: Simple `SelectPrompt` (existing) or Bubble Tea list for filtered branches
- **Artifact multi-select**: New Bubble Tea model — checkboxes with `[x] spec.md (4 comments)` format
- **Comment processing loop**: Sequential display with 3-option prompt (Process / Skip / Quit)
- **Resolution multi-select**: Reuse artifact multi-select model for comment selection
- **Confirm prompts**: Existing `tui.ConfirmPrompt()` for commit/push/stash

### 4. Revision Prompt Template

Embedded via `//go:embed prompt.tmpl` in `pkg/cli/revise/prompt.go`.

```gotemplate
You are assisting with document revision for spec "{{.SpecKey}}".

## Artifacts to Revise
The following files contain comments that need to be addressed:
{{- range .Comments}}
- {{.FilePath}}
{{- end}}

You have full access to read and edit these files in the workspace.

## Comments to Address
For each comment below:
1. Read the target location in the document
2. Analyze the reviewer's feedback
3. Generate 2-3 distinct edit suggestions
4. Use AskUserQuestion to present options and get the user's preference
5. Apply the chosen edit

{{- range .Comments}}

### Comment {{.Index}}
- **File**: {{.FilePath}}
- **Target**: {{if .Target}}"{{.Target}}"{{else}}General feedback{{end}}
- **Feedback**: "{{.Feedback}}"
{{- if .Guidance}}
- **Author Guidance**: "{{.Guidance}}"
{{- end}}
{{- end}}

## Important Instructions
- ALWAYS use AskUserQuestion before making any edit
- Present clear, distinct options for each comment
- Apply edits incrementally, one comment at a time
- After all edits, summarize what was changed
- Do NOT modify files beyond what the comments request

Begin processing Comment 1.
```

### 5. Agent Launch

Extend existing `AgentLauncher.LaunchWithPrompt(prompt string)`:
```go
func (l *AgentLauncher) LaunchWithPrompt(prompt string) error {
    cmd := exec.Command(l.Command, prompt)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Dir = l.Dir
    return cmd.Run()
}
```

Critical: prompt is passed as a positional argument (not stdin) to preserve TTY interactivity.

### 6. Full Command Flow

```
1. Auth check (auth.GetValidAccessToken)
2. Branch selection (resolve spec_key)
3. Branch checkout if needed (stash → checkout)
4. Fetch comments (project → spec → change → review_comments)
5. Fast exit if no unresolved comments
6. Artifact multi-select (show only artifacts with comments + counts)
7. Comment processing loop (process/skip/quit per comment)
8. If no comments processed → exit
9. Generate combined prompt (render template)
10. Show token estimate + warnings
11. Open editor (temp file → $EDITOR → read back)
12. Confirm/re-edit/cancel prompt
13. Launch agent OR write prompt to file
14. Post-agent: show git status summary
15. Offer commit + push
16. Comment resolution multi-select
17. Resolve selected comments (PATCH API)
18. Session end (stash pop reminder if applicable)
```

## Complexity Tracking

No violations. The feature follows existing patterns throughout:
- PostgREST client follows `session.MetadataClient`
- Command follows `commands/session.go` structure
- TUI follows `tui/sl_new.go` Bubble Tea patterns
- Agent launch follows `launcher/launcher.go` patterns
