---
description: Manage specification dependencies using the SpecLedger CLI deps commands.
---

## User Input

```text
$ARGUMENTS
```

You **MUST** consider the user input before proceeding (if not empty).

## Overview

The SpecLedger CLI provides dependency management for specifications. Dependencies are external specs that your project references or builds upon. They are cached locally for offline use and easy LLM access.

## When to Use

Use this command when:
- **Adding dependencies** - Your spec references other specifications
- **Listing dependencies** - Understanding what specs your project depends on
- **Downloading dependencies** - Fetching and caching external specs for offline use
- **Updating dependencies** - Keeping specs at latest compatible versions
- **Removing dependencies** - Cleaning up unused spec references

## How Caching Works

Dependencies are cached locally at `~/.cache/specledger/` (similar to Go's module cache):

- **Automatic download**: `sl deps resolve` downloads all dependencies
- **Offline access**: Once cached, specs are available offline
- **LLM integration**: Cached specs can be easily read and referenced by AI agents
- **Version pinning**: Each cached version has a cryptographic hash (in spec.sum)

## Commands Reference

### sl deps add

Add a specification dependency to your project.

**Usage:**
```bash
sl deps add <repo-url> [branch] [spec-path] [--alias <name>]
```

**Examples:**
```bash
# Basic usage - uses default branch and spec.md
sl deps add git@github.com:org/api-spec

# With specific branch
sl deps add git@github.com:org/api-spec v1.0

# With specific spec path
sl deps add git@github.com:org/api-spec main specs/api.md

# With alias for easy reference
sl deps add git@github.com:org/api-spec --alias api
```

**What happens:**
- The dependency is added to `specs/spec.mod`
- It will be downloaded when you run `sl deps resolve`
- Once downloaded, it's cached locally for offline use

### sl deps list

List all declared dependencies.

**Usage:**
```bash
sl deps list
```

**Output:**
```
Dependencies (3):

1. git@github.com:org/api-spec
   Version: main
   Spec: spec.md
   Cached: ✓
   Alias: api

2. git@github.com:org/auth-spec
   Version: v2.0
   Spec: spec.md
   Cached: ✓

3. git@github.com:org/db-spec
   Version: main
   Spec: schema.md
   Cached: ✗ (not downloaded yet)
```

### sl deps resolve

Download and cache all dependencies.

**Usage:**
```bash
sl deps resolve
```

**What happens:**
- Downloads all dependencies from `specs/spec.mod`
- Validates versions and commits
- Caches them locally at `~/.cache/specledger/`
- Updates `specs/spec.sum` with cryptographic hashes

**After running this:**
- Dependencies are available offline
- LLMs can read and reference the cached specs
- Reproducible builds with locked hashes

### sl deps remove

Remove a dependency.

**Usage:**
```bash
sl deps remove <repo-url>
```

**Example:**
```bash
sl deps remove git@github.com:org/api-spec
```

**Note:** The local cache is kept (for potential future use).

### sl deps update

Update dependencies to latest compatible versions.

**Usage:**
```bash
sl deps update [repo-url]
```

**Examples:**
```bash
# Update all dependencies
sl deps update

# Update specific dependency
sl deps update git@github.com:org/api-spec
```

## Referencing Dependencies

Once a dependency is added and resolved, reference it in your specifications:

**In spec.md:**
```markdown
## API Integration

This component depends on: `api` (git@github.com:org/api-spec.git/main/spec.md)

Key requirements from api:
- Authentication: Must use OAuth2 flow
- Rate limits: 100 req/min per user
```

**In templates:**
```markdown
## Database Schema

Extends base schema defined in: `db` (git@github.com:org/db-spec.git/main/schema.md)

Additional fields:
- user_preferences: JSONB
```

## Reading Cached Dependencies

For LLMs and agents, cached dependencies can be read directly:

```bash
# Cached dependencies are at:
~/.cache/specledger/<repo-url>/<commit>/<spec-path>

# Example:
~/.cache/specledger/github.com/org/api-spec/a1b2c3d/spec.md
```

## Best Practices

1. **Use meaningful aliases** - Short, memorable names for frequently referenced specs
2. **Pin versions** - Use specific tags (not just `main`) for reproducibility
3. **Document references** - Always note which parts of external specs you use
4. **Run resolve regularly** - `sl deps resolve` downloads and caches dependencies
5. **Commit spec.sum** - Lockfile ensures reproducible builds
6. **Use aliases for LLMs** - Makes it easier for AI to understand dependencies

## Error Handling

**Common errors:**

- **"Not a SpecLedger project"** - Navigate to a project directory or run `sl new` first
- **"Invalid repository URL"** - Use `git@` format (e.g., `git@github.com:org/spec`)
- **"Dependency not found"** - Check the URL and branch in `specs/spec.mod`
- **"Not cached"** - Run `sl deps resolve` to download dependencies

## Getting Help

```bash
sl deps --help          # Show all deps commands
sl deps add --help      # Show add command help
sl deps list --help     # Show list command help
```
