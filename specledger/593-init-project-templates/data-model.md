# Data Model: Project Template & Coding Agent Selection

**Feature**: 593-init-project-templates
**Date**: 2026-02-20

This document defines the key entities, their fields, relationships, and validation rules for the template and agent selection feature.

---

## Entity: Project Template

**Description**: Represents a predefined project structure with directories, files, and configuration for a specific use case (e.g., Full-Stack, ML Image Processing).

### Fields

| Field | Type | Required | Description | Validation |
|-------|------|----------|-------------|------------|
| `ID` | string | Yes | Unique template identifier (kebab-case) | Must match pattern `^[a-z][a-z0-9-]*$`, max 50 chars |
| `Name` | string | Yes | Human-readable template name | Non-empty, max 100 chars |
| `Description` | string | Yes | One-line description of template purpose | Non-empty, max 200 chars |
| `Version` | string | Yes | Semantic version of template structure | Must match semver pattern `^\d+\.\d+\.\d+$` |
| `Path` | string | Yes | Relative path within embedded templates/ directory | Non-empty, must exist in embedded FS |
| `Technologies` | []string | No | Key technologies used in template (e.g., "Go", "React", "Kafka") | Max 10 items, each max 50 chars |
| `Structure` | []string | Yes | Top-level directories created by template | Non-empty array, each item max 200 chars |

### Go Struct

```go
type ProjectTemplate struct {
    ID           string   `yaml:"id"`
    Name         string   `yaml:"name"`
    Description  string   `yaml:"description"`
    Version      string   `yaml:"version"`
    Path         string   `yaml:"path"`
    Technologies []string `yaml:"technologies,omitempty"`
    Structure    []string `yaml:"structure"`
}
```

### Relationships

- **Has Many**: Template Files (embedded in binary via //go:embed)
- **Referenced By**: ProjectMetadata (via `project.template` field)

---

## Entity: Coding Agent Configuration

**Description**: Represents the configuration for an AI coding agent (Claude Code, OpenCode, or None).

### Fields

| Field | Type | Required | Description | Validation |
|-------|------|----------|-------------|------------|
| `ID` | string | Yes | Unique agent identifier | Must match existing agent |
| `Name` | string | Yes | Human-readable agent name | Non-empty |
| `Description` | string | Yes | One-line description | Non-empty |
| `ConfigDir` | string | No | Config directory name | Must start with ".", empty for "None" |
| `SupportsSessionCapture` | bool | No | Session capture support | Defaults to false |

### Go Struct

```go
type CodingAgentConfig struct {
    ID                      string
    Name                    string
    Description             string
    ConfigDir               string
    SupportsSessionCapture  bool
}
```

---

## Entity: Project Metadata (v1.1.0)

**Description**: Stored in specledger.yaml with project identification, timestamps, template ID, and agent ID.

### New Fields (v1.1.0)

| Field | Type | Required | Description | Migration |
|-------|------|----------|-------------|-----------|
| `project.id` | uuid.UUID | Yes | Unique project identifier | Auto-generated if missing |
| `project.template` | string | No | Selected template ID | New field |
| `project.agent` | string | No | Selected agent ID | New field |

### YAML Example

```yaml
version: 1.1.0
project:
  id: 550e8400-e29b-41d4-a716-446655440000
  name: my-fullstack-app
  short_code: mfa
  template: full-stack
  agent: claude-code
  created: 2026-02-20T10:30:00Z
  modified: 2026-02-20T10:30:00Z
  version: 0.1.0
```

### Migration (v1.0.0 → v1.1.0)

1. Generate UUID v4 using `uuid.New()`
2. Set `project.template = ""` (omitempty)
3. Set `project.agent = ""` (omitempty)
4. Update `version = "1.1.0"`
5. Save metadata immediately

---

## Entity: Session Capture Hook

**Description**: Configuration in .claude/settings.json for automatic session recording by project UUID.

### JSON Example

```json
{
  "saveTranscripts": true,
  "hooks": {
    "PostToolUse": [{
      "matcher": "Bash",
      "hooks": [{
        "type": "command",
        "command": "sl session capture --project-id 550e8400-e29b-41d4-a716-446655440000"
      }]
    }]
  }
}
```

### Generation

**Trigger**: Project creation when `agent == "claude-code"`

**Process**:
1. Read project UUID from metadata
2. Generate settings.json with PostToolUse hook
3. Write to `.claude/settings.json` (0644 permissions)
