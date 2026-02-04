# SpecLedger Repository Skills Template

This template provides context about all available skills in the SpecLedger repository. Use this to understand the tooling landscape when working on SpecLedger projects.

## Repository Overview

SpecLedger is a unified CLI for project bootstrap and specification dependency management. It includes:

- **CLI**: `sl` command for all operations
- **Beads**: Graph-based issue tracker for persistent work tracking
- **Specs**: External specification dependency management
- **TUI**: Beautiful terminal UI with gum

## Available Skills

### 1. [deps-manager](.claude/skills/deps-manager/SKILL.md)

**Purpose**: Manage external specification dependencies

**Commands**:
- `sl deps add <repo-url>` - Add a specification dependency
- `sl deps list` - List all declared dependencies
- `sl deps resolve` - Resolve all dependencies and generate spec.sum
- `sl deps update` - Update dependencies to latest compatible versions
- `sl deps remove <repo-url> <spec-path>` - Remove a dependency

**When to use**:
- Adding/Removing dependencies
- Resolving spec dependencies
- Understanding dependency trees

**Quick example**:
```bash
sl deps add github.com/org/project-spec main specs/api.md --alias api
sl deps list
sl deps resolve
```

---

### 2. [bd-issue-tracking](.claude/skills/bd-issue-tracking/SKILL.md)

**Purpose**: Track complex, multi-session work with dependency graphs using bd (beads) issue tracker

**When to use**:
- Multi-session work spanning multiple days
- Complex dependencies with blockers/prerequisites
- Knowledge work with fuzzy boundaries
- Side quests that might pause main task
- Project memory that needs to survive compaction

**Key commands**:
```bash
bd ready --json              # See available work
bd create "Title" -d "Description"  # Create new issue
bd show issue-id             # View issue details
bd update issue-id --status in_progress  # Start working
bd close issue-id --reason "Reason"     # Mark complete
```

**Alternative to TodoWrite**:
- TodoWrite: Short-term (this session)
- bd: Long-term (survives compaction)

---

### 3. [gum-tui](.claude/skills/gum-tui/SKILL.md)

**Purpose**: Use gum for beautiful terminal UI prompts

**When to use**:
- Interactive project setup (`sl new`)
- Any terminal interaction that needs better UX

**Features**:
- `gum input` - Text input with placeholders
- `gum confirm` - Yes/no confirmation
- `gum choose` - Selection menus with color highlighting

**Quick example**:
```bash
sl new  # Uses gum for beautiful TUI
```

---

## Choosing the Right Tool

### deps-manager vs Manual Specs

**Use deps-manager when**:
- You need to reference external specs across multiple projects
- You want reproducible spec resolution
- You need transitive dependency tracking
- You're working on specification-heavy projects

**Use manual specs when**:
- Specs are self-contained and don't need versioning
- You want simple markdown references
- Specs are not shared across projects

### bd-issue-tracking vs TodoWrite

**Use bd when**:
- Work spans multiple sessions
- Dependencies/blockers exist
- Context will be needed after compaction
- Exploratory or fuzzy work

**Use TodoWrite when**:
- Single-session tasks
- Linear execution with no branching
- Immediate context available

**Decision criteria**:
```
❓ Will I need this context in 2 weeks?
   Yes → bd
   No → Continue

❓ Does this have blockers/dependencies?
   Yes → bd
   No → Continue

❓ Is this linear with no branching?
   No → TodoWrite
   Yes → Continue
```

### gum-tui vs Basic CLI

**Use gum-tui when**:
- You want beautiful interactive prompts
- Terminal has proper color support
- User experience matters

**Use basic CLI when**:
- CI/CD or non-interactive environment
- Terminal has limited capabilities
- Simpler is better

---

## Common Workflows

### 1. Project Setup with Dependencies

```bash
# Create project with TUI
sl new

# Add dependencies
sl deps add github.com/org/base-spec main specs/base.md --alias base

# Resolve to fetch
sl deps resolve

# Verify
sl deps list
```

### 2. Feature Development

```bash
# Create task issue
bd create "Implement user authentication" -p 0

# Add dependency if needed
sl deps add github.com/org/auth-spec main specs/auth.md --alias auth

# Start working
bd update <issue-id> --status in_progress

# Code...
bd show <issue-id> --notes "COMPLETED: [notes] IN PROGRESS: [notes] NEXT: [notes]"

# Close when done
bd close <issue-id> --reason "Done"
```

### 3. Spec Development

```bash
# Create spec
bd create "API specification" -d "Define all API endpoints"

# Add external dependencies
sl deps add github.com/org/endpoint-specs main specs/endpoints.md --alias endpoints

# In your spec:
# {{ $endpoints := (lookup "github.com/org/endpoint-specs" "main" "specs/endpoints.yml") }}
# {{ $endpoints }}
```

### 4. Update Workflow

```bash
# Check what's ready
bd ready

# Create/update issues
bd create "Update to latest deps"

# Add dep issues if found
sl deps add github.com/org/new-dep main specs/new.md

# Resolve and update
sl deps resolve
sl deps update
```

---

## Integration Points

### deps-manager + bd-issue-tracking

When discovering issues during dependency work:
```bash
# Create issue for discovered problem
bd create "Found: Invalid spec URL"

# Link issues
bd dep add <current-issue> <dependency-issue> --type discovered-from
```

### deps-manager + gum-tui

Interactive dependency management:
```bash
# Beautiful prompts for project setup
sl new

# Uses gum input/confirm/choose
```

---

## Files and Locations

| File | Purpose | Related Skill |
|------|---------|---------------|
| `specs/spec.mod` | Dependency manifest | deps-manager |
| `specs/spec.sum` | Lockfile with hashes | deps-manager |
| `.specledger/` | Beads database | bd-issue-tracking |
| `.beads/` | Beads database (project-local) | bd-issue-tracking |
| `~/.beads/` | Beads database (global) | bd-issue-tracking |

---

## Getting Help

### deps-manager
- See [deps-manager skill](.claude/skills/deps-manager/SKILL.md) for detailed documentation
- Run `sl deps --help` for command-specific help

### bd-issue-tracking
- See [bd-issue-tracking skill](.claude/skills/bd-issue-tracking/SKILL.md) for detailed documentation
- Run `bd --help` for command-specific help

### gum-tui
- Run `gum --help` for gum-specific help
- Terminal uses gum when available, falls back to basic prompts

---

## Next Steps

1. **Start a new project**: `sl new`
2. **Track your work**: `bd create "First task"`
3. **Manage dependencies**: `sl deps add <repo>`
4. **Resolve specs**: `sl deps resolve`

See individual skills for detailed documentation.

---

**Template Last Updated**: 2026-02-01
**Repository**: SpecLedger
**CLI Version**: 1.0.0
