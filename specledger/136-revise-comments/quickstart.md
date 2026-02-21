# sl revise — Quickstart Guide

This quickstart shows how the user stories translate into CLI commands and interactive flows. It includes sample invocations, expected output, and the full workflow for processing review comments with an LLM agent.

## Prerequisites

```bash
# Authenticate (one-time)
sl auth login

# Verify
sl auth status
# Status: Signed in
# Email:  you@example.com
# Token:  Valid (expires in 29m0s)
```

---

## Interactive Workflow

> **Integration Test Target: MVP** (US1, US2, US3)
>
> Steps 1-3 below form the core integration test suite covering branch selection,
> comment fetching, artifact selection, and interactive comment processing.

### 1. Start a Revise Session (US1)

**Scenario 1a: On a feature branch (most common)**

```bash
# Already on branch 009-feature-name
sl revise

# Output:
# ✓ Authenticated as you@example.com
#
# Branch: 009-feature-name
# Continue with this branch? [Y/n]: y
#
# Fetching comments for 009-feature-name...
# ✓ Found 6 unresolved comments across 3 artifacts
```

**Scenario 1b: Specify branch directly**

```bash
sl revise 009-feature-name

# Output:
# ✓ Authenticated as you@example.com
# Fetching comments for 009-feature-name...
# ✓ Found 6 unresolved comments across 3 artifacts
```

