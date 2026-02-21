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

## R10a: Existing Test Infrastructure

**Decision**: Implement functionality first. Defer PostgREST mocking and test suite to a follow-up phase.

**Rationale**: The codebase has zero HTTP mocking infrastructure. The `session` package — which makes identical PostgREST calls — has no tests. Establishing HTTP mocking patterns is valuable but should not block the feature implementation.

**Current state**:
- **Unit tests**: Standard `testing.T` with table-driven tests, scattered across `pkg/cli/metadata/`, `pkg/cli/prerequisites/`, `pkg/deps/`, `pkg/embedded/`
- **Integration tests**: Build the `sl` binary, run as subprocess, verify output/filesystem (`tests/integration/`)
- **HTTP mocking**: None. No `httptest.NewServer`, no custom `RoundTripper`, no mock generators
- **Session package tests**: None exist despite making PostgREST and Storage API calls
- **Test deps**: `stretchr/testify v1.11.1` available as indirect dep but rarely used
- **CI**: GitHub Actions runs `make test` (unit) and golangci-lint. `make test-integration` exists but not in CI pipeline

**Testing approach for this feature (phased)**:

Phase A (with implementation):
- `--auto` mode with snapshot testing (prompt output to stdout is deterministic)
- Template rendering unit tests (pure function, no HTTP)
- Token estimation unit tests (pure function)
- Fixture parsing unit tests (JSON deserialization)
- Git helper tests (can use `t.TempDir()` with git init)

Phase B (follow-up):
- PostgREST client tests using `httptest.NewServer` to mock Supabase responses
- Auth retry tests (mock 401 → refresh → retry)
- Integration tests (build binary, run `sl revise --auto fixture.json`, verify stdout)
- Establish reusable `MockPostgREST` helper for the session package as well

**Mock pattern to establish**:
```go
func newMockPostgREST(t *testing.T, routes map[string]interface{}) *httptest.Server {
    return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Match route by path + query params
        // Return canned JSON response
    }))
}
```

## R11: TUI Approach — Sequential with `huh` Forms

**Decision**: Use `charmbracelet/huh` form library for sequential flow. Defer multi-pane TUI to a future sprint.

**Rationale**: The existing codebase uses forward-only Bubble Tea flows (`sl_new.go`, `sl_init.go`) with custom form logic. `huh` provides multi-select, text input, confirm, and group-based step navigation out of the box — replacing ~200 lines of custom form code with ~30 lines of `huh` calls. Multi-pane TUI (tree view + detail + controls) would require ~10-15 days vs ~3-5 days for `huh` forms.

**Current TUI state**: Forward-only, no back navigation, custom radio/checkbox components, `bubbles/textinput` for text entry. No `huh` in go.mod yet — needs `go get github.com/charmbracelet/huh`.

**Multi-pane design (deferred)**: For a future sprint, the revise TUI could be refactored into:
- Left pane: Artifact tree (using `tree-bubble` or custom)
- Top-right: Comment position indicator + comment detail
- Middle-right: Comment context (file content around `selected_text`)
- Bottom: Controls (process/skip/quit) + guidance text input
- Full prompt editor view with regenerate warning

Reference implementations: [leg100/pug](https://github.com/leg100/pug) (662 stars), [KevM/bubbleo](https://github.com/KevM/bubbleo) (68 stars, NavStack pattern).

**Alternatives considered**: Full multi-pane TUI (deferred due to effort), `tview` (different framework, not Charm ecosystem), plain Bubble Tea without `huh` (more custom code).

## R12: Token Estimation — Simple Heuristic

**Decision**: Use `len(text) / 3.5` character-based heuristic (rounded up). No external library.

**Rationale**: Anthropic does NOT publish a local tokenizer for Claude 3+ models. Their official recommendation for local estimation is ~3.5 characters per token (~20% error margin). All Go tokenizer libraries (`hupe1980/go-tiktoken`, `tiktoken-go/tokenizer`, `pkoukk/tiktoken-go`) use pre-Claude-3 or OpenAI encodings that are equally approximate. The heuristic has zero binary size impact and zero dependencies.

**Libraries evaluated**:
- `hupe1980/go-tiktoken` (20 stars): Has native `claude` encoding but based on pre-Claude-3 tokenizer. +4-6MB binary.
- `tiktoken-go/tokenizer` (421 stars): Uses `cl100k_base` (OpenAI). +4MB binary. Clean API.
- `pkoukk/tiktoken-go` (885 stars): Most popular. Requires network call or +4MB offline loader.
- Anthropic Token Counting API: Exact but requires network call.

**Implementation**:
```go
func EstimateTokens(text string) int {
    return int(math.Ceil(float64(len(text)) / 3.5))
}
```

## R13: Automation Mode — Fixture File Design

**Decision**: Support `sl revise --auto <fixture.json>` for non-interactive mode.

**Rationale**: Enables CI integration, repeatable testing, and batch processing. The fixture file maps comment identifiers to processing decisions and guidance.

**Fixture format**:
```json
{
  "branch": "009-feature-name",
  "comments": [
    {
      "file_path": "specledger/009-xxx/spec.md",
      "selected_text": "some text the reviewer highlighted",
      "guidance": "Replace with more specific language"
    },
    {
      "file_path": "specledger/009-xxx/plan.md",
      "selected_text": "another passage",
      "guidance": ""
    }
  ]
}
```

**Matching strategy**: Comments are matched by `file_path` + `selected_text` (not by UUID), since UUIDs are internal and not exposed to users. If multiple comments match the same file_path + selected_text pair, all are processed.

**Flags**:
- `--auto <fixture.json>`: Non-interactive mode, skips all TUI prompts
- `--dry-run`: Output generated prompt to stdout/file, don't launch agent or resolve comments (works with or without `--auto`)

## R14: Auth Auto-Retry on 401

**Decision**: Wrap all PostgREST API calls with auto-retry on 401/PGRST303.

**Rationale**: The coding agent session can run for 30+ minutes. Supabase access tokens have a short TTL (~30min). When the agent exits and `sl revise` resumes for the resolve step, the token will likely be expired. The existing `auth.GetValidAccessToken()` function already handles auto-refresh via the refresh token. Wrapping API calls to retry on 401 makes this transparent.

**Implementation**: Add a `doWithRetry` helper to the `ReviseClient`:
```go
func (c *ReviseClient) doWithRetry(req *http.Request) (*http.Response, error) {
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode == 401 || resp.StatusCode == 403 {
        resp.Body.Close()
        // Refresh token and retry
        token, err := auth.GetValidAccessToken()
        if err != nil {
            return nil, fmt.Errorf("re-authentication failed: %w", err)
        }
        req.Header.Set("Authorization", "Bearer "+token)
        return c.client.Do(req)
    }
    return resp, nil
}
```
