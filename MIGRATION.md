# Migration Guide: From Legacy CLI to Unified CLI

This guide helps you transition from the old bash `sl` script and standalone Go `specledger` CLI to the new unified CLI tool.

## Overview

The new unified CLI consolidates both the bash `sl` bootstrap script and the Go `specledger` CLI into a single binary called `sl`. The legacy `specledger` command still works as an alias for backward compatibility.

## Key Changes

### Command Renaming

| Old Command | New Command |
|-------------|-------------|
| `./sl` (bash script) | `sl` (Go binary) |
| `./specledger` (Go binary) | `sl` (Go binary) |
| `./specledger new` | `sl new` |
| `./specledger deps ...` | `sl deps ...` |

### Backward Compatibility

The legacy `specledger` command continues to work as an alias:

```bash
sl new        # Works (was: ./specledger new)
sl deps list  # Works (was: ./specledger deps list)
```

## Migration Steps

### Step 1: Install the New CLI

#### From Source

```bash
# Clone and build
git clone https://github.com/your-org/specledger.git
cd specledger
make build

# Copy to your path
cp bin/sl ~/bin/
# Or use sudo for system-wide install
sudo cp bin/sl /usr/local/bin/
```

#### Using Installation Script

```bash
curl -fsSL https://raw.githubusercontent.com/your-org/specledger/main/scripts/install.sh | sudo bash
```

#### Using Go

```bash
go install github.com/your-org/specledger/cmd/sl@latest
```

### Step 2: Verify Installation

```bash
# Check version
sl --version

# Check help
sl --help
```

### Step 3: Update Your Scripts

Replace any references to the old `sl` script with the new `sl` command:

#### Before (bash script)

```bash
#!/bin/bash
# Old script location
./sl new --project-name myproject --short-code mp
```

#### After (new CLI)

```bash
#!/bin/bash
# Use the new CLI
sl new --project-name myproject --short-code mp
```

### Step 4: Update CI/CD Pipelines

Update your CI/CD configuration files:

#### GitHub Actions (before)

```yaml
- name: Create project
  run: ./sl new --ci --project-name myproject --short-code mp
```

#### GitHub Actions (after)

```yaml
- name: Create project
  run: sl new --ci --project-name myproject --short-code mp
```

#### GitLab CI (before)

```yaml
create_project:
  script:
    - ./sl new --ci --project-name myproject --short-code mp
```

#### GitLab CI (after)

```yaml
create_project:
  script:
    - sl new --ci --project-name myproject --short-code mp
```

### Step 5: Update Shell Aliases

If you have custom shell aliases for `sl`, you can simplify them:

#### Before

```bash
alias sl='./sl'
alias sl-new='./sl new'
```

#### After

```bash
# The binary is now named sl
alias sl-new='sl new'
```

Or you can use the alias directly:

```bash
alias sl-new='sl --help'  # Shows help
```

## Breaking Changes

### No Breaking Changes for CLI

There are **no breaking changes** to the CLI interface. All existing commands continue to work exactly as before.

### Minor Changes

1. **Binary Name Changed**: From `specledger` to `sl`
2. **Installation Methods Added**: Multiple new ways to install
3. **TUI Available**: Interactive terminal UI now available

## Common Migration Issues

### Issue: "command not found: sl"

**Cause**: The `sl` binary is not in your PATH.

**Solution**:
```bash
# Check if binary exists
ls ~/bin/sl
# or
ls /usr/local/bin/sl

# If it exists, add to PATH
export PATH="$HOME/bin:$PATH"
# Add to your shell config: ~/.bashrc or ~/.zshrc
echo 'export PATH="$HOME/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Issue: Old `sl` script still works

**Cause**: The old bash script is still in your PATH.

**Solution**:
```bash
# Remove old script
rm ~/bin/sl

# Or rename to something else
mv ~/bin/sl ~/bin/sl.old
```

### Issue: Permissions error

**Cause**: The binary is not executable.

**Solution**:
```bash
chmod +x ~/bin/sl
# or
sudo chmod +x /usr/local/bin/sl
```

## Testing Your Migration

After migrating, run these tests to ensure everything works:

```bash
# 1. Check version
sl --version
# Expected: sl version 1.0.0

# 2. Check help
sl --help
# Expected: Shows all commands

# 3. Create a project
sl new --ci --project-name migration-test --short-code mt
# Expected: Project created successfully

# 4. Verify project
ls ~/demos/migration-test
# Expected: Project directory exists

# 5. Use dependency commands (if in a project)
cd ~/demos/migration-test
sl deps list
# Expected: No errors, lists dependencies (empty initially)

# 6. Clean up test project
rm -rf ~/demos/migration-test
```

## Rollback Plan

If you need to rollback to the old CLI:

### Option 1: Keep Both (Recommended)

The new `sl` is just a different binary name. Keep your old `specledger` installation and `sl` bash script for now:

```bash
# You can run both:
./specledger new      # Old CLI
sl new                # New CLI
```

### Option 2: Restore Old Scripts

```bash
# Restore the old sl script
git checkout HEAD -- sl
chmod +x sl

# Or download from backup
curl -O https://raw.githubusercontent.com/your-org/specledger/main/sl
chmod +x sl
```

### Option 3: Uninstall New CLI

```bash
# Remove new binary
sudo rm /usr/local/bin/sl
# or
rm ~/bin/sl
```

## Feature Comparison

| Feature | Old CLI | New CLI |
|---------|---------|---------|
| Interactive TUI | No | Yes |
| Non-interactive mode | No | Yes (`--ci`) |
| Command name | `specledger` | `sl` (with `specledger` alias) |
| Installation | Manual | Multiple methods |
| Error messages | Basic | Detailed with suggestions |

## Next Steps

1. **Read the README**: Check out the [README.md](README.md) for full documentation
2. **Try the TUI**: Run `sl new` in your terminal
3. **Update Documentation**: Update your team's docs to reference the new CLI
4. **Share Feedback**: Report any issues or feature requests

## Getting Help

- **Documentation**: [README.md](README.md)
- **Issues**: [GitHub Issues](https://github.com/your-org/specledger/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/specledger/discussions)

## Questions?

If you encounter any issues during migration:

1. Check this migration guide
2. Read the troubleshooting section
3. Open an issue on GitHub
4. Join the discussion in GitHub Discussions

---

**Migration Version**: 1.0.0
**Updated**: 2026-01-31