> **Note:** If the specified branch differs from your current branch, stash handling
> applies (see [Branch Switching (US7)](#branch-switching-us7) below).

**Scenario 1c: On main — pick from branches with comments**

```bash
sl revise

# Output:
# ✓ Authenticated as you@example.com
#
# You are on main. Select a branch to revise:
#
#   › 009-feature-name         (6 comments)
#     006-oss-public-dashboard (8 comments)
#     011-streamline-onboarding (3 comments)
#
# Use ↑↓ to navigate, Enter to select
```

**Scenario 1d: Not authenticated**

```bash
sl revise

# ✗ Not authenticated. Run `sl auth login` first.
```

**Scenario 1e: No unresolved comments**

```bash
sl revise

# ✓ Authenticated as you@example.com
# Branch: 009-feature-name
# Fetching comments...
# ✓ No unresolved comments found. Nothing to do.
```

---

### 2. Select Artifacts (US2)

After comments are fetched, only artifacts with unresolved comments are shown:

```bash
# Select artifacts to process (toggle with Space, confirm with Enter):
#
#   [x] specledger/009-feature-name/spec.md       (4 comments)
#   [ ] specledger/009-feature-name/plan.md        (1 comment)
#   [x] specledger/009-feature-name/data-model.md  (1 comment)
#
# 2 artifacts selected, 5 comments to process
```

If all comments are already resolved:

```bash
# ✓ All comments are resolved. Nothing to process.
```

---

### 3. Process Comments (US3)

For each comment on the selected artifacts, you choose an action:

```bash
# ═══════════════════════════════════════════════════════════════════════════════
# Comment 1 of 5 | spec.md
# ═══════════════════════════════════════════════════════════════════════════════
# Author: so0k | 2026-02-19 12:42
# ───────────────────────────────────────────────────────────────────────────────
# Selected: "when artifact content fails to load? Display an error message
#            with a retry option"
# ───────────────────────────────────────────────────────────────────────────────
# Feedback: "this is unclear, artifact content is statically pre-generated
#            in place - there is no retry"
# ───────────────────────────────────────────────────────────────────────────────
#
# Action? [p]rocess / [s]kip / [q]uit: p
# Guidance (optional, press Enter to skip): Remove retry language, clarify static generation
#
# ✓ Queued for processing with guidance
```

**Skip a comment:**

```bash
# ═══════════════════════════════════════════════════════════════════════════════
# Comment 2 of 5 | spec.md
# ═══════════════════════════════════════════════════════════════════════════════
# Author: so0k | 2026-02-19 12:39
# ───────────────────────────────────────────────────────────────────────────────
# Selected: "and existing comments"
# ───────────────────────────────────────────────────────────────────────────────
# Feedback: "this is correct, these are beads comments in JSONL"
# ───────────────────────────────────────────────────────────────────────────────
#
# Action? [p]rocess / [s]kip / [q]uit: s
#
# ⊘ Skipped
```

**Quit early (proceeds to prompt generation with comments processed so far):**

```bash
# Action? [p]rocess / [s]kip / [q]uit: q
#
# Exiting loop. 3 comments processed, 2 remaining.
# Proceeding to prompt generation with 3 processed comments...
```

**Edge case — selected text not found in file:**

```bash
# ═══════════════════════════════════════════════════════════════════════════════
# Comment 3 of 5 | spec.md
# ═══════════════════════════════════════════════════════════════════════════════
# Author: Ngoc Tran | 2026-02-19 12:37
# ───────────────────────────────────────────────────────────────────────────────
# ⚠ Original selected text not found in current file version
# Selected: "some text that was since removed"
# ───────────────────────────────────────────────────────────────────────────────
# Feedback: "this is not good"
# ───────────────────────────────────────────────────────────────────────────────
```

---

> **Integration Test Target: MVP Phase 2** (US4, US5)
>
> Steps 4-5 cover prompt generation, editor launch, and agent execution.

### 4. Generate and Edit Revision Prompt (US4)

After the processing loop, a combined prompt is generated:

```bash
# ───────────────────────────────────────────────────────────────────────────────
# Revision Summary
# ───────────────────────────────────────────────────────────────────────────────
# Processed: 3 comments
# Skipped:   2 comments
#
# Estimated tokens: ~1,240 (within recommended range)
#
# Edit prompt before launching agent? [Y/n]: y
```

Your configured editor opens with the rendered prompt:

```bash
# Launching $EDITOR (vim)...
# [Editor opens with revision prompt content]
# [User reviews, modifies if needed, saves and closes]
```

After the editor closes:

```bash
# ✓ Prompt updated (1,180 tokens)
#
# What would you like to do?
#   › Launch Claude Code with this prompt
#     Re-edit prompt
#     Write prompt to file
#     Cancel
```

**Token warnings:**

```bash
# ⚠ Estimated tokens: ~45 (very short — may lack sufficient context for the agent)
```

```bash
# ⚠ Estimated tokens: ~12,400 (very long — may reduce agent effectiveness)
```

**Editor not found:**

```bash
# ⚠ No editor found ($EDITOR/$VISUAL not set, vi not available)
# Prompt displayed below. Enter a filename to save it:
# Filename: revision-prompt.md
# ✓ Prompt saved to revision-prompt.md
```

---

### 5. Launch Coding Agent (US5)

```bash
# Launching Claude Code with prompt...
#
# [Claude Code starts as interactive subprocess]
# [Agent reads artifacts, generates edit suggestions]
# [Agent uses AskUserQuestion for each comment]
# [User approves/rejects edits interactively]
# [Agent exits when done]
```

**Agent not configured:**

```bash
# ⚠ No coding agent configured.
#   Configure one in specledger.yaml or set SPECLEDGER_AGENT environment variable.
#   Install Claude Code: npm install -g @anthropic-ai/claude-code
#
# Enter a filename to write the prompt to: revision-prompt.md
# ✓ Prompt saved to revision-prompt.md
```

---

> **Integration Test Target: Phase 3** (US6)
>
> Steps 6-7 cover the post-agent commit/push and comment resolution flow.

### 6. Commit and Push Changes (US6)

After the agent exits, you see a summary of changes:

```bash
# ───────────────────────────────────────────────────────────────────────────────
# Agent session complete. Changed files:
# ───────────────────────────────────────────────────────────────────────────────
#   M specledger/009-feature-name/spec.md
#   M specledger/009-feature-name/data-model.md
#
# Commit and push these changes? [Y/n]: y
# Commit message: Address review feedback on spec and data model
#
# ✓ Committed: a1b2c3d "Address review feedback on spec and data model"
# ✓ Pushed to origin/009-feature-name
```

**No changes on disk (agent committed itself or no changes needed):**

```bash
# ───────────────────────────────────────────────────────────────────────────────
# Agent session complete. No uncommitted changes detected.
# ───────────────────────────────────────────────────────────────────────────────
# Proceeding to comment resolution...
```

**Skip committing:**

```bash
# Commit and push these changes? [Y/n]: n
#
# ⚠ Changes not committed. Resolving comments without pushing may lead
#   to inconsistencies on the remote.
#
# Proceed to resolve comments anyway? [y/N]: n
# Unresolved comments remain. Re-run `sl revise` after pushing to resolve them.
```

---

### 7. Resolve Comments (US6)

After commit/push (or if skipped), select which comments to resolve:

```bash
# ───────────────────────────────────────────────────────────────────────────────
# Mark comments as resolved
# ───────────────────────────────────────────────────────────────────────────────
# Select comments to resolve (toggle with Space, confirm with Enter):
#
#   [x] spec.md: "this is unclear, artifact content is statically..."
#   [x] spec.md: "replace: ...by making \"public\" project..."
#   [ ] data-model.md: "this is not good"
#
# 2 of 3 selected

# Resolve 2 comments? [Y/n]: y
# ✓ Resolved 2 comments
# ⊘ 1 comment left unresolved
```

**Session end:**

```bash
# ───────────────────────────────────────────────────────────────────────────────
# ✓ Revise session complete
# ───────────────────────────────────────────────────────────────────────────────
# Branch:     009-feature-name
# Processed:  3 comments
# Resolved:   2 comments
# Files:      2 modified
# Commit:     a1b2c3d
#
# 1 unresolved comment remains. Re-run `sl revise` to continue.
```

---

## Branch Switching (US7)

When you select a branch different from your current one:

**No uncommitted changes:**

```bash
# Switching to branch 009-feature-name...
# ✓ Checked out 009-feature-name
```

**Uncommitted changes detected:**

```bash
# ⚠ You have uncommitted changes:
#   M  pkg/cli/commands/auth.go
#   M  pkg/cli/session/capture.go
#
# What would you like to do?
#   › Stash changes and switch
#     Abort (stay on current branch)
#     Continue anyway (risk conflicts)
#
# ✓ Changes stashed
# ✓ Checked out 009-feature-name
```

**At session end (if stashed):**

```bash
# ⚠ You have stashed changes. Run `git stash pop` to restore them.
```

**Remote-only branch:**

```bash
# Branch 042-new-feature exists on remote but not locally.
# Fetching and checking out...
# ✓ Created local branch 042-new-feature tracking origin/042-new-feature
```

---

## Automation Mode (US8)

For CI pipelines, testing, and non-interactive use.

### Fixture File Format

```json
{
  "branch": "009-feature-name",
  "comments": [
    {
      "file_path": "specledger/009-feature-name/spec.md",
      "selected_text": "when artifact content fails to load",
      "guidance": "Remove retry language — content is statically pre-generated"
    },
    {
      "file_path": "specledger/009-feature-name/spec.md",
      "selected_text": "by making project",
      "guidance": "Replace with: by making \"public\" project"
    }
  ]
}
```

### Generate Prompt (stdout)

```bash
sl revise --auto fixture.json

# Outputs the rendered revision prompt to stdout (no agent launched, no
# comments resolved). Deterministic output for snapshot testing.
```

### Snapshot Testing

```bash
# Generate and save expected output
sl revise --auto fixture.json > expected-prompt.txt

# After code changes, verify prompt hasn't regressed
sl revise --auto fixture.json > actual-prompt.txt
diff expected-prompt.txt actual-prompt.txt
```

### Fixture Warnings

```bash
sl revise --auto fixture.json

# ⚠ Comment not found (skipped): file_path="spec.md", selected_text="removed text"
# ⚠ Comment already resolved (skipped): file_path="spec.md", selected_text="old feedback"
# ✓ Matched 2 of 4 fixture comments
# [prompt output follows on stdout]
```

---

## Summary Mode (US9)

Compact, non-interactive listing of unresolved comments for use by other tools (e.g., `/specledger.clarify`).

```bash
sl revise --summary

# Output (stdout):
# specledger/009-feature-name/spec.md:10-15  "Consider adding more detail"    (so0k)
# specledger/009-feature-name/spec.md:42     "Missing performance reqs"       (reviewer)
# specledger/009-feature-name/plan.md:—      "Needs architecture diagram"     (Ariel)
#
# 3 unresolved comments across 2 artifacts
```

**Auth failure (silent exit for agent integration):**

```bash
sl revise --summary
# [no output, exit code 1]
# Calling agent (e.g., /specledger.clarify) can gracefully fall back to local-only analysis
```

### Integration with `/specledger.clarify`

When a user runs `/specledger.clarify` inside Claude Code, the clarify prompt instructs the agent to:

1. Run `sl revise --summary` to fetch reviewer feedback
2. If successful, present comments via AskUserQuestion multi-select
3. Incorporate selected reviewer feedback into the clarification session
4. If auth fails, proceed with local-only spec analysis (no error shown)

```bash
# User inside Claude Code:
/specledger.clarify

# Agent executes internally:
# $ sl revise --summary
# [gets comment listing]
#
# Agent presents to user:
# "I found 3 reviewer comments. Select which to include in this clarification:"
#   [x] spec.md:10-15 — "Consider adding more detail" (so0k)
#   [ ] spec.md:42    — "Missing performance reqs" (reviewer)
#   [x] plan.md       — "Needs architecture diagram" (Ariel)
#
# Agent proceeds with its spec analysis, incorporating the 2 selected comments
# as additional context for ambiguity resolution.
```

---

## Dry Run (Interactive)

Go through the full interactive flow but write the prompt to a file instead of launching the agent:

```bash
sl revise --dry-run

# [Normal interactive flow: branch selection, artifact selection, comment processing]
# ...
#
# Enter a filename to save the prompt: revision-prompt.md
# ✓ Prompt saved to revision-prompt.md (1,240 tokens)
# No agent launched. No comments resolved.
```

---

## Error Handling

**Token expired mid-session:**

```bash
# [Transparent — auto-refreshes on 401 and retries the request]
# No user action needed unless refresh token is also expired:
#
# ✗ Session expired. Run `sl auth login` to re-authenticate.
```

**Network error:**

```bash
# ✗ Failed to fetch comments: connection refused
#   Check your network connection and try again.
```

**File referenced by comment no longer exists:**

```bash
# ═══════════════════════════════════════════════════════════════════════════════
# Comment 4 of 5 | old-spec.md
# ═══════════════════════════════════════════════════════════════════════════════
# ⚠ File not found locally: specledger/009-feature-name/old-spec.md
# ...
```

---

## Command Reference

```
sl revise [branch-name] [flags]

Arguments:
  branch-name    Target branch (optional; auto-detected from current branch)

Flags:
  --auto <file>  Non-interactive mode with fixture file (outputs prompt to stdout)
  --dry-run      Interactive flow but write prompt to file instead of launching agent
  --summary      Compact comment listing to stdout (non-interactive, for agent use)

Examples:
  sl revise                          # Auto-detect branch, interactive
  sl revise 009-feature-name         # Specify branch, interactive
  sl revise --dry-run                # Interactive, save prompt to file
  sl revise --auto fixture.json      # Non-interactive, prompt to stdout
  sl revise --summary                # Compact listing for agent integration
```
