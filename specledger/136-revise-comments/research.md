# Research: 136-revise-comments

**Date**: 2026-02-20
**Branch**: `136-revise-comments`

## Prior Work

- **specledger.revise.md**: Existing Claude Code slash command at `.claude/commands/specledger.revise.md`. Uses shell CURL commands to fetch/resolve comments. No Go implementation. References both `comments` (issue comments, beads) and `review_comments` (artifact comments) tables — only `review_comments` is in scope for this feature.
- **008-cli-auth**: The `sl auth` command group is fully implemented. Provides `auth.GetValidAccessToken()` with auto-refresh, `auth.GetSupabaseURL()`, `auth.GetSupabaseAnonKey()`. The revise command will follow the same pattern.
- **009-add-login-and-comment-commands**: Previous iteration on comment-related CLI functionality. Some commands may have been partially implemented.
- **Existing launcher package**: `pkg/cli/launcher/launcher.go` provides `AgentLauncher` with `Launch()` method. Currently launches the agent with no arguments. Needs extension to pass prompt as positional argument (per user-provided sample code).

## R1: Database Schema — `review_comments` vs `comments`

**Decision**: Only query `review_comments` table. The `comments` table holds beads issue-tracker comments (task progress logs) — not artifact review feedback.

**Rationale**: Confirmed via Supabase queries:
- `comments` table FK → `issues.id` (beads issues). Sample data shows status updates like "Phase 7 COMPLETE".
- `review_comments` table FK → `changes.id` → `specs.id` → `projects.id`. Sample data shows artifact-level review feedback with `file_path`, `selected_text`, `content`.

**Alternatives considered**: Querying both tables (per existing revise.md) — rejected because issue comments are beads-internal and not relevant to artifact revision.

## R2: API Query Chain for Fetching Comments

**Decision**: Stepwise PostgREST queries: project → spec → change → review_comments.

**Rationale**: PostgREST resource embedding with nested filtering is unreliable for this join depth. The stepwise approach (4 sequential HTTP calls) is more debuggable and follows the existing `session` package pattern.

**Query chain**:
1. `GET /rest/v1/projects?repo_owner=eq.{owner}&repo_name=eq.{name}&select=id` → `project_id`
2. `GET /rest/v1/specs?project_id=eq.{pid}&spec_key=eq.{key}&select=id` → `spec_id`
3. `GET /rest/v1/changes?spec_id=eq.{sid}&select=id,head_branch,state` → `change_id`
4. `GET /rest/v1/review_comments?change_id=eq.{cid}&is_resolved=eq.false&select=*&order=created_at.asc` → comments

**For branch listing (US1.2)**: Aggregate client-side after fetching all unresolved comments for the project. A server-side RPC function would be more efficient but adds a migration dependency — defer to future optimization.

**Alternatives considered**: Single embedded resource query, Postgres RPC function — rejected for simplicity.

## R3: Comment Resolution via PostgREST

**Decision**: `PATCH /rest/v1/review_comments?id=eq.{uuid}` with body `{"is_resolved": true}`.

**Rationale**: Direct PostgREST PATCH. RLS policies on `review_comments` require authenticated user to be a `project_member` — the same JWT used for fetching will work for resolving. The operation is idempotent.

## R4: Existing Codebase Patterns for New Command

**Decision**: Follow existing command conventions exactly.

**Key patterns**:
- Command defined as `var VarReviseCmd` in `pkg/cli/commands/revise.go`
- Registered in `cmd/sl/main.go` via `rootCmd.AddCommand(commands.VarReviseCmd)`
- Auth check: `token, err := auth.GetValidAccessToken()`
- PostgREST client: Follow `session.MetadataClient` pattern (base URL, anon key, http.Client)
- TUI: Bubble Tea for multi-select, `tui.ConfirmPrompt()` / `tui.SelectPrompt()` for simple prompts
- Output: `ui.PrintSuccess()`, `ui.PrintError()`, `ui.Info()` from `pkg/cli/ui/colors.go`
- Git: `session.GetCurrentBranch(cwd)` (exec-based), `deps.OpenRepository()` (go-git based)

