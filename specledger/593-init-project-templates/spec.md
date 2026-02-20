# Feature Specification: Project Template & Coding Agent Selection

**Feature Branch**: `593-init-project-templates`
**Created**: 2026-02-20
**Status**: Draft
**Input**: User description: "Enable developers to select from 7 business-defined project templates (General Purpose, Full-Stack, Batch Data Processing, Real-Time Workflow, ML Image Processing, Real-Time Data Pipeline, AI Chatbot) and choose their preferred coding agent (Claude Code, OpenCode, or None) during the interactive `sl new` command. The system generates a unique project ID (UUID v4) for each project to enable session storage and tracking in Supabase."

**Note**: Infrastructure setup (AWS CDK definitions) will be handled by a separate `sl infra` command in a future feature, keeping project initialization focused and fast.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Select Project Template During Initialization (Priority: P1)

A developer runs `sl new` and interactively selects from 7 available project templates to match their project type. The system displays each template with its name, one-line description, and key technology characteristics. The developer navigates with arrow keys and confirms their selection.

**Why this priority**: This is the core value proposition - enabling developers to start with a project structure tailored to their use case rather than a generic template. Without this, the feature provides no value.

**Independent Test**: Can be fully tested by running `sl new`, selecting a template (e.g., "Full-Stack Application"), completing the flow, and verifying the created project contains the expected directory structure and files for that template type.

**Acceptance Scenarios**:

1. **Given** a developer runs `sl new` interactively, **When** they reach the template selection step, **Then** they see 7 template options with names, descriptions, and technology tags
2. **Given** the template selection screen is displayed, **When** the developer presses up/down arrow keys, **Then** the selection cursor moves and wraps around at boundaries
3. **Given** a template is highlighted, **When** the developer presses Enter, **Then** the system stores the selection and proceeds to the agent selection step
4. **Given** the developer selects "Full-Stack Application" template, **When** project creation completes, **Then** the project contains backend/ and frontend/ directories with appropriate starter files

---

### User Story 2 - Select Coding Agent Configuration (Priority: P1)

A developer selects their preferred coding agent (Claude Code, OpenCode, or None) during project initialization. The system creates agent-specific configuration directories and files based on the selection. For Claude Code, the system also creates a settings file with session capture hooks pre-configured.

**Why this priority**: This is essential for the feature's integration with AI coding workflows. Different teams use different agents, and the project must support their choice from the start.

**Independent Test**: Can be fully tested by running `sl new`, selecting any template, choosing an agent option (e.g., "OpenCode"), and verifying the created project contains the correct agent configuration directory (.opencode/) with appropriate files.

**Acceptance Scenarios**:

1. **Given** a developer completes template selection, **When** they reach the agent selection step, **Then** they see 3 agent options: Claude Code (default), OpenCode, and None
2. **Given** the agent selection screen is displayed, **When** the developer presses Enter on "Claude Code", **Then** the system creates a .claude/ directory with commands and skills subdirectories
3. **Given** the developer selects "Claude Code", **When** project creation completes, **Then** the project contains .claude/settings.json with session capture hooks configured for the project ID
4. **Given** the developer selects "OpenCode", **When** project creation completes, **Then** the project contains .opencode/ directory with commands and skills, plus opencode.json configuration file
5. **Given** the developer selects "None", **When** project creation completes, **Then** the project contains no agent-specific directories

---

### User Story 3 - Generate Unique Project ID (Priority: P1)

Every new project created with `sl new` receives a unique UUID v4 identifier stored in the project metadata. This ID enables session tracking and storage in Supabase, allowing sessions to be organized by project.

**Why this priority**: This is a prerequisite for session storage functionality. Without unique project IDs, there's no way to associate captured sessions with specific projects in Supabase.

**Independent Test**: Can be fully tested by creating multiple projects with `sl new` and verifying each project's specledger.yaml contains a unique UUID in the project.id field.

**Acceptance Scenarios**:

1. **Given** a developer creates a new project, **When** project initialization completes, **Then** the specledger.yaml file contains a project.id field with a valid UUID v4
2. **Given** a developer creates 1000 different projects, **When** checking all project IDs, **Then** all IDs are unique with no collisions
3. **Given** a project with ID "550e8400-e29b-41d4-a716-446655440000", **When** Claude Code captures a session, **Then** the session is stored in Supabase under this project ID

