# Implementation Plan: Revise Comments CLI Command

**Branch**: `136-revise-comments` | **Date**: 2026-02-20 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `specledger/136-revise-comments/spec.md`

## Summary

Add `sl revise` as a new Cobra subcommand in the Go CLI binary. The command fetches unresolved artifact comments from Supabase (via PostgREST), presents an interactive TUI for selecting artifacts and processing comments, generates a combined revision prompt from an embedded Go template, launches the user's editor for prompt refinement, spawns the configured coding agent (Claude Code) with the prompt, and offers to commit/push changes and resolve comments afterward.

## Technical Context

**Language/Version**: Go 1.24.2
**Primary Dependencies**: Cobra (CLI), Bubble Tea + Bubbles + Lipgloss + **Huh** (TUI/forms), go-git/v5, `net/http` (PostgREST), `text/template` (prompt rendering), `os/exec` (editor + agent launch)
**Storage**: Supabase PostgreSQL via PostgREST REST API (remote); no local persistence
**Testing**: `go test` with table-driven tests for API client, template rendering, token estimation; integration tests for the full command flow
**Target Platform**: macOS, Linux (CLI binary distributed via GoReleaser + Homebrew)
**Project Type**: Single Go module (existing CLI binary)
**Performance Goals**: Comment fetch <5s (SC-001). SC-002 dropped — interactive/agent steps are user-paced.
**Constraints**: Must preserve TTY when spawning editor and agent; one new dependency (`charmbracelet/huh` for forms) beyond what's in go.mod
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
│   ├── client.go                    # PostgREST client with auto-retry on 401
│   ├── types.go                     # ReviewComment, ProcessedComment, RevisionContext, AutoFixture structs
│   ├── prompt.go                    # Template rendering + token estimation (len/3.5 heuristic)
│   ├── editor.go                    # Editor launch (temp file + $EDITOR)
│   ├── git.go                       # Branch switching, stash, commit/push helpers
│   ├── automation.go                # Fixture file parsing and non-interactive flow
│   └── prompt.tmpl                  # Embedded Go template for revision prompt
├── launcher/
│   └── launcher.go                  # MODIFY: Add LaunchWithPrompt(prompt string) method
└── tui/
    └── revise_forms.go              # NEW: huh-based forms for artifact selection, comment processing, resolution

tests/
├── unit/
│   └── revise/
│       ├── prompt_test.go           # Phase A: Template rendering + token estimation (pure functions)
│       ├── automation_test.go       # Phase A: Fixture parsing + snapshot tests
│       └── client_test.go           # Phase B: PostgREST client tests (httptest mock)
└── integration/
    └── revise_test.go               # Phase B: Build binary, run --auto, verify stdout
