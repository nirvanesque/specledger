# SpecLedger Architecture

## Overview

SpecLedger is a **thin wrapper CLI** that orchestrates tool installation and project bootstrapping while delegating SDD (Specification-Driven Development) workflows to user-chosen frameworks.

## Design Philosophy

**Do One Thing Well**: SpecLedger focuses on:
- ✅ Project bootstrapping and initialization
- ✅ Tool prerequisite checking and installation
- ✅ Specification dependency management
- ✅ Framework neutrality

**Delegate, Don't Duplicate**: SpecLedger does NOT:
- ❌ Implement its own SDD workflow commands
- ❌ Duplicate functionality from Spec Kit or OpenSpec
- ❌ Force a specific development methodology

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                         SpecLedger CLI                       │
│                      (sl command)                           │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐           │
│  │  sl new    │  │  sl init   │  │  sl doctor │           │
│  └─────┬──────┘  └─────┬──────┘  └─────┬──────┘           │
│        │               │               │                   │
│        ▼               ▼               ▼                   │
│  ┌──────────────────────────────────────────────────┐     │
│  │         Prerequisites Checker                   │     │
│  │  - Detects mise, bd, perles                     │     │
│  │  - Auto-installs via mise                        │     │
│  │  - Checks frameworks (specify, openspec)         │     │
│  └───────────────────┬──────────────────────────────┘     │
│                      │                                    │
│                      ▼                                    │
│  ┌──────────────────────────────────────────────────┐     │
│  │         YAML Metadata System                     │     │
│  │  - specledger/specledger.yaml                    │     │
│  │  - Project info, framework choice, dependencies  │     │
│  └───────────────────┬──────────────────────────────┘     │
│                      │                                    │
└──────────────────────┼────────────────────────────────────┘
                       │
       ┌───────────────┼───────────────┐
       ▼               ▼               ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│   Spec Kit  │ │  OpenSpec   │ │    None     │
│  (specify)  │ │  (openspec) │ │             │
└─────────────┘ └─────────────┘ └─────────────┘
       │               │               │
       ▼               ▼               ▼
┌─────────────────────────────────────────────────────────┐
│                  User's Workflow                       │
│  - Specification, Planning, Implementation              │
│  - Issue Tracking (bd - beads)                         │
└─────────────────────────────────────────────────────────┘
```

## Components

### 1. CLI Commands (`pkg/cli/commands/`)

| Command | File | Purpose |
|---------|------|---------|
| `sl new` | `bootstrap.go` | Create new project with interactive TUI |
| `sl init` | `bootstrap.go` | Initialize SpecLedger in existing repo |
| `sl doctor` | `doctor.go` | Check tool installation status |
| `sl migrate` | `migrate.go` | Convert .mod to YAML format |
| `sl deps` | `deps.go` | Manage dependencies (add/list/remove) |

### 2. Prerequisites Checker (`pkg/cli/prerequisites/`)

**Purpose**: Detect and ensure required tools are installed

**Core Tools (Required)**:
- `mise` - Version manager for tools
- `bd` (beads) - Issue tracking
- `perles` - Workflow automation

**Framework Tools (Optional)**:
- `specify` - Spec Kit CLI
- `openspec` - OpenSpec CLI

**Behavior**:
- Interactive mode: Prompts user to install missing tools
- CI mode: Auto-installs without prompts
- Provides clear error messages with install instructions

### 3. Metadata System (`pkg/cli/metadata/`)

**Purpose**: Manage project configuration in YAML format

**Files**:
- `schema.go` - Go structs for YAML schema
- `yaml.go` - Read/write YAML files
- `migration.go` - Convert legacy .mod to YAML

**Schema** (`specledger/specledger.yaml`):
```yaml
version: "1.0.0"
project:
  name: string
  short_code: string
  created: timestamp
  modified: timestamp
  version: string
framework:
  choice: speckit | openspec | both | none
  installed_at: timestamp (optional)
dependencies:
  - url: git URL
    branch: string (default: main)
    path: string (default: spec.md)
    alias: string
    resolved_commit: SHA hash