---

### User Story 4 - Non-Interactive Template Selection (Priority: P2)

A developer or CI/CD system runs `sl new` with `--template` and `--agent` flags to create a project without interactive prompts. This enables automated project creation in scripts and pipelines.

**Why this priority**: While interactive mode is the primary use case, non-interactive mode is essential for automation and CI/CD workflows. Teams need to script project creation.

**Independent Test**: Can be fully tested by running `sl new --project-name test --project-dir /tmp/test --short-code t --template full-stack --agent opencode` and verifying the project is created without any prompts.

**Acceptance Scenarios**:

1. **Given** a CI/CD script runs `sl new` with all required flags, **When** the command executes, **Then** the project is created without any interactive prompts
2. **Given** a developer provides `--template invalid-template`, **When** the command executes, **Then** the system returns an error message listing available templates and exits with non-zero code
3. **Given** a developer provides `--template full-stack --agent claude-code`, **When** project creation completes, **Then** the project structure matches exactly what would be created in interactive mode with the same selections
4. **Given** a developer runs `sl new` with --force flag and target directory exists, **When** the command executes, **Then** the system deletes the existing directory completely and creates a fresh project from scratch

---

### User Story 5 - List Available Templates (Priority: P2)

A developer runs `sl new --list-templates` to see all available project templates with their descriptions and characteristics before deciding which to use.

**Why this priority**: This enables developers to make informed decisions about template selection without entering the full interactive flow. It's a discovery feature that improves usability.

**Independent Test**: Can be fully tested by running `sl new --list-templates` and verifying the output displays all 7 templates with their IDs, names, descriptions, and technology tags.

**Acceptance Scenarios**:

1. **Given** a developer runs `sl new --list-templates`, **When** the command executes, **Then** the output displays all 7 templates with their IDs, names, and descriptions
2. **Given** the template list is displayed, **When** the developer reads the output, **Then** each template shows its technology characteristics (e.g., "Tech: Go, TypeScript, React, REST API")
3. **Given** the list is displayed, **When** the output ends, **Then** the command exits with code 0 and displays usage instructions

---

### User Story 6 - Backward Compatibility with Current Behavior (Priority: P1)

When a developer runs `sl new` and accepts all default selections (General Purpose template + Claude Code agent), the resulting project structure is identical to the current `sl new` behavior. Existing projects without UUIDs are automatically assigned one when first loaded.

**Why this priority**: This ensures zero breaking changes for existing users and workflows. Teams should be able to upgrade without any disruption to their current processes.

**Independent Test**: Can be fully tested by comparing the output of current `sl new` with new `sl new` (accepting defaults) using a directory diff tool and verifying they are identical except for the new metadata fields.

**Acceptance Scenarios**:

1. **Given** a developer runs new `sl new` and accepts default template and agent, **When** project creation completes, **Then** the directory structure matches current `sl new` output exactly
2. **Given** an existing project with specledger.yaml version 1.0.0 (no UUID), **When** the project is loaded by any `sl` command, **Then** the system generates a UUID and updates the metadata to version 1.1.0
3. **Given** the metadata schema changes to version 1.1.0, **When** loading old projects, **Then** the system continues to work without errors and migrates metadata transparently

---

### User Story 7 - Template-Specific Claude Code Settings (Priority: P2)

When a developer selects Claude Code as their agent, the system creates a .claude/settings.json file with session capture hooks configured to use the project's unique ID. This ensures all coding sessions are automatically captured and organized by project in Supabase.

**Why this priority**: This automates the integration between project creation and session tracking. Developers don't need to manually configure session capture for each project.

**Independent Test**: Can be fully tested by creating a project with Claude Code selected, reading .claude/settings.json, and verifying it contains the session capture hook with the correct project ID.

**Acceptance Scenarios**:

1. **Given** a developer selects Claude Code agent, **When** project creation completes, **Then** .claude/settings.json exists with saveTranscripts set to true
2. **Given** .claude/settings.json is created, **When** reading the file, **Then** it contains a PostToolUse hook for Bash commands that runs `sl session capture`
3. **Given** the settings file contains the project ID, **When** a coding session runs in this project, **Then** session transcripts are automatically uploaded to Supabase under the correct project namespace

---

### Edge Cases

