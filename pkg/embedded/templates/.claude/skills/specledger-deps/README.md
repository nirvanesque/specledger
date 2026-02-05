# SpecLedger Dependency Management Skill

This skill documents the SpecLedger CLI `deps` commands for managing external specification dependencies.

## Commands

| Command | Description |
|---------|-------------|
| `sl deps add <repo-url>` | Add a specification dependency |
| `sl deps list` | List all declared dependencies |
| `sl deps resolve` | Download and cache dependencies |
| `sl deps update` | Update dependencies to latest versions |
| `sl deps remove <repo-url>` | Remove a dependency |

## Quick Start

```bash
# Add a dependency
sl deps add git@github.com:org/project-spec main spec.md --alias myproject

# List all dependencies
sl deps list

# Download and cache
sl deps resolve

# Update to latest versions
sl deps update
```

## See Also

- **[specledger-issue-tracking](../specledger-issue-tracking/README.md)** - Track work across sessions with dependency graphs