```

**Structure Decision**: Single Go module, new `pkg/cli/revise/` package for domain logic, command handler in `pkg/cli/commands/revise.go`. Follows existing patterns (auth, session).

## Implementation Priority

**Phase A — Functionality first** (this sprint):
1. PostgREST client (`client.go`) + types (`types.go`)
2. Command handler (`commands/revise.go`) + registration in `main.go`
3. Git helpers (`git.go`) — branch detection, stash, checkout, commit/push
4. TUI forms (`revise_forms.go`) — huh-based interactive prompts
5. Prompt template (`prompt.go` + `prompt.tmpl`) + token estimation
6. Editor launch (`editor.go`)
7. Agent launcher extension (`LaunchWithPrompt`)
8. Automation mode (`automation.go`) + `--auto` / `--dry-run` flags
9. Lightweight tests: template rendering, token estimation, fixture parsing, snapshot tests via `--auto`

**Phase B — Test infrastructure** (follow-up sprint):
1. Establish `httptest.NewServer` mock pattern for PostgREST
2. PostgREST client unit tests (mock fetch, resolve, 401 retry)
3. Reusable `MockPostgREST` helper (benefits session package too)
4. Integration tests (build binary, run `sl revise --auto`, verify stdout)
5. CI pipeline update for integration tests
6. Regex unit tests for `pkg/cli/git` and `pkg/cli/session` — see TD-3 in Tech Debt section

## Key Design Decisions

### 1. PostgREST Client with Auth Auto-Retry (`pkg/cli/revise/client.go`)

Follow the `session.MetadataClient` pattern with an added `doWithRetry` wrapper:
- `ReviseClient` struct with `baseURL`, `anonKey`, `*http.Client`
- Methods: `GetProject()`, `GetSpec()`, `GetChange()`, `FetchComments()`, `ResolveComment()`, `ListSpecsWithComments()`
- Each method makes one HTTP call, returns typed structs
- **All API calls wrapped with `doWithRetry`**: On 401/PGRST303, refresh the access token via `auth.GetValidAccessToken()` and retry once. This handles token expiry during long agent sessions.
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

### 3. TUI Components — `charmbracelet/huh` Forms

New dependency: `go get github.com/charmbracelet/huh`

All interactive prompts use `huh` form groups (replacing custom Bubble Tea models and `tui.ConfirmPrompt()`):

- **Branch selector**: `huh.NewSelect[string]()` with branch names + unresolved comment counts
- **Artifact multi-select**: `huh.NewMultiSelect[string]()` with options like `spec.md (4 comments)`
- **Comment processing**: `huh.NewSelect[string]()` per comment with Process / Skip / Quit options
- **Guidance input**: `huh.NewText()` (multi-line) when user selects "Process"
- **Commit confirm**: `huh.NewConfirm()` for commit/push decision
- **Resolution multi-select**: `huh.NewMultiSelect[string]()` for selecting comments to resolve
- **Stash confirm**: `huh.NewSelect[string]()` with Stash / Abort / Continue options

`huh` forms support forward/back navigation between groups and embed natively in Bubble Tea programs.

**Future sprint**: Multi-pane TUI with artifact tree (left), comment detail (top-right), context (middle), controls (bottom). See research.md R11 for design notes and reference implementations.

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
6. Track every option proposed to user and their choices in a dedicated revision log, create if it doesn't exist yet.

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
- Apply edits incrementally, one comment at a time across all impacted artifacts
- After all edits, summarize what was changed
- Do NOT modify files beyond what the comments request

Begin processing first document and comment.
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

### 6. Token Estimation

Simple character-based heuristic — no external library:
```go
func EstimateTokens(text string) int {
    return int(math.Ceil(float64(len(text)) / 3.5))
}
```
Per Anthropic's recommendation: ~3.5 characters per token, ~20% error margin. Zero binary size impact. See research.md R12 for alternatives evaluated.

### 7. Automation Mode (`pkg/cli/revise/automation.go`)

Fixture file format:
```json
{
  "branch": "009-feature-name",
  "comments": [
    {"file_path": "specledger/009-xxx/spec.md", "selected_text": "highlighted text", "guidance": "Fix this"},
    {"file_path": "specledger/009-xxx/plan.md", "selected_text": "another passage", "guidance": ""}
  ]
}
```

Flags:
- `--auto <fixture.json>`: Non-interactive mode. Match comments by `file_path` + `selected_text`, generate prompt, print to **stdout**, exit. No agent launch, no resolution. Deterministic output enables **snapshot testing**.
- `--dry-run`: Interactive mode variant. Goes through the full interactive flow (select, process, edit) but writes prompt to a file instead of launching the agent. No resolution.
- `--summary`: Compact non-interactive listing of unresolved comments to stdout. One line per comment: `file_path:line  "selected_text"  (author)`. On auth failure, exits silently (exit code 1, no stdout). Designed for agent integration (e.g., `/specledger.clarify` prompt calls `sl revise --summary`).

In `--auto` mode, steps 7-8 are fixture-driven and step 12 prints to stdout and exits. All subsequent steps (editor, agent, commit, resolve) are skipped.

### 9. Post-Agent: No Changes on Disk

When the agent exits with no uncommitted file changes (agent committed itself, or no changes were needed), skip the commit/push step and proceed directly to comment resolution. The user should always be prompted to resolve comments regardless of whether files changed.

### 8. Full Command Flow

```
1. Parse flags (--auto, --dry-run, --summary)
2. Auth check (auth.GetValidAccessToken)
2a. If --summary → resolve branch → fetch comments → print compact listing → exit (silent exit code 1 on auth failure)
3. Branch selection (resolve spec_key; from fixture.branch if --auto)
4. Branch checkout if needed (stash → checkout)
5. Fetch comments (project → spec → change → review_comments) [auto-retry on 401]
6. Fast exit if no unresolved comments
7. Artifact multi-select (huh form) — or fixture-driven if --auto
8. Comment processing loop (huh form per comment) — or fixture-driven if --auto
9. If no comments processed → exit
10. Generate combined prompt (render template)
11. Show token estimate + warnings (len/3.5 heuristic)
12. If --auto → print prompt to stdout → exit (enables snapshot testing)
13. Open editor (temp file → $EDITOR → read back)
14. Confirm/re-edit/cancel prompt (huh confirm)
15. If --dry-run → prompt for filename to write to → exit
16. Launch agent OR prompt for filename to write to
17. Post-agent: check git status
18. If changes on disk → show summary + offer commit/push (huh confirm)
18a. If no changes on disk → skip commit/push, proceed to resolve
19. Refresh auth token (auto-retry handles this transparently)
20. Comment resolution multi-select (huh form)
21. Resolve selected comments (PATCH API) [auto-retry on 401]
22. Session end (stash pop reminder if applicable)
```

## Future Enhancements

Items identified during review feedback, deferred from this sprint:

- **Agent-driven comment management — `sl comments` command group** *(high impact)*: Instead of `sl revise` handling comment resolution after the agent exits, a dedicated `sl comments` command group gives both users and agents fine-grained control over comment workflows:

  **Planned subcommands**:
  - `sl comments pull` — Fetch and cache unresolved comments for the current spec/branch
  - `sl comments list [--artifact <path>]` — List unresolved comments, optionally filtered by artifact
  - `sl comments summarize [--artifact <path>]` — Summarize comments grouped by artifact with counts
  - `sl comments resolve <id-or-match> --reason "..."` — Resolve a comment with an explanation (matches by file_path + selected_text, no UUIDs exposed)
  - `sl comments reply <id-or-match> --body "..."` — Add a reply to a comment without resolving

  **Agentic workflow**: The agent can call these commands from within its shell session to pull context, resolve comments incrementally, and provide structured reasons — eliminating the need for `sl revise` to handle post-agent resolution.

  **Impact on current flow**: Steps 17-22 (post-agent commit/push/resolve in `sl revise`) become a **fallback** for when `sl comments` commands are not available or the agent doesn't use them.

  **Estimated effort**: 5-8 days (command group scaffolding + individual subcommands + prompt template updates)
  **Prerequisite**: Core `sl revise` flow (this sprint) must be working first

- **Export resolve file on auth expiry**: When the refresh token is also expired during the resolve step, export a JSON file listing comments to resolve. The user can re-authenticate and run `sl revise --auto resolve-file.json` to complete resolution.
- **Multi-pane TUI**: Rich TUI with artifact tree (left), comment detail (right), controls (bottom), and free navigation between views. See research.md R11 for design notes and reference implementations.

---

## Tech Debt Identified — Sprint 136 (2026-02)

Items discovered during implementation and manual testing of the `sl revise` flow. Not blocking the initial release but should be addressed in follow-up sprints.

### TD-1 — Editor launch UX: confirm before opening, name the editor

**Problem**: The editor opens immediately after comment processing with no transition message. The user has no indication which editor will be used or what to do when they exit.

**Current behaviour**:
```
[comment processing loop ends]
[vim opens immediately with no preamble]
[user saves and exits]
What would you like to do? › Launch / Re-edit / Write to file / Cancel
```

**Expected behaviour**:
```
──────────────────────────────────────────────────────────
3 comments processed. Revision prompt is ready.
Opening vim — review and refine the prompt, then save and exit to continue.
──────────────────────────────────────────────────────────
[vim opens]
[user saves and exits]
What would you like to do? › Launch / Re-edit / Write to file / Cancel
```

**Implementation**:
- In `editAndConfirmPrompt`, before calling `revise.EditPrompt`: call `detectEditor()` (already exists in `editor.go`), print the transition message naming the editor.
- The `detectEditor()` call currently happens inside `EditPrompt`; expose it or duplicate the lookup.
- No new dependencies.

---

### TD-2 — Model selection: let user choose Claude model before agent launch

**Problem**: The agent is always launched with its default model. For spec revision work, the user may want to pick a specific model (e.g., `claude-opus-4-5` for deep reasoning vs `claude-haiku-4-5` for speed).

**Investigation needed**:
- Does `claude` CLI support `--list-models` or equivalent to enumerate available models at runtime?
- Does `claude` CLI support a `--model <id>` flag that can be prepended to the prompt argument?
- If yes to both: fetch model list, present `huh.NewSelect`, pass `--model <id>` when constructing the `exec.Command` in `AgentLauncher.LaunchWithPrompt`.
- If `--list-models` is not available: hardcode a curated list of current Claude models (updated with each release) as a fallback.

**Scope**: Change is localised to `AgentLauncher.LaunchWithPrompt` in `pkg/cli/launcher/launcher.go` and the call site in `commands/revise.go` Step 9.

---

### TD-3 — Regex unit tests: `pkg/cli/git` and `pkg/cli/metadata` gaps

**Problem**: Three packages use compiled regexes with no test coverage. A regression in `repoURLRe` (URL parsing) was only caught by a manual `sl revise` run — it was not caught by any automated test.

**Packages and patterns without tests**:

| File | Pattern | Risk |
|------|---------|------|
| `pkg/cli/git/git.go` | `repoURLRe` — parses GitHub remote URL to extract owner/repo | High — used on every `sl revise` invocation |
| `pkg/cli/git/git.go` | `featureBranchRe` — detects feature branches for auto-confirm | Medium — affects branch selection UX |
| `pkg/cli/session/capture.go` | `gitCommitPattern`, `gitAmendPattern` — detects git commits in shell history | Medium — missed detection = missed session checkpoint |
| `pkg/cli/commands/bootstrap_helpers.go` | `placeholderPattern` — detects unfilled template slots | Low — bootstrap-only path |

**Recommended tests** (all table-driven):
- `TestGetRepoOwnerName`: standard SSH, SSH config alias (`github.com-<suffix>`), HTTPS, `.git` suffix / no suffix, invalid URL
- `TestIsFeatureBranch`: `136-name` → true, `main` → false, `99-short` → false (< 3 digits), `1234-long` → true
- `TestGitCommitPattern`: `git commit -m "..."` → match, `git commit --amend` → match, `git push` → no match
- `TestPlaceholderPattern`: `[FOO]` → match, `[AB]` → no match (< 3 chars), `[foo]` → no match (lowercase)

**Where to add**: New file `pkg/cli/git/git_test.go`; extend `pkg/cli/session/capture_test.go` (create if needed).

## Complexity Tracking

No violations. The feature follows existing patterns throughout:
- PostgREST client follows `session.MetadataClient`
- Command follows `commands/session.go` structure
- TUI follows `tui/sl_new.go` Bubble Tea patterns
- Agent launch follows `launcher/launcher.go` patterns
