# Quickstart: Project Template & Coding Agent Selection

**Feature**: 593-init-project-templates
**Date**: 2026-02-20

## Developer Quickstart

### 1. Run Interactive Template Selection

```bash
sl new
```

Follow the prompts:
1. Enter project name
2. Choose project directory
3. Enter short code (2-4 letters)
4. **NEW**: Select from 7 project templates
5. Select constitution principles
6. **NEW**: Choose coding agent (Claude Code, OpenCode, None)
7. Confirm and create

### 2. Non-Interactive Mode

```bash
sl new \
  --project-name my-app \
  --project-dir /path/to/parent \
  --short-code ma \
  --template full-stack \
  --agent claude-code
```

### 3. List Available Templates

```bash
sl new --list-templates
```

Expected output:
```
Available Templates:
  general-purpose      General Purpose Go CLI/Library
  full-stack           Full-Stack Application (Go + React)
  batch-processing     Batch Data Processing (Temporal)
  realtime-workflow    Real-Time Workflow (Temporal)
  ml-image             ML Image Processing (PyTorch)
  realtime-pipeline    Real-Time Data Pipeline (Kafka)
  ai-chatbot           AI Chatbot (LangChain)
```

## Implementation Phases

### Phase 1: Extend Metadata Schema ✅

1. Add `github.com/google/uuid` dependency
2. Extend `ProjectInfo` struct with UUID, template, agent fields
3. Update schema version to 1.1.0
4. Implement v1.0.0 → v1.1.0 migration

**Files Modified**:
- `pkg/cli/metadata/schema.go`
- `pkg/cli/metadata/migration.go` (new)

### Phase 2: Create Template Structures ✅

1. Create 6 new template directories under `pkg/embedded/templates/`
2. Rename `specledger/` to `general-purpose/`
3. Update manifest.yaml with all 7 templates
4. Create README and starter files for each template

**Files Modified**:
- `pkg/embedded/templates/manifest.yaml`
- `pkg/embedded/templates/general-purpose/` (renamed)
- `pkg/embedded/templates/full-stack/` (new)
- `pkg/embedded/templates/batch-processing/` (new)
- `pkg/embedded/templates/realtime-workflow/` (new)
- `pkg/embedded/templates/ml-image/` (new)
- `pkg/embedded/templates/realtime-pipeline/` (new)
- `pkg/embedded/templates/ai-chatbot/` (new)

### Phase 3: Update TUI Flow ✅

1. Add template selection step after `stepShortCode`
2. Add agent selection step (modify existing `stepAgentPreference`)
3. Load templates from manifest
4. Update confirmation view

**Files Modified**:
- `pkg/cli/tui/sl_new.go`

### Phase 4: Implement Bootstrap Logic ✅

1. Generate UUID on project creation
2. Copy template files based on selection
3. Create agent config directories
4. Generate .claude/settings.json with project UUID
5. Write metadata with template and agent fields

**Files Modified**:
- `pkg/cli/commands/bootstrap_helpers.go`
- `pkg/cli/commands/new.go`

### Phase 5: Add CLI Flags ✅

1. Add `--template` flag
2. Add `--agent` flag
3. Add `--list-templates` flag
4. Implement validation and error messages

**Files Modified**:
- `pkg/cli/commands/new.go`

### Phase 6: Testing ✅

1. Unit tests for UUID generation
2. Unit tests for metadata migration
3. Integration tests for all templates
4. Integration tests for all agents
5. Backward compatibility tests

**Files Created**:
- `tests/integration/template_selection_test.go`
- `tests/integration/agent_selection_test.go`
- `tests/integration/uuid_generation_test.go`
- `tests/integration/metadata_migration_test.go`

## Testing Checklist

- [ ] All 7 templates create correct directory structures
- [ ] All 3 agents create correct configuration files
- [ ] Default selections (general-purpose + claude-code) match old behavior
- [ ] UUID generation creates unique IDs (test 10,000 iterations)
- [ ] Metadata migration from v1.0.0 works without errors
- [ ] Non-interactive mode works with all flag combinations
- [ ] --list-templates displays all templates correctly
- [ ] Session capture hooks are generated with correct UUID
- [ ] TTY detection requires flags in non-interactive environments

## References

- Feature Spec: `spec.md`
- Research: `research.md`
- Data Model: `data-model.md`
- Implementation Plan: `plan.md`