## R5: Launching Agent with Prompt

**Decision**: Extend `AgentLauncher` to accept a prompt string, passed as a positional argument to the agent command.

**Rationale**: Per user-provided sample code, Claude Code does not accept stdin or piped input when launched from Go — the prompt must be a positional argument to preserve TTY interactivity. The existing `Launch()` method has no prompt parameter.

**Implementation approach**: Add `LaunchWithPrompt(prompt string) error` method to `AgentLauncher`:
```go
cmd := exec.Command(l.Command, prompt)
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
```

**Alternatives considered**: Stdin pipe (breaks TTY), temp file argument (unnecessary indirection).

## R6: Embedded Go Template for Revision Prompt

**Decision**: Embed the revision prompt template in Go binary using `//go:embed` directive, render with `text/template`.

**Rationale**: Follows existing embedded templates pattern in `pkg/embedded/`. The template combines artifact context, comment details, user guidance, and processing instructions for the coding agent.

**Template structure** (refined from user-provided sample):
- Document context with artifact file path
- Comments section listing each processed comment with target, feedback, and guidance
- Instructions for the agent (analyze, suggest edits, use AskUserQuestion)

## R7: Token Count Estimation

**Decision**: Simple heuristic: `len(prompt) / 4` as approximate token count.

**Rationale**: The 4-characters-per-token heuristic is a well-known rough approximation for English text. Good enough for warnings about prompt size. More precise tokenizers (tiktoken, etc.) would add a dependency for marginal benefit.

**Thresholds**: Under 100 tokens → "may lack context" warning. Over 8000 tokens → "may reduce agent effectiveness" warning.

## R8: Editor Launch for Prompt Editing

**Decision**: Write prompt to temp file, launch `$EDITOR` / `$VISUAL` / `vi`, read back modified content.

**Rationale**: Standard Unix pattern used by `git commit`, `kubectl edit`, etc. The temp file is created in `os.TempDir()` with a `.md` extension for syntax highlighting.

**Implementation**:
```go
tmpFile, _ := os.CreateTemp("", "sl-revise-*.md")
tmpFile.Write([]byte(prompt))
tmpFile.Close()
editor := os.Getenv("EDITOR")
if editor == "" { editor = os.Getenv("VISUAL") }
if editor == "" { editor = "vi" }
cmd := exec.Command(editor, tmpFile.Name())
cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
cmd.Run()
modified, _ := os.ReadFile(tmpFile.Name())
os.Remove(tmpFile.Name())
```

## R9: Git Operations for Branch Switching

**Decision**: Use `exec.Command("git", ...)` for branch operations (stash, checkout, fetch), consistent with `session` package.

**Rationale**: The `go-git` library in the codebase is used for clone/open operations but shell git is already used for branch detection in `session.GetCurrentBranch()`. Branch switching with stash handling is simpler with shell commands.

**Operations needed**:
- `git status --porcelain` — detect uncommitted changes
- `git stash` — stash changes
- `git fetch origin <branch>` — fetch remote branch
- `git checkout <branch>` / `git checkout -b <branch> origin/<branch>` — switch/create tracking branch
- `git add <files>` + `git commit -m "..."` + `git push origin HEAD` — commit/push flow

## R10: Project ID Resolution from Git Remote

**Decision**: Parse git remote URL to extract `repo_owner` and `repo_name`, then query `projects` table.

**Rationale**: The session package already has `session.GetProjectID(cwd)` which does exactly this. Reuse that function.

**Verified**: `SELECT id FROM projects WHERE repo_owner = 'specledger' AND repo_name = 'specledger'` returns `7109364a-2ebc-451f-b052-2fbe5453459e`.