- What happens when a developer presses Ctrl+C during template selection? System exits gracefully without creating any files.
- How does the system handle template directory creation when the target directory already exists? System prompts for confirmation to overwrite or aborts if `--force` flag is not provided.
- What happens if UUID generation fails? System logs an error and aborts project creation (though UUID generation should never fail with crypto/rand).
- How does the system handle migration of projects with version 1.0.0 metadata that have no UUID? System auto-generates UUID on first load and updates metadata to version 1.1.0.
- What happens if agent-specific configuration files fail to copy? System logs a warning but continues with project creation, treating it as a non-fatal error.
- How does the system behave when run in a non-TTY environment without required flags? System returns an error message stating that interactive terminal is required or all flags must be provided.
- What happens when --force flag is used with an existing directory? System deletes the entire existing directory and recreates it from scratch with the selected template, providing a clean slate (destructive operation).
- Can developers change templates after project creation? No - template selection is permanent once the project is created. Changing templates would require creating a new project and manually migrating code.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST display 7 project template options during interactive `sl new` command: General Purpose, Full-Stack Application, Batch Data Processing, Real-Time Workflow, ML Image Processing, Real-Time Data Pipeline, AI Chatbot
- **FR-002**: System MUST display for each template: name, one-line description, and key technology characteristics
- **FR-003**: System MUST allow navigation through template options using up/down arrow keys with wraparound
- **FR-004**: System MUST display 3 coding agent options: Claude Code (default), OpenCode, None
- **FR-005**: System MUST generate a unique UUID v4 for each new project
- **FR-006**: System MUST store project UUID in specledger.yaml under project.id field
- **FR-007**: System MUST store selected template ID in specledger.yaml under project.template field
- **FR-008**: System MUST store selected agent ID in specledger.yaml under project.agent field
- **FR-009**: System MUST update metadata schema version to 1.1.0 for new projects
- **FR-010**: System MUST create template-specific directory structures based on selected template
- **FR-011**: System MUST create .claude/ directory with commands/ and skills/ subdirectories when Claude Code is selected
- **FR-012**: System MUST create .opencode/ directory with commands/ and skills/ subdirectories plus opencode.json when OpenCode is selected
- **FR-013**: System MUST create .claude/settings.json with session capture hooks configured when Claude Code is selected
- **FR-014**: System MUST include project UUID in .claude/settings.json for session organization
- **FR-015**: System MUST NOT create agent-specific directories when "None" is selected
- **FR-016**: System MUST support `--template <id>` flag for non-interactive template selection
- **FR-017**: System MUST support `--agent <id>` flag for non-interactive agent selection
- **FR-018**: System MUST support `--list-templates` flag to display all available templates
- **FR-019**: System MUST validate template ID against available templates and show error if invalid
- **FR-020**: System MUST validate agent ID against available agents and show error if invalid
- **FR-021**: System MUST require all flags (--project-name, --project-dir, --short-code, --template, --agent) when run in non-TTY environment
- **FR-022**: System MUST produce identical output for "General Purpose + Claude Code" selections as current `sl new` behavior (plus infra/ directory)
- **FR-023**: System MUST auto-generate UUID for existing projects (version 1.0.0) when first loaded and update to version 1.1.0
- **FR-024**: System MUST use github.com/google/uuid library for UUID generation
- **FR-025**: System MUST create projects without agent configuration files when "None" is selected, only creating specledger.yaml and template structure
- **FR-026**: System MUST delete entire existing directory when --force flag is provided, recreating project from scratch
- **FR-027**: System MUST prevent template changes after project creation - template selection is immutable
- **FR-028**: System MUST log all template operations using structured logging (template selection, file operations, errors) to stdout

### Key Entities

