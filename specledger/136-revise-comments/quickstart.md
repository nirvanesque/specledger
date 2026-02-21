# Quickstart: 136-revise-comments

## Prerequisites

- Go 1.24+ installed
- `sl auth login` completed (valid Supabase credentials)
- On the `136-revise-comments` branch

## Build & Run

```bash
# Build the CLI
go build -o sl ./cmd/sl/

# Run revise on current branch
./sl revise

# Run revise on a specific branch
./sl revise 009-feature-name
```

## Development Workflow

### 1. Create the revise package

```bash
mkdir -p pkg/cli/revise
```

New files:
- `pkg/cli/revise/types.go` — Data types (ReviewComment, ProcessedComment, RevisionContext)
- `pkg/cli/revise/client.go` — PostgREST client (fetch comments, resolve)
- `pkg/cli/revise/prompt.go` — Template rendering + token estimation
- `pkg/cli/revise/prompt.tmpl` — Embedded Go template
- `pkg/cli/revise/editor.go` — Editor launch helper
- `pkg/cli/revise/git.go` — Git branch/stash/commit helpers

### 2. Create the command

- `pkg/cli/commands/revise.go` — Cobra command definition + main flow
- Register in `cmd/sl/main.go`: `rootCmd.AddCommand(commands.VarReviseCmd)`

### 3. Extend the launcher

- Add `LaunchWithPrompt(prompt string) error` to `pkg/cli/launcher/launcher.go`

### 4. Add TUI components

- `pkg/cli/tui/revise_select.go` — Multi-select Bubble Tea model for artifact/comment selection

## Testing

```bash
# Unit tests
go test ./pkg/cli/revise/...

# Build and test manually
go build -o sl ./cmd/sl/ && ./sl revise
```

## Key API Endpoints

See [contracts/postgrest-api.md](./contracts/postgrest-api.md) for full API details.

Quick test with curl:
```bash
# Get your token
TOKEN=$(sl auth token)
URL=$(sl auth supabase --url)
KEY=$(sl auth supabase --key)

# Fetch unresolved comments for a spec
curl -s "$URL/rest/v1/review_comments?change_id=eq.<change_id>&is_resolved=eq.false" \
  -H "Authorization: Bearer $TOKEN" \
  -H "apikey: $KEY" | jq
```