```

### 4. TUI (`pkg/cli/tui/`)

**Framework**: Bubble Tea for terminal UI

**Components**:
- `sl_new.go` - Interactive bootstrap flow
- Framework selection (checkboxes for Spec Kit, OpenSpec)
- Project name, short code input
- Project directory selection

### 5. Embedded Templates (`pkg/embedded/templates/`)

**Purpose**: Files copied to new projects

**Structure**:
```
templates/
├── .claude/
│   ├── commands/specledger.{deps,adopt,resume}.md
│   └── skills/specledger-{deps,issue-tracking}/
├── .beads/
│   └── config.yaml
├── specledger/
│   ├── specledger.yaml (YAML metadata)
│   └── AGENTS.md
└── mise.toml (with commented framework options)
```

## Decision Records

### Why Thin Wrapper Architecture?

**Problem**: Previous SpecLedger duplicated SDD workflow functionality from Spec Kit.

**Solution**: Redesign as orchestrator that:
1. Sets up projects with proper tooling
2. Manages spec dependencies (unique value)
3. Lets users choose their SDD framework
4. Doesn't compete with frameworks

**Trade-offs**:
- ✅ Clearer boundaries and responsibilities
- ✅ Users can switch frameworks easily
- ✅ Less code to maintain
- ❌ Need to explain framework choice

### Why mise for Tool Management?

**Decision**: Use mise as the universal tool installer

**Rationale**:
- Single tool for all dependencies (Go, Python, Node, binaries)
- Cross-platform support
- Declarative configuration (`mise.toml`)
- Easy to add custom tools

**Alternatives Considered**:
- `go install` - Go-only, requires GOPATH setup
- Homebrew apt - Platform-specific
- Manual installation - Poor UX

### Why YAML Instead of .mod?

**Decision**: Replace `.mod` file format with YAML

**Rationale**:
- Human-readable and editable
- Standard format with good Go support
- Easy to extend with new fields
- Better validation and error messages

**Migration Path**:
- `sl migrate` command converts .mod to YAML
- Automatic detection with deprecation warning
- Backward compatibility during transition

### Why Check Prerequisites?

**Decision**: Actively check and install required tools

**Rationale**:
- Reduces onboarding friction
- Prevents cryptic errors from missing tools
- Consistent development environments
- Better CI/CD experience

## Extension Points

### Adding New Prerequisite Tools

Edit `pkg/cli/prerequisites/checker.go`:
```go
var RequiredTools = []Tool{
    {Name: "mise", Category: ToolCategoryCore},
    {Name: "bd", Category: ToolCategoryCore},
    {Name: "perles", Category: ToolCategoryCore},
    // Add new tool here
}
```

### Adding New SDD Frameworks

1. Update TUI (`pkg/cli/tui/sl_new.go`) with framework option
2. Update mise.toml template with installation command
3. Update `pkg/cli/prerequisites/checker.go` with detection
4. Update documentation

### Extending YAML Schema

Edit `pkg/cli/metadata/schema.go`:
```go
type ProjectMetadata struct {
    Version    string        `yaml:"version"`
    Project    ProjectInfo   `yaml:"project"`
    Framework  FrameworkInfo `yaml:"framework"`
    Dependencies []Dependency `yaml:"dependencies,omitempty"`
    // Add new fields here
}
```

## Data Flow

### Creating a New Project (`sl new`)

```
User runs "sl new"
        │
        ▼
┌───────────────────┐
│  TUI collects     │
│  - Project name   │
│  - Short code     │
│  - Framework(s)   │
│  - Directory      │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│  Check Prerequisites │
│  - Detect mise    │
│  - Prompt install │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│  Create Project  │
│  - Copy templates │
│  - Write YAML     │
│  - Run mise trust │
└───────────────────┘
```

### Adding Dependencies (`sl deps add`)

```
User runs "sl deps add <url>"
        │
        ▼
┌───────────────────┐
│  Load YAML        │
│  specledger.yaml  │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│  Validate         │
│  - Check URL      │
│  - Check duplicate│
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│  Append to list   │
│  dependencies[]   │
└─────────┬─────────┘
          │
          ▼
┌───────────────────┐
│  Save YAML        │
└───────────────────┘
```

## Dependencies Between Projects

SpecLedger enables **spec dependencies** - projects can reference specifications from other repositories:

```
Project A (specledger/specledger)
├── specledger.yaml
└── dependencies:
    └── url: git@github.com:user/common-spec
      branch: main
      path: spec.md
      alias: common

Project B (user/my-api)
├── specledger.yaml
└── dependencies:
    ├── url: git@github.com:user/common-spec
    │   alias: common
    └── url: git@github.com:specledger/specledger
        alias: sl-base
```

**Cached locally at**: `~/.specledger/cache/<domain>/<org>/<repo>/<commit>/`

## Security Considerations

1. **Tool Installation**: SpecLedger only runs `mise install` for user-confirmed tools
2. **Git Operations**: Uses shallow clones for security
3. **Cache Validation**: SHA-256 hashes for content verification
4. **No Remote Execution**: All operations are local

## Performance

- **Bootstrap**: <3 minutes including tool installation
- **sl doctor**: <2 seconds for tool detection
- **YAML parsing**: <10ms for typical projects
- **Dependency resolution**: Depends on git clone times

## Future Enhancements

- [ ] Dependency graph visualization (`sl graph show`)
- [ ] Automatic dependency updates
- [ ] Global configuration file
- [ ] Project templates
- [ ] CI/CD integration best practices