- **Project Template**: Represents a predefined project structure with directories, files, and configuration for a specific use case (e.g., Full-Stack, ML Image Processing). Contains template ID, name, description, technology characteristics, and embedded file paths.
- **Coding Agent Configuration**: Represents the configuration for an AI coding agent (Claude Code, OpenCode, or None). Contains agent ID, name, description, and configuration directory name.
- **Project Metadata**: Stored in specledger.yaml, contains project identification (UUID, name, short code), creation timestamps, selected template ID, selected agent ID, and schema version.
- **Session Capture Hook**: Configuration in .claude/settings.json that triggers `sl session capture` after Bash tool usage, enabling automatic session recording organized by project UUID.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Developers can complete the full interactive template and agent selection flow in under 60 seconds
- **SC-002**: All 7 project templates produce valid, non-empty project structures with appropriate directories and starter files
- **SC-003**: Projects created with default selections (General Purpose + Claude Code) are identical to current `sl new` output except for new metadata fields
- **SC-004**: Non-interactive mode with all required flags completes project creation without any prompts or user interaction
- **SC-005**: Every project created receives a unique UUID with zero collisions across 10,000 test projects
- **SC-006**: Template structures are recognizable and distinct - developers can identify the template type by examining the directory structure
- **SC-007**: Agent selection produces correct configuration directories - Claude Code creates .claude/, OpenCode creates .opencode/, None creates neither
- **SC-008**: Session capture hooks are automatically configured for Claude Code projects, requiring no manual setup by developers
- **SC-009**: Existing projects without UUIDs are successfully migrated when loaded, with no data loss or errors
- **SC-010**: Template list display (--list-templates) shows all 7 options in under 100 milliseconds
- **SC-011**: Template selection is immutable - attempting to change template after creation results in clear error message directing developer to create new project
- **SC-012**: All template operations are logged with structured output including timestamps, operation types, file paths, and error details

### Previous work

This feature extends and builds upon:

#### Epic: Project Initialization & Configuration

- **Feature 011-streamline-onboarding**: Established the interactive TUI flow with Bubble Tea framework, text input components, and step-based state machine for project creation
- **Feature 005-embedded-templates**: Created the embedded template system with manifest-based discovery, file copying with pattern matching, and `//go:embed` filesystem integration
- **Feature 004-thin-wrapper-redesign**: Defined the ProjectMetadata structure, specledger.yaml schema, and metadata management system
- **Feature 010-checkpoint-session-capture**: Implemented session capture functionality that this feature integrates with through .claude/settings.json hooks

## Clarifications

### Session 2026-02-20

- Q: What should happen when a project directory already exists and developer runs 'sl new' with --force flag? → A: Delete entire directory and recreate from scratch (Clean slate approach - removes all existing files and creates new project)
- Q: Should developers be able to change/upgrade project template after initial creation? → A: No - template is immutable after creation (Template selection is permanent, avoids complex migration logic)
- Q: What observability should be built into template operations for debugging and monitoring? → A: Structured logs only (Log template selection, file operations, errors to stdout/file using standard logging)

**Note**: Infrastructure-related clarifications (AWS credentials, secrets management) were removed as infrastructure setup has been moved to a separate `sl infra` command feature.

## Dependencies & Assumptions

### Dependencies

- **github.com/google/uuid v1.6.0**: Required for UUID v4 generation with cryptographically secure randomness
- **Existing TUI Framework**: Depends on Bubble Tea (v1.3.10), Bubbles (v0.21.1), and Lipgloss (v1.1.0) already in use
- **Existing Template System**: Depends on embedded template loading system from feature 005
- **Session Capture Command**: Depends on `sl session capture` command from feature 010

### Assumptions

- All 7 template structures will be embedded in the binary at compile time, increasing binary size by approximately 10-15MB
- Template files use standard directory layouts recognizable by developers familiar with Go, React, Airflow, Temporal, TensorFlow, Kafka ecosystems
- Session capture hooks assume developers are running Claude Code in projects where this feature is used
- UUID collision probability is negligible (122 bits of randomness per UUID v4)
- Terminal environments support arrow key navigation and standard ANSI escape codes for TUI rendering
- Non-interactive mode is primarily used in CI/CD environments where all project parameters can be predetermined
- Metadata migration from v1.0.0 to v1.1.0 is backward compatible - old fields remain unchanged, new fields are added
- OpenCode configuration format follows similar patterns to Claude Code for ease of migration between agents
- Template characteristics (technology tags) are concise enough to display in a single line without overwhelming the UI
- The --force flag performs destructive deletion of existing directories - no backup or merge is attempted
- Template selection is permanent and immutable after project creation - no template migration/upgrade commands exist
- Structured logging captures all template operations, file copies, and errors for debugging
- Infrastructure setup (AWS CDK, compute, storage) will be handled by a separate `sl infra` command in a future feature
