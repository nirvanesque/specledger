---
name: specledger-deps
description: Manage specification dependencies using the SpecLedger CLI deps commands. Use when working with external specification dependencies, dependency resolution, or linking specs across repositories.
---

# SpecLedger Dependency Management

## Overview

The SpecLedger CLI provides commands to manage external specification dependencies. Dependencies are specified in `specledger/specledger.mod`, downloaded and cached locally at `~/.specledger/cache/` (similar to Go's module cache), and can be easily referenced by LLMs.

## When to Use

Use this skill when:
- **Adding new dependencies** - Referencing specs from other repositories
- **Listing dependencies** - Understanding what specs a project depends on
- **Downloading dependencies** - Fetching and caching external specs for offline use
- **Updating dependencies** - Keeping specs at the latest compatible version
- **Removing dependencies** - Cleaning up unused spec references

## How Caching Works

Dependencies are cached locally at `~/.specledger/cache/`:

- **Automatic download**: `sl deps resolve` downloads all dependencies
- **Offline access**: Once cached, specs are available offline
- **LLM integration**: Cached specs can be easily read and referenced by AI agents
- **Version pinning**: Each cached version has a cryptographic hash (in specledger.sum)

```
~/.specledger/cache/
└── github.com/
    └── org/
        └── api-spec/
            └── a1b2c3d/           # Commit hash
                └── spec.md        # Cached spec file
```

## CLI Commands Reference

| Command | Description |
|---------|-------------|
| `sl deps add <repo-url>` | Add a specification dependency |
| `sl deps list` | List all declared dependencies |
| `sl deps resolve` | Download and cache all dependencies |
| `sl deps update` | Update dependencies to latest versions |
| `sl deps remove <repo-url>` | Remove a dependency |

### Adding Dependencies

**Basic usage:**
```bash
sl deps add git@github.com:org/project-spec
```

**With branch and spec path:**
```bash
sl deps add git@github.com:org/project-spec main specs/api.md
```

**With alias:**
```bash
sl deps add git@github.com:org/project-spec --alias myapi
```

**What happens:**
- The dependency is added to `specledger/specledger.mod`
- It will be downloaded when you run `sl deps resolve`
- Once downloaded, it's cached locally for offline use

### Listing Dependencies

```bash
sl deps list
```

Output:
```
Dependencies (2):

1. git@github.com:org/project-spec
   Version: main
   Spec: spec.md
   Cached: ✓
   Alias: myapi

2. git@github.com:org/referenced-spec
   Version: v1.0
   Spec: spec.md
   Cached: ✗ (not downloaded yet)
```

### Resolving (Downloading) Dependencies

```bash
sl deps resolve
```

This command:
1. Reads `specledger/specledger.mod`
2. Fetches external specifications from Git
3. Validates versions and commits
4. Caches them locally at `~/.specledger/cache/`
5. Generates `specledger/specledger.sum` with cryptographic hashes

**After running this:**
- Dependencies are available offline
- LLMs can read and reference the cached specs
- Reproducible builds with locked hashes

### Updating Dependencies

```bash
# Update all dependencies to latest compatible versions
sl deps update

# Update specific dependency
sl deps update git@github.com:org/project-spec
```

### Removing Dependencies

```bash
sl deps remove git@github.com:org/project-spec
```

**Note:** The local cache is kept (for potential future use).

## Dependency Workflow

### Typical Project Setup

```bash
# 1. Initialize a new project
sl new

# 2. Add your first dependency
sl deps add git@github.com:org/spec-base --alias base

# 3. Resolve to fetch and cache
sl deps resolve

# 4. Add more dependencies as needed
sl deps add git@github.com:org/ui-specs

# 5. Resolve again
sl deps resolve
```

### Working with Dependencies in Specs

**In markdown specs:**
```markdown
# User Authentication Spec

## Dependencies
- `base` (git@github.com:org/spec-base/main/specs/auth-base.md)

## Content
This spec extends the base authentication specification...
```

**For LLM reference:**
```markdown
## API Integration

This component depends on: `api` (git@github.com:org/api-spec/main/specs/api.md)

Cached at: ~/.specledger/cache/github.com/org/api-spec/<commit>/specs/api.md
```

## Dependency Format

**spec.mod file structure:**
```yaml
manifest_version: 1
dependencies:
  - repository_url: git@github.com:org/project-spec
    version: main
    spec_path: spec.md
    alias: myproject
    added_at: 2024-01-15T10:00:00Z
```

**spec.sum file structure:**
```json
{
  "lockfile_version": "1",
  "dependencies": [
    {
      "repository_url": "git@github.com:org/project-spec",
      "commit_hash": "abc123def456",
      "content_hash": "xyz789",
      "spec_path": "spec.md",
      "branch": "main",
      "size": 1024,
      "fetched_at": "2024-01-15T10:00:00Z"
    }
  ],
  "generated_at": "2024-01-15T10:00:00Z"
}
```

## Advanced Patterns

### Multiple Aliases for Same Repo

```bash
sl deps add git@github.com:org/shared main specs/common.md --alias common
sl deps add git@github.com:org/shared main specs/ui.md --alias ui
```

### Dependency Version Locking

```bash
# Add with specific version
sl deps add git@github.com:org/project-spec v1.2.3

# This locks to the commit at v1.2.3
# To update, run: sl deps update
```

## Troubleshooting

### "Not a SpecLedger project"

Make sure you're in a project directory with `specledger.mod`:
```bash
ls specledger.mod
sl deps list
```

### "Dependency not found"

```bash
# Check if added correctly
sl deps list

# Check the spec.mod file directly
cat specledger/specledger.mod
```

### "Failed to resolve dependencies"

```bash
# Check internet connection and repo URLs
# Verify the repo exists at the specified path
git ls-remote git@github.com:org/project-spec.git
```

### "Invalid repository URL"

Use valid Git URLs:
- ✅ `git@github.com:org/project.git`
- ✅ `https://github.com/org/project.git`

## Best Practices

1. **Use aliases** for common dependencies to make specs more readable
2. **Lock versions** when you need reproducibility (use tags instead of `main`)
3. **Resolve regularly** to keep `spec.sum` up to date
4. **Commit spec.sum** to ensure reproducible builds
5. **Document dependencies** in your README with their purpose

## Files Reference

| File | Description |
|------|-------------|
| `specledger/specledger.mod` | Dependency manifest and project metadata |
| `specledger/specledger.sum` | Lockfile with resolved commits and hashes |
| `~/.specledger/cache/` | Local cache of downloaded dependencies |

## Related Skills

- **[bd-issue-tracking](.claude/skills/bd-issue-tracking/SKILL.md)** - Use for tracking work across sessions

## Quick Commands Cheat Sheet

```bash
# Add dependency
sl deps add <repo-url> [branch] [spec-path] [--alias <name>]

# List dependencies
sl deps list

# Resolve (download) all
sl deps resolve

# Update all
sl deps update [repo-url]

# Remove dependency
sl deps remove <repo-url>
```
