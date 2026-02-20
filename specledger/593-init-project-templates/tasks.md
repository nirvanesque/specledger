# Tasks Index: Project Template & Coding Agent Selection

**Feature**: 593-init-project-templates
**Status**: REVIEW - Tasks generated but NOT yet created in issue tracker

⚠️ **IMPORTANT**: This file shows the planned task structure. Review it first, then run the creation commands at the bottom.

## Feature Tracking

* **Epic ID**: TBD (will be created as `SL-xxxxxx`)
* **User Stories Source**: `specledger/593-init-project-templates/spec.md`
* **Research Inputs**: `specledger/593-init-project-templates/research.md`
* **Planning Details**: `specledger/593-init-project-templates/plan.md`
* **Data Model**: Not yet created (Phase 1 deliverable)
* **Contract Definitions**: Not yet created (Phase 1 deliverable)

## Epic Overview

**Title**: Project Template & Coding Agent Selection

**Description**: Enable developers to select from 7 business-defined project templates (General Purpose, Full-Stack, Batch Data, Real-Time Workflow, ML Image, Real-Time Data Pipeline, AI Chatbot) and choose their preferred coding agent (Claude Code, OpenCode, None) during interactive `sl new` command. System generates unique UUID v4 for each project to enable session storage and tracking.

**Labels**: `spec:593-init-project-templates`, `component:cli`, `component:tui`, `component:templates`

**Priority**: 1 (High - multiple P1 user stories)

---

## Task Organization Summary

### Total Tasks: 47 tasks across 10 phases

**Phase Breakdown**:
- Phase 1 (Setup): 3 tasks - Project initialization
- Phase 2 (Foundational): 6 tasks - Core infrastructure blocking all user stories
- Phase 3 (US1 - Template Selection): 7 tasks - Core template selection feature
- Phase 4 (US2 - Agent Configuration): 5 tasks - Agent selection and config
- Phase 5 (US3 - UUID Generation): 4 tasks - Project ID system
- Phase 6 (US6 - Backward Compatibility): 3 tasks - Migration and compatibility
- Phase 7 (US4 - Non-Interactive Mode): 4 tasks - CLI flags
- Phase 8 (US5 - List Templates): 2 tasks - Template discovery
- Phase 9 (US7 - Claude Settings): 3 tasks - Session capture integration
- Phase 10 (Polish): 10 tasks - Cross-cutting concerns

**User Story Priorities**:
- **P1 Stories** (MVP scope): US1, US2, US3, US6 - 19 implementation tasks
- **P2 Stories** (Enhancement): US4, US5, US7 - 9 implementation tasks
- **Foundation**: 9 infrastructure tasks (blocking all stories)
- **Polish**: 10 quality/documentation tasks

---

## CLI Commands Reference

After reviewing this file, execute these commands to create all issues:

```bash
# 1. Create Epic
sl issue create "Project Template & Coding Agent Selection" \
  --description "Enable template/agent selection in sl new with 7 templates, 3 agents, UUID generation" \
  --type epic \
  --labels "spec:593-init-project-templates,component:cli,component:tui" \
  --priority 1

# Note the Epic ID (SL-xxxxxx) and use it as --parent for all features below

# 2. Create Phase Features (replace EPIC_ID with actual ID)
# Phase 1: Setup
# Phase 2: Foundational
# Phase 3-9: User Stories
# Phase 10: Polish

# 3. Create all tasks under their respective phase features
# (See detailed commands in "Task Creation Commands" section below)
```

---

## Phase 1: Setup (Shared Infrastructure) 🔧

**Purpose**: Project initialization and dependency setup
**Feature ID**: TBD
**Labels**: `phase:setup`, `spec:593-init-project-templates`
**Dependencies**: None - can start immediately

### Tasks (3)

#### T001: Add github.com/google/uuid dependency
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:foundation`, `component:dependencies`, `spec:593-init-project-templates`
- **Description**: Need UUID generation capability for unique project IDs (FR-005). Current go.mod lacks UUID library.
- **Design**: Add `github.com/google/uuid v1.6.0` to go.mod using `go get`. Verify YAML marshaling works with uuid.UUID type.
- **Acceptance**: go.mod contains uuid dependency, go mod tidy succeeds, import works in test file
- **Definition of Done**:
  - [ ] go.mod contains github.com/google/uuid v1.6.0
  - [ ] go mod tidy completes without errors
  - [ ] UUID can be imported and used in Go code

#### T002: Create pkg/models directory structure
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `story:foundation`, `component:structure`, `spec:593-init-project-templates`
- **Description**: Need centralized location for data models (TemplateDefinition, AgentConfig). Currently no pkg/models/ directory.
- **Design**: Create `pkg/models/` directory. Add .gitkeep or README explaining purpose.
- **Acceptance**: pkg/models/ directory exists, can create .go files within it
- **Definition of Done**:
  - [ ] pkg/models/ directory created
  - [ ] Directory is tracked in git

#### T003: Update .gitignore for development artifacts
- **Type**: task
- **Priority**: 3 (low)
- **Labels**: `story:foundation`, `component:config`, `spec:593-init-project-templates`
- **Description**: Ensure build artifacts, test files, and IDE files are ignored during development.
- **Design**: Review .gitignore, add patterns: *.test, *.out, .DS_Store, .idea/, .vscode/. Preserve existing patterns.
- **Acceptance**: Common development artifacts don't appear in git status
- **Definition of Done**:
  - [ ] .gitignore includes test artifacts (*.test, *.out)
  - [ ] .gitignore includes IDE files (.idea/, .vscode/, .DS_Store)

**Checkpoint**: Setup complete - foundation tasks can begin

---

## Phase 2: Foundational (Blocking Prerequisites) ⚠️ CRITICAL

**Purpose**: Core infrastructure that MUST complete before ANY user story implementation
**Feature ID**: TBD
**Labels**: `phase:foundational`, `spec:593-init-project-templates`
**Dependencies**: Phase 1 (Setup) must complete first

**⚠️ BLOCKING**: No user story work can begin until ALL foundational tasks complete

### Tasks (6)

#### T004: Create TemplateDefinition model struct
- **Type**: task
- **Priority**: 0 (critical)
- **Labels**: `story:foundation`, `component:models`, `fr:FR-001`, `fr:FR-002`, `spec:593-init-project-templates`
- **Dependencies**: T001 (uuid dependency)
- **Description**: Need data structure to represent project templates with metadata (FR-001, FR-002). All template features depend on this type.
- **Design**: Create `pkg/models/template.go`. Define struct: ID (string), Name (string), Description (string), Characteristics ([]string), Path (string), IsDefault (bool). Add Validate() method checking ID format (kebab-case), name length (1-50 chars), description (1-200 chars), max 6 characteristics. Add String() method for logging.
- **Acceptance**: TemplateDefinition struct defined, Validate() rejects invalid templates, String() formats correctly
- **Definition of Done**:
  - [ ] TemplateDefinition struct created with all 6 fields
  - [ ] Validate() method rejects invalid ID formats
  - [ ] Validate() enforces name and description length limits
  - [ ] Unit tests for validation pass

#### T005: Create AgentConfig model struct
- **Type**: task
- **Priority**: 0 (critical)
- **Labels**: `story:foundation`, `component:models`, `fr:FR-004`, `spec:593-init-project-templates`
- **Description**: Need data structure for agent configurations (FR-004). Agent selection feature depends on this type.
- **Design**: Create `pkg/models/agent.go`. Define struct: ID (string), Name (string), Description (string), ConfigDir (string, e.g. ".claude"). Add Validate() method checking ID (kebab-case), name (1-50 chars), description (1-200 chars). Add HasConfig() bool method returning true if ConfigDir != "". Add SupportedAgents() []AgentConfig returning hardcoded list: claude-code, opencode, none. Add GetAgentByID(id string) (*AgentConfig, error) and DefaultAgent() AgentConfig.
- **Acceptance**: AgentConfig struct defined, SupportedAgents() returns 3 options, GetAgentByID() validates IDs
- **Definition of Done**:
  - [ ] AgentConfig struct created with all fields
  - [ ] SupportedAgents() returns 3 agents (Claude Code, OpenCode, None)
  - [ ] GetAgentByID() finds valid IDs and rejects invalid ones
  - [ ] HasConfig() correctly identifies agents with/without config directories

#### T006: Extend ProjectMetadata schema to v1.1.0
- **Type**: task
- **Priority**: 0 (critical)
- **Labels**: `story:foundation`, `component:metadata`, `fr:FR-006`, `fr:FR-007`, `fr:FR-008`, `fr:FR-009`, `spec:593-init-project-templates`
- **Dependencies**: T001 (uuid dependency)
- **Description**: Need to store project UUID, template, and agent in metadata (FR-006, FR-007, FR-008, FR-009). All user stories need this schema extension.
- **Design**: Modify `pkg/cli/metadata/schema.go`. Add to ProjectInfo: ID (uuid.UUID yaml:"id"), Template (string yaml:"template,omitempty"), Agent (string yaml:"agent,omitempty"). Update Validate() to check ID != uuid.Nil. Update NewProjectMetadata() to generate UUID with uuid.New(), set default template="", agent="". Bump MetadataVersion constant from "1.0.0" to "1.1.0".
- **Acceptance**: ProjectInfo has 3 new fields, Validate() requires UUID, new projects get auto-generated UUIDs
- **Definition of Done**:
  - [ ] ProjectInfo.ID field added (uuid.UUID type)
  - [ ] ProjectInfo.Template and ProjectInfo.Agent fields added
  - [ ] MetadataVersion constant updated to "1.1.0"
  - [ ] Validate() rejects nil UUIDs
  - [ ] NewProjectMetadata() generates UUID automatically

#### T007: Implement metadata migration v1.0.0 → v1.1.0
- **Type**: task
- **Priority**: 0 (critical)
- **Labels**: `story:US6`, `component:metadata`, `fr:FR-023`, `spec:593-init-project-templates`
- **Dependencies**: T006 (schema extension)
- **Description**: Need to migrate existing v1.0.0 projects without breaking them (FR-023). Backward compatibility requires automatic migration.
- **Design**: Create `pkg/cli/metadata/migration.go` or add to yaml.go. In Load() function, after unmarshal: if metadata.Version == "1.0.0" && metadata.Project.ID == uuid.Nil, generate new UUID, set metadata.Version = "1.1.0", save metadata. If Template or Agent fields empty, leave them empty (optional fields with omitempty).
- **Acceptance**: Loading v1.0.0 metadata auto-generates UUID, updates version, saves successfully
- **Definition of Done**:
  - [ ] Load() detects v1.0.0 projects without UUID
  - [ ] Auto-generates UUID for old projects
  - [ ] Updates version to v1.1.0 automatically
  - [ ] Saves migrated metadata without data loss

#### T008: Update manifest.yaml structure for templates
- **Type**: task
- **Priority**: 0 (critical)
- **Labels**: `story:foundation`, `component:templates`, `fr:FR-001`, `spec:593-init-project-templates`
- **Description**: Need manifest format supporting template metadata (FR-001). Template loading system depends on manifest structure.
- **Design**: Modify `pkg/embedded/templates/manifest.yaml`. Change structure from flat playbooks list to include template metadata. Add fields: id, name, description, characteristics[], path, is_default. Keep version at top level. Update embedded.go parsing if needed to handle new structure.
- **Acceptance**: Manifest parses with template metadata, LoadTemplates() returns TemplateDefinition objects
- **Definition of Done**:
  - [ ] manifest.yaml includes template metadata fields
  - [ ] Existing playbook loader can parse new format
  - [ ] Template characteristics stored as YAML array

#### T009: Create template loader in playbooks package
- **Type**: task
- **Priority**: 0 (critical)
- **Labels**: `story:foundation`, `component:templates`, `component:playbooks`, `spec:593-init-project-templates`
- **Dependencies**: T004 (TemplateDefinition), T008 (manifest)
- **Description**: Need function to load templates from embedded manifest. All template operations depend on this loader.
- **Design**: Add `LoadTemplates() ([]models.TemplateDefinition, error)` to `pkg/cli/playbooks/embedded.go` or new file `pkg/cli/playbooks/templates.go`. Read manifest.yaml from TemplatesFS, unmarshal to struct with Templates []TemplateDefinition field, validate each template, return slice. Cache results for performance.
- **Acceptance**: LoadTemplates() returns all templates from manifest, validates each, handles errors gracefully
- **Definition of Done**:
  - [ ] LoadTemplates() function implemented
  - [ ] Returns all templates from manifest.yaml
  - [ ] Validates each template on load
  - [ ] Caches results to avoid repeated reads

**Checkpoint**: Foundation complete - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Template Selection (Priority: P1) 🎯 MVP

**Goal**: Developer can interactively select from 7 project templates during `sl new`
**Feature ID**: TBD
**Labels**: `phase:us1`, `story:US1`, `spec:593-init-project-templates`
**Dependencies**: Phase 2 (Foundational) must complete

**Independent Test**: Run `sl new`, navigate template list with arrow keys, select any template, complete flow. Verify project created with correct template structure.

### Implementation Tasks (7)

#### T010: Add stepTemplate constant to TUI
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US1`, `component:tui`, `fr:FR-003`, `spec:593-init-project-templates`
- **Dependencies**: T009 (template loader)
- **Description**: Need new TUI step for template selection (FR-003). Must integrate into existing step sequence.
- **Design**: Modify `pkg/cli/tui/sl_new.go`. Add `stepTemplate` constant after stepShortCode, before stepPlaybook. Update step order: stepProjectName → stepDirectory → stepShortCode → stepTemplate → stepPlaybook → ... Renumber subsequent steps if needed.
- **Acceptance**: stepTemplate constant exists, step order correct, TUI flow progresses through new step
- **Definition of Done**:
  - [ ] stepTemplate constant added with correct value
  - [ ] Step inserted in correct position in flow
  - [ ] Subsequent steps renumbered if necessary

#### T011: Add template state to TUI Model struct
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US1`, `component:tui`, `fr:FR-001`, `spec:593-init-project-templates`
- **Dependencies**: T004 (TemplateDefinition), T010 (stepTemplate)
- **Description**: Need to store template selection state in TUI model (FR-001).
- **Design**: Modify `pkg/cli/tui/sl_new.go` Model struct. Add fields: `templates []models.TemplateDefinition` (loaded from manifest), `selectedTemplateIndex int` (current cursor position, default 0). Initialize templates in InitialModel() by calling playbooks.LoadTemplates(), set selectedTemplateIndex to default template index (find IsDefault=true or 0).
- **Acceptance**: Model stores template list and selection index, InitialModel() populates templates
- **Definition of Done**:
  - [ ] Model.templates field added
  - [ ] Model.selectedTemplateIndex field added
  - [ ] InitialModel() loads templates from manifest
  - [ ] Default template pre-selected

#### T012: Implement template selection View rendering
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US1`, `component:tui`, `fr:FR-002`, `spec:593-init-project-templates`
- **Dependencies**: T011 (template state)
- **Description**: Need to display template options with descriptions (FR-002, SC-001).
- **Design**: Add `viewTemplateSelection()` method to `pkg/cli/tui/sl_new.go`. For each template in m.templates: render cursor ("›" if selected), radio button ("◉" selected, "○" unselected), template name (bold if selected), description (subtle color), characteristics as "Tech: X, Y, Z" (subtle color). Add help text: "[↑/↓ to navigate, Enter to select]". Use lipgloss gold #13 for selected, gray #240 for unselected.
- **Acceptance**: Template list displays with names, descriptions, tech tags, correct styling, cursor navigation
- **Definition of Done**:
  - [ ] All 7 templates displayed with names and descriptions
  - [ ] Selected template highlighted with cursor and radio button
  - [ ] Technology characteristics shown for each template
  - [ ] Help text displayed at bottom

#### T013: Implement template selection Update navigation
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US1`, `component:tui`, `fr:FR-003`, `spec:593-init-project-templates`
- **Dependencies**: T011 (template state)
- **Description**: Need arrow key navigation through templates (FR-003, US1 acceptance 2).
- **Design**: In Update() method switch on m.step, add case stepTemplate. Handle tea.KeyMsg: "up"/"k" decrements selectedTemplateIndex with wraparound to len(templates)-1 if <0. "down"/"j" increments with wraparound to 0 if >= len(templates). Return updated model.
- **Acceptance**: Up/down arrows move selection, wraps at boundaries, visual feedback immediate
- **Definition of Done**:
  - [ ] Up arrow decrements selection with wraparound
  - [ ] Down arrow increments selection with wraparound
  - [ ] Selection updates immediately in view
  - [ ] Keyboard shortcuts (k/j) work as alternatives

#### T014: Implement template selection confirmation
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US1`, `component:tui`, `fr:FR-003`, `fr:FR-007`, `spec:593-init-project-templates`
- **Dependencies**: T013 (navigation)
- **Description**: Need Enter key to confirm template selection and store in answers map (FR-007, US1 acceptance 3).
- **Design**: In Update() stepTemplate case, handle tea.KeyEnter. Store templates[selectedTemplateIndex].ID in m.answers["template"]. Validate selection index in bounds. Advance to next step (m.step = stepAgent or stepPlaybook). Log selection.
- **Acceptance**: Enter stores template ID in answers, advances to next step, invalid index handled safely
- **Definition of Done**:
  - [ ] Enter key stores template ID in answers map
  - [ ] Flow advances to next step (agent selection)
  - [ ] Selection validated before storing
  - [ ] Template ID (e.g., "full-stack") stored correctly

#### T015: Add template selection to confirmation review
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `story:US1`, `component:tui`, `spec:593-init-project-templates`
- **Dependencies**: T014 (confirmation)
- **Description**: Need to display selected template in final confirmation screen before project creation.
- **Design**: Modify confirmation view in `sl_new.go`. Add line showing "Template: {template name from ID}" between short code and playbook/agent. Look up template name from m.templates using m.answers["template"] ID.
- **Acceptance**: Confirmation screen shows template selection, name displayed (not just ID)
- **Definition of Done**:
  - [ ] Confirmation view displays selected template name
  - [ ] Template shown in correct position in summary
  - [ ] Template name resolved from ID correctly

#### T016: Create 7 template directory structures
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US1`, `component:templates`, `fr:FR-010`, `spec:593-init-project-templates`
- **Dependencies**: T008 (manifest structure)
- **Description**: Need actual template directories with starter files (FR-010, SC-002). Templates must represent industry-standard structures.
- **Design**: Create directories in `pkg/embedded/templates/`: general-purpose/ (copy current specledger), full-stack/ (backend/ + frontend/), batch-data/ (workflows/ + cmd/), realtime-workflow/ (workflows/ + activities/), ml-image/ (src/data + src/models), realtime-data/ (cmd/producer + cmd/consumer + internal/kafka), ai-chatbot/ (src/agents + src/integrations). Each includes: README.md, .gitignore, starter files, directory structure per research.md. Update manifest.yaml with all 7 definitions.
- **Acceptance**: All 7 template directories exist with complete structure, manifest lists all, embedded correctly
- **Definition of Done**:
  - [ ] All 7 template directories created with structures from research.md
  - [ ] Each template has README.md explaining structure
  - [ ] manifest.yaml includes all 7 template definitions
  - [ ] Templates embed correctly in binary (test with go:embed)

#### T017: Add structured logging for template selection
- **Type**: task
- **Priority**: 3 (low)
- **Labels**: `story:US1`, `component:logging`, `fr:FR-028`, `spec:593-init-project-templates`
- **Dependencies**: T014 (confirmation)
- **Description**: Need logging for template operations (FR-028, SC-012). Debugging and monitoring require logs.
- **Design**: Add log/slog logging in template selection flow. Log: template list loaded (count), template selected (ID, name), template validation (success/failure). Use structured fields: slog.Info("template selected", "id", templateID, "name", templateName).
- **Acceptance**: Template operations logged with structured fields, logs include relevant context
- **Definition of Done**:
  - [ ] Template list loading logged
  - [ ] Template selection logged with ID and name
  - [ ] Structured logging fields used (not string formatting)

**Checkpoint**: Template selection complete and independently testable

---

## Phase 4: User Story 2 - Agent Configuration (Priority: P1)

**Goal**: Developer can select coding agent (Claude Code, OpenCode, None) with correct config files created
**Feature ID**: TBD
**Labels**: `phase:us2`, `story:US2`, `spec:593-init-project-templates`
**Dependencies**: Phase 2 (Foundational)

**Independent Test**: Run `sl new`, select any template, select each agent option, verify correct config directories created (.claude/, .opencode/, or none).

### Implementation Tasks (5)

#### T018: Add stepAgent constant to TUI
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US2`, `component:tui`, `fr:FR-004`, `spec:593-init-project-templates`
- **Dependencies**: T009 (foundational complete)
- **Description**: Need TUI step for agent selection (FR-004, US2 acceptance 1).
- **Design**: Modify `pkg/cli/tui/sl_new.go`. Add stepAgent constant after stepTemplate. Update Model struct: add `agents []models.AgentConfig`, `selectedAgentIndex int`. Initialize in InitialModel(): agents = models.SupportedAgents(), selectedAgentIndex = 0 (Claude Code default).
- **Acceptance**: stepAgent constant exists, agent state in Model, agents list populated with 3 options
- **Definition of Done**:
  - [ ] stepAgent constant added after stepTemplate
  - [ ] Model.agents and Model.selectedAgentIndex fields added
  - [ ] InitialModel() loads 3 agents (Claude Code, OpenCode, None)
  - [ ] Claude Code pre-selected as default

#### T019: Implement agent selection View and Update
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US2`, `component:tui`, `fr:FR-004`, `spec:593-init-project-templates`
- **Dependencies**: T018 (stepAgent)
- **Description**: Need agent selection UI with navigation (US2 acceptance 2).
- **Design**: Add viewAgentSelection() rendering 3 agents with cursor, radio buttons, names, descriptions (similar to template view). In Update() case stepAgent: handle up/down keys with wraparound through 3 agents. Enter stores agents[selectedAgentIndex].ID in m.answers["agent"], advances step.
- **Acceptance**: 3 agents displayed, navigation works, Enter confirms selection
- **Definition of Done**:
  - [ ] All 3 agents displayed with names and descriptions
  - [ ] Arrow key navigation with wraparound
  - [ ] Enter stores agent ID and advances flow

#### T020: Implement agent config directory creation
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US2`, `component:bootstrap`, `fr:FR-011`, `fr:FR-012`, `fr:FR-015`, `spec:593-init-project-templates`
- **Dependencies**: T019 (agent selection)
- **Description**: Need to create agent-specific directories based on selection (FR-011, FR-012, FR-015, US2 acceptance 3, 4, 5).
- **Design**: Modify `pkg/cli/commands/bootstrap.go` or bootstrap_helpers.go. After project directory created, read m.answers["agent"]. Get agent via models.GetAgentByID(). If agent.HasConfig() true, create agent.ConfigDir directory (e.g., ".claude/" or ".opencode/"). Copy agent template files from embedded templates. If agent.ConfigDir empty (none), skip creation. Handle errors as non-fatal (log warning, continue).
- **Acceptance**: Claude Code creates .claude/ with commands/ and skills/, OpenCode creates .opencode/, None creates nothing
- **Definition of Done**:
  - [ ] .claude/ directory created when Claude Code selected
  - [ ] .opencode/ directory created when OpenCode selected
  - [ ] No agent directories when None selected
  - [ ] Agent config files copied correctly

#### T021: Create OpenCode config template structure
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `story:US2`, `component:templates`, `fr:FR-012`, `spec:593-init-project-templates`
- **Dependencies**: T016 (template directories)
- **Description**: Need OpenCode configuration files to copy (FR-012, US2 acceptance 4). OpenCode requires opencode.json and directory structure.
- **Design**: Create `pkg/embedded/templates/agents/opencode/` directory with: .opencode/commands/ (port from .claude/commands/), .opencode/skills/ (port from .claude/skills/), opencode.json (JSON schema reference), AGENTS.md template. Copy logic reads from agents/opencode/ when OpenCode selected.
- **Acceptance**: OpenCode template exists with all required files, copies correctly to project
- **Definition of Done**:
  - [ ] agents/opencode/ directory created with structure
  - [ ] .opencode/commands/ ported from .claude/commands/
  - [ ] opencode.json template created with schema reference
  - [ ] AGENTS.md template created

#### T022: Store agent selection in metadata
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US2`, `component:metadata`, `fr:FR-008`, `spec:593-init-project-templates`
- **Dependencies**: T019 (agent selection), T006 (metadata schema)
- **Description**: Need to record agent choice in specledger.yaml (FR-008, SC-007).
- **Design**: In bootstrap.go, after reading m.answers["agent"], set metadata.Project.Agent = agentID before saving. Validate agentID against known agents (claude-code, opencode, none).
- **Acceptance**: specledger.yaml contains correct agent field, validates known agents only
- **Definition of Done**:
  - [ ] Agent ID stored in metadata.Project.Agent
  - [ ] Metadata validates agent ID
  - [ ] Agent field appears in specledger.yaml

**Checkpoint**: Agent configuration complete and independently testable

---

## Phase 5: User Story 3 - UUID Generation (Priority: P1)

**Goal**: Every new project receives unique UUID v4 stored in metadata
**Feature ID**: TBD
**Labels**: `phase:us3`, `story:US3`, `spec:593-init-project-templates`
**Dependencies**: Phase 2 (Foundational)

**Independent Test**: Create multiple projects, verify each has unique UUID in specledger.yaml project.id field.

### Implementation Tasks (4)

#### T023: Generate UUID in NewProjectMetadata
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US3`, `component:metadata`, `fr:FR-005`, `fr:FR-006`, `spec:593-init-project-templates`
- **Dependencies**: T006 (metadata schema)
- **Description**: Need automatic UUID generation for new projects (FR-005, FR-006, US3 acceptance 1).
- **Design**: Modify `pkg/cli/metadata/yaml.go` NewProjectMetadata() function. Add line: `Project: ProjectInfo{ID: uuid.New(), ...}`. UUID generated via github.com/google/uuid. If uuid.New() panics (extremely rare), let it fail fast.
- **Acceptance**: New projects get UUID automatically, UUID is cryptographically random, no collisions in 1000 projects
- **Definition of Done**:
  - [ ] NewProjectMetadata() calls uuid.New()
  - [ ] UUID stored in ProjectInfo.ID field
  - [ ] UUID appears in saved specledger.yaml

#### T024: Add UUID validation
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `story:US3`, `component:metadata`, `spec:593-init-project-templates`
- **Dependencies**: T023 (UUID generation)
- **Description**: Need to validate UUID presence on metadata load (US3 acceptance 1, SC-005).
- **Design**: In metadata Validate() method, add check: if m.Version == "1.1.0" && m.Project.ID == uuid.Nil, return error "project ID required for v1.1.0". Allow nil UUID for v1.0.0 (migration handles it).
- **Acceptance**: v1.1.0 metadata requires UUID, validation rejects nil UUIDs, v1.0.0 allowed without UUID
- **Definition of Done**:
  - [ ] Validate() rejects v1.1.0 metadata with nil UUID
  - [ ] Error message clear and actionable
  - [ ] v1.0.0 projects still load (for migration)

#### T025: Verify UUID uniqueness
- **Type**: task
- **Priority**: 3 (low)
- **Labels**: `story:US3`, `component:testing`, `spec:593-init-project-templates`
- **Dependencies**: T023 (UUID generation)
- **Description**: Need test ensuring UUID collision probability is negligible (SC-005).
- **Design**: Create test in metadata package: generate 10,000 UUIDs with uuid.New(), store in map, verify no duplicates. Test UUID YAML marshaling/unmarshaling (round-trip preserves value).
- **Acceptance**: 10,000 UUIDs all unique, YAML round-trip preserves UUID value
- **Definition of Done**:
  - [ ] Test generates 10,000 UUIDs without collisions
  - [ ] UUID marshals to YAML correctly (hyphenated format)
  - [ ] UUID unmarshals from YAML correctly

#### T026: Document UUID purpose in metadata schema
- **Type**: task
- **Priority**: 3 (low)
- **Labels**: `story:US3`, `component:documentation`, `spec:593-init-project-templates`
- **Dependencies**: T023 (UUID generation)
- **Description**: Need documentation explaining UUID field purpose (US3 acceptance 3).
- **Design**: Add Go doc comment to ProjectInfo.ID field: "ID is a unique project identifier (UUID v4) used to organize sessions in Supabase storage. Auto-generated on project creation, never changes." Update schema.go package docs if present.
- **Acceptance**: ID field has clear doc comment explaining purpose and lifecycle
- **Definition of Done**:
  - [ ] ProjectInfo.ID field has Go doc comment
  - [ ] Comment explains UUID purpose (session tracking)
  - [ ] Comment states UUID auto-generated and immutable

**Checkpoint**: UUID generation complete and independently testable

---

## Phase 6: User Story 6 - Backward Compatibility (Priority: P1)

**Goal**: Default selections produce identical output to current sl new, old projects migrate smoothly
**Feature ID**: TBD
**Labels**: `phase:us6`, `story:US6`, `spec:593-init-project-templates`
**Dependencies**: Phase 2 (Foundational), Phase 3 (US1), Phase 4 (US2), Phase 5 (US3)

**Independent Test**: Create project with defaults (General Purpose + Claude Code), compare to current sl new output byte-for-byte (except metadata). Load old v1.0.0 project, verify UUID added and version updated.

### Implementation Tasks (3)

#### T027: Set default template to general-purpose
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US6`, `component:tui`, `component:templates`, `fr:FR-022`, `spec:593-init-project-templates`
- **Dependencies**: T011 (template state), T016 (templates created)
- **Description**: Need General Purpose template as default (FR-022, US6 acceptance 1, SC-003).
- **Design**: In manifest.yaml, set is_default: true for general-purpose template. In InitialModel(), find default template: for i, t := range templates { if t.IsDefault { selectedTemplateIndex = i; break } }. Fallback to index 0 if no default found.
- **Acceptance**: General Purpose pre-selected, matches current sl new structure exactly
- **Definition of Done**:
  - [ ] general-purpose template has is_default: true in manifest
  - [ ] InitialModel() selects default template automatically
  - [ ] General Purpose structure identical to current specledger playbook

#### T028: Test backward compatibility of default flow
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US6`, `component:testing`, `fr:FR-022`, `spec:593-init-project-templates`
- **Dependencies**: T027 (default template), T018 (default agent)
- **Description**: Need automated test verifying defaults match current behavior (SC-003, US6 acceptance 1).
- **Design**: Create integration test in tests/integration/: create project with all defaults (press Enter through TUI, or use flags with general-purpose + claude-code), compare directory structure to baseline snapshot (current sl new output). Ignore specledger.yaml differences (new fields). Use go-cmp or similar for recursive comparison.
- **Acceptance**: Test passes, directory structure identical to current sl new (excluding metadata)
- **Definition of Done**:
  - [ ] Integration test compares default output to baseline
  - [ ] Test ignores new metadata fields (id, template, agent)
  - [ ] Directory structure and files byte-for-byte identical

#### T029: Implement and test v1.0.0 → v1.1.0 migration
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US6`, `component:metadata`, `fr:FR-023`, `spec:593-init-project-templates`
- **Dependencies**: T007 (migration logic), T024 (UUID validation)
- **Description**: Need to verify old projects migrate correctly (FR-023, US6 acceptance 2, 3, SC-009).
- **Design**: Create integration test: save v1.0.0 specledger.yaml (no id/template/agent), call Load(), verify: ID generated (not nil), Version updated to "1.1.0", Template and Agent empty, all original fields preserved. Test Save() writes correct format. Test migration idempotent (loading twice doesn't change UUID).
- **Acceptance**: v1.0.0 projects load, get UUID, update to v1.1.0, no data loss, idempotent
- **Definition of Done**:
  - [ ] Test loads v1.0.0 project successfully
  - [ ] UUID auto-generated on load
  - [ ] Version updated to v1.1.0
  - [ ] Original project data preserved
  - [ ] Migration idempotent (same UUID on repeat loads)

**Checkpoint**: Backward compatibility verified and tested

---

## Phase 7: User Story 4 - Non-Interactive Mode (Priority: P2)

**Goal**: Developers can use --template and --agent flags for CI/CD automation
**Feature ID**: TBD
**Labels**: `phase:us4`, `story:US4`, `spec:593-init-project-templates`
**Dependencies**: Phase 2 (Foundational), Phase 3 (US1), Phase 4 (US2)

**Independent Test**: Run `sl new --project-name test --project-dir /tmp/test --short-code t --template full-stack --agent opencode`, verify project created without prompts.

### Implementation Tasks (4)

#### T030: Add --template and --agent CLI flags
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US4`, `component:cli`, `fr:FR-016`, `fr:FR-017`, `spec:593-init-project-templates`
- **Dependencies**: None (can be done early)
- **Description**: Need CLI flags for non-interactive mode (FR-016, FR-017, US4 acceptance 1).
- **Design**: Modify `sl new` command definition (pkg/cli/commands/ or cmd/). Add flags: --template string (template ID, e.g., "full-stack"), --agent string (agent ID, e.g., "claude-code"). Register with cobra: cmd.Flags().String("template", "", "Project template ID"), cmd.Flags().String("agent", "", "Coding agent ID").
- **Acceptance**: Flags accepted by CLI, values accessible in command handler
- **Definition of Done**:
  - [ ] --template flag added with help text
  - [ ] --agent flag added with help text
  - [ ] Flags registered with cobra CLI framework

#### T031: Implement TTY detection and flag validation
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US4`, `component:cli`, `fr:FR-021`, `spec:593-init-project-templates`
- **Dependencies**: T030 (flags added)
- **Description**: Need to require flags in non-TTY environment (FR-021, US4 acceptance 1).
- **Design**: In `sl new` command handler, use isatty or similar library to check if stdin is TTY. If !isatty && (template flag empty || agent flag empty), return error: "non-interactive mode requires --template and --agent flags". If TTY, run TUI as normal (flags optional, override TUI defaults).
- **Acceptance**: Non-TTY requires flags, TTY allows interactive, error messages clear
- **Definition of Done**:
  - [ ] TTY detection implemented (isatty or terminal.IsTerminal)
  - [ ] Non-TTY environment requires both flags
  - [ ] Clear error message when flags missing
  - [ ] TTY mode still works interactively

#### T032: Validate template and agent flag values
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US4`, `component:cli`, `fr:FR-019`, `fr:FR-020`, `spec:593-init-project-templates`
- **Dependencies**: T031 (flag validation)
- **Description**: Need to reject invalid template/agent IDs (FR-019, FR-020, US4 acceptance 2).
- **Design**: If flags provided, validate before creating project. Load templates via LoadTemplates(), check if flag value matches any template.ID. Check agent via models.GetAgentByID(). If invalid, return error listing valid options: "unknown template: X. Available: general-purpose, full-stack, ..." Include help text: "use --list-templates for details".
- **Acceptance**: Invalid IDs rejected with helpful error, valid options listed, exit code non-zero
- **Definition of Done**:
  - [ ] Invalid template ID shows error with available options
  - [ ] Invalid agent ID shows error with available options
  - [ ] Command exits with non-zero code on invalid input
  - [ ] Error messages reference --list-templates for more info

#### T033: Bypass TUI when flags provided
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `story:US4`, `component:cli`, `component:tui`, `spec:593-init-project-templates`
- **Dependencies**: T032 (validation)
- **Description**: Need to skip TUI and use flag values directly (US4 acceptance 3, SC-004).
- **Design**: In command handler, if template and agent flags provided and valid, build answers map directly: answers["template"] = templateFlag, answers["agent"] = agentFlag. Pass to bootstrap without running TUI. Still prompt for name/dir/shortcode if missing (or require those flags too in non-TTY).
- **Acceptance**: Flags bypass TUI, project created immediately, output matches interactive mode
- **Definition of Done**:
  - [ ] Flags bypass TUI when provided
  - [ ] Answers map populated from flags
  - [ ] Bootstrap creates project correctly
  - [ ] Output identical to interactive mode with same selections

**Checkpoint**: Non-interactive mode complete and tested

---

## Phase 8: User Story 5 - List Templates (Priority: P2)

**Goal**: Developers can discover available templates with --list-templates flag
**Feature ID**: TBD
**Labels**: `phase:us5`, `story:US5`, `spec:593-init-project-templates`
**Dependencies**: Phase 2 (Foundational)

**Independent Test**: Run `sl new --list-templates`, verify all 7 templates displayed with descriptions and tech tags.

### Implementation Tasks (2)

#### T034: Add --list-templates flag
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `story:US5`, `component:cli`, `fr:FR-018`, `spec:593-init-project-templates`
- **Dependencies**: T009 (template loader)
- **Description**: Need discovery flag for available templates (FR-018, US5 acceptance 1).
- **Design**: Add flag to `sl new`: cmd.Flags().Bool("list-templates", false, "List available project templates"). In command handler, if flag true, call LoadTemplates(), display formatted list, exit 0.
- **Acceptance**: Flag recognized, triggers template listing, exits without creating project
- **Definition of Done**:
  - [ ] --list-templates flag added
  - [ ] Flag triggers template listing
  - [ ] Command exits cleanly after listing

#### T035: Format and display template list
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `story:US5`, `component:cli`, `fr:FR-018`, `spec:593-init-project-templates`
- **Dependencies**: T034 (flag added)
- **Description**: Need formatted output showing template details (US5 acceptance 2, 3, SC-010).
- **Design**: When --list-templates active, format output: "Available Project Templates:\n\n  {id}{default marker}\n    {description}\n    Tech: {characteristics}\n\n". Mark default with " (default)" suffix. Add footer: "Use: sl new --template <id>". Measure performance (should be <100ms per SC-010).
- **Acceptance**: All 7 templates listed, default marked, descriptions shown, tech tags included, exits 0
- **Definition of Done**:
  - [ ] Output shows all 7 templates with IDs and names
  - [ ] Descriptions displayed for each template
  - [ ] Technology tags shown (e.g., "Tech: Go, React, PostgreSQL")
  - [ ] Default template marked
  - [ ] Usage instructions at bottom
  - [ ] Performance <100ms (SC-010)

**Checkpoint**: Template discovery complete

---

## Phase 9: User Story 7 - Claude Code Settings (Priority: P2)

**Goal**: Claude Code selection auto-configures session capture with project UUID
**Feature ID**: TBD
**Labels**: `phase:us7`, `story:US7`, `spec:593-init-project-templates`
**Dependencies**: Phase 4 (US2), Phase 5 (US3)

**Independent Test**: Create project with Claude Code agent, read .claude/settings.json, verify PostToolUse hook with `sl session capture` command and correct structure.

### Implementation Tasks (3)

#### T036: Create .claude/settings.json template
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `story:US7`, `component:templates`, `fr:FR-013`, `spec:593-init-project-templates`
- **Dependencies**: T021 (agent templates)
- **Description**: Need settings template with session capture config (FR-013, US7 acceptance 1, 2).
- **Design**: Create `pkg/embedded/templates/agents/claude/settings.json` with structure: `{"saveTranscripts": true, "transcriptsDirectory": "~/.claude/sessions", "hooks": {"PostToolUse": [{"matcher": "Bash", "hooks": [{"type": "command", "command": "sl session capture"}]}]}}`. This is JSON template copied when Claude Code selected.
- **Acceptance**: Template JSON valid, includes session capture hook, copies to .claude/ directory
- **Definition of Done**:
  - [ ] settings.json template created in agents/claude/
  - [ ] JSON structure includes PostToolUse hook
  - [ ] Hook triggers on Bash tool usage
  - [ ] Command is "sl session capture"

#### T037: Copy settings.json when Claude Code selected
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `story:US7`, `component:bootstrap`, `fr:FR-013`, `fr:FR-014`, `spec:593-init-project-templates`
- **Dependencies**: T036 (settings template), T020 (agent creation)
- **Description**: Need to install settings.json during project creation (FR-014, US7 acceptance 3, SC-008).
- **Design**: In bootstrap agent config creation (T020 code), if agent.ID == "claude-code", copy settings.json from embedded template to projectDir/.claude/settings.json. Set file permissions 0644.
- **Acceptance**: settings.json appears in .claude/ when Claude Code selected, not present for other agents
- **Definition of Done**:
  - [ ] settings.json copied to .claude/ for Claude Code
  - [ ] File not created for OpenCode or None
  - [ ] File permissions set to 0644
  - [ ] File readable by Claude Code CLI

#### T038: Document session capture integration
- **Type**: task
- **Priority**: 3 (low)
- **Labels**: `story:US7`, `component:documentation`, `spec:593-init-project-templates`
- **Dependencies**: T037 (settings.json copying)
- **Description**: Need documentation explaining session capture auto-configuration (US7 acceptance 3).
- **Design**: Add comment to settings.json template explaining hook purpose. Create or update template README.md files explaining session capture integration: "This project uses Claude Code with automatic session capture. All terminal interactions are saved to ~/.claude/sessions and uploaded to Supabase for project tracking."
- **Acceptance**: Documentation explains session capture, developers understand auto-configuration
- **Definition of Done**:
  - [ ] settings.json includes comment explaining hook
  - [ ] Template README mentions session capture
  - [ ] Instructions clear for developers unfamiliar with feature

**Checkpoint**: Claude Code integration complete

---

## Phase 10: Polish & Cross-Cutting Concerns

**Purpose**: Quality, documentation, and cross-cutting improvements
**Feature ID**: TBD
**Labels**: `phase:polish`, `spec:593-init-project-templates`
**Dependencies**: All user story phases (3-9) complete

### Tasks (10)

#### T039: Create integration tests for all template types
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `component:testing`, `spec:593-init-project-templates`
- **Dependencies**: T016 (templates created)
- **Description**: Need test coverage for all 7 template creations (SC-002, SC-006).
- **Design**: Create `tests/integration/templates_test.go`. For each template ID: call project creation, verify key directories exist (specific to template per research.md), verify starter files present, verify metadata correct. Table-driven test with template ID → expected directories map.
- **Acceptance**: All 7 templates tested, key directories verified, tests pass
- **Definition of Done**:
  - [ ] Integration test covers all 7 templates
  - [ ] Each template verified for key directories
  - [ ] Test uses table-driven approach
  - [ ] All tests pass

#### T040: Create integration tests for all agent configs
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `component:testing`, `spec:593-init-project-templates`
- **Dependencies**: T020 (agent creation), T021 (OpenCode template)
- **Description**: Need test coverage for all 3 agent configurations (SC-007).
- **Design**: Create `tests/integration/agents_test.go`. For each agent: create project, verify correct directories (.claude/, .opencode/, or none), verify config files present, verify settings.json for Claude Code. Table-driven test with agent ID → expected directories.
- **Acceptance**: All 3 agents tested, config directories verified, tests pass
- **Definition of Done**:
  - [ ] Integration test covers all 3 agents
  - [ ] Claude Code verified for .claude/ and settings.json
  - [ ] OpenCode verified for .opencode/
  - [ ] None verified for no agent directories

#### T041: Add TUI navigation tests
- **Type**: task
- **Priority**: 3 (low)
- **Labels**: `component:testing`, `component:tui`, `spec:593-init-project-templates`
- **Dependencies**: T013 (template navigation), T019 (agent navigation)
- **Description**: Need unit tests for TUI step navigation (SC-001).
- **Design**: Create `pkg/cli/tui/sl_new_test.go` or extend existing. Test: arrow key navigation moves selection, wraparound works, Enter confirms, answers map populated. Test Model.Update() directly without running full tea.Program.
- **Acceptance**: Navigation tests pass, wraparound verified, selection confirmed
- **Definition of Done**:
  - [ ] Tests for template selection navigation
  - [ ] Tests for agent selection navigation
  - [ ] Wraparound behavior tested
  - [ ] Selection confirmation tested

#### T042: Add --force flag documentation and implementation
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `component:cli`, `fr:FR-026`, `spec:593-init-project-templates`
- **Dependencies**: None (can be done early)
- **Description**: Need --force flag for directory overwrite (edge case, clarification answer 1, FR-026).
- **Design**: Add flag: cmd.Flags().Bool("force", false, "Overwrite existing directory"). In bootstrap, before creating project directory, check if exists. If exists && !force, return error. If exists && force, delete directory recursively (os.RemoveAll), then create fresh.
- **Acceptance**: --force deletes and recreates, without flag shows error if directory exists
- **Definition of Done**:
  - [ ] --force flag added to sl new
  - [ ] Existing directory detected and rejected without flag
  - [ ] --force deletes existing directory completely
  - [ ] Project created fresh after deletion

#### T043: Enforce template immutability
- **Type**: task
- **Priority**: 3 (low)
- **Labels**: `component:cli`, `fr:FR-027`, `spec:593-init-project-templates`
- **Dependencies**: T022 (template stored in metadata)
- **Description**: Need to prevent template changes after creation (clarification answer 2, FR-027, SC-011).
- **Design**: If attempting to add "sl template change" command in future, return error. Document in metadata that Template field is immutable. No code change needed now, just documentation and test that changing metadata.Project.Template manually has no effect on project structure.
- **Acceptance**: Documentation states template immutable, no template change command exists
- **Definition of Done**:
  - [ ] Comment in schema.go states Template field immutable
  - [ ] No command to change template exists
  - [ ] README or docs mention template immutability

#### T044: Add error handling for agent config copy failures
- **Type**: task
- **Priority**: 2 (normal)
- **Labels**: `component:bootstrap`, `spec:593-init-project-templates`
- **Dependencies**: T020 (agent config creation)
- **Description**: Need graceful handling of agent config failures (edge case from spec).
- **Design**: In agent config copying code (T020), wrap file operations in error handling. If copy fails, log warning: slog.Warn("failed to create agent config", "agent", agentID, "error", err). Continue project creation (non-fatal). Final output shows warning but reports success.
- **Acceptance**: Agent config failures logged but don't abort project creation
- **Definition of Done**:
  - [ ] Agent config copy failures logged as warnings
  - [ ] Project creation continues despite agent failures
  - [ ] User sees warning in output

#### T045: Update AGENTS.md template for all templates
- **Type**: task
- **Priority**: 3 (low)
- **Labels**: `component:templates`, `component:documentation`, `spec:593-init-project-templates`
- **Dependencies**: T016 (templates created)
- **Description**: Need AGENTS.md context file for each template explaining structure and technologies.
- **Design**: For each of 7 templates, create or update AGENTS.md file with: project description, technology stack, directory structure explanation, development commands, testing approach. Use consistent format across all templates. Include markers for manual additions: <!-- MANUAL ADDITIONS START --> <!-- MANUAL ADDITIONS END -->.
- **Acceptance**: All 7 templates have AGENTS.md, content explains structure, consistent format
- **Definition of Done**:
  - [ ] AGENTS.md created for all 7 templates
  - [ ] Each explains template-specific structure
  - [ ] Technology stack documented per template
  - [ ] Manual addition markers present

#### T046: Update command help text for new flags
- **Type**: task
- **Priority**: 3 (low)
- **Labels**: `component:cli`, `component:documentation`, `spec:593-init-project-templates`
- **Dependencies**: T030 (flags added), T034 (list-templates flag)
- **Description**: Need clear help text for all new flags.
- **Design**: Update `sl new` command help text (cobra Long/Short descriptions). Explain --template, --agent, --list-templates, --force flags. Add examples: "sl new --template full-stack --agent opencode". Document non-interactive mode requirements.
- **Acceptance**: sl new --help shows clear flag descriptions with examples
- **Definition of Done**:
  - [ ] Help text documents all new flags
  - [ ] Examples provided for common usage
  - [ ] Non-interactive mode requirements explained

#### T047: Create quickstart.md developer guide
- **Type**: task
- **Priority**: 3 (low)
- **Labels**: `component:documentation`, `spec:593-init-project-templates`
- **Dependencies**: All implementation complete
- **Description**: Need developer guide for extending templates and TUI (plan.md Phase 1 deliverable).
- **Design**: Create `specledger/593-init-project-templates/quickstart.md` with: how to add new template (directory, manifest, testing), how to add new TUI step, how to test template creation, common patterns, troubleshooting. Include code examples.
- **Acceptance**: Quickstart.md provides step-by-step guide for common developer tasks
- **Definition of Done**:
  - [ ] quickstart.md created with all sections
  - [ ] Instructions for adding new template
  - [ ] Instructions for extending TUI
  - [ ] Code examples included
  - [ ] Troubleshooting section present

#### T048: Run final integration test suite
- **Type**: task
- **Priority**: 1 (high)
- **Labels**: `component:testing`, `spec:593-init-project-templates`
- **Dependencies**: All implementation tasks complete
- **Description**: Need comprehensive test run verifying all success criteria (SC-001 through SC-012).
- **Design**: Run full test suite: go test ./... -v. Verify: TUI flow <60s, all templates create correctly, backward compatibility, non-interactive mode, UUID uniqueness, agent configs, template immutability, logging. Create test script automating all success criteria verification.
- **Acceptance**: All tests pass, all success criteria verified, no regressions
- **Definition of Done**:
  - [ ] Full test suite runs without failures
  - [ ] All 12 success criteria verified
  - [ ] Test script created for repeatable verification
  - [ ] Performance benchmarks met (SC-001, SC-010)

**Checkpoint**: Feature complete, all success criteria met, ready for production

---

## Dependencies & Execution Order

### Critical Path (MVP - Must Complete in Order)

```
Phase 1 (Setup) → Phase 2 (Foundational) → [Phase 3 (US1), Phase 4 (US2), Phase 5 (US3)] → Phase 6 (US6)
```

**Explanation**:
- Phase 1 (Setup): Must complete first - establishes dependencies and structure
- Phase 2 (Foundational): BLOCKS all user stories - creates core infrastructure
- Phases 3-5 (P1 User Stories): Can execute in parallel after Phase 2 completes
- Phase 6 (US6): Requires Phases 3-5 for backward compatibility testing

### Secondary Features (Can Defer)

```
Phase 7 (US4 - Non-Interactive): Can start after Phases 3-4 complete
Phase 8 (US5 - List Templates): Can start after Phase 2 complete
Phase 9 (US7 - Claude Settings): Can start after Phases 4-5 complete
```

### Parallel Execution Opportunities

**After Phase 2 completes**, these can run in parallel with proper team coordination:
- **Team A**: Phase 3 (US1 - Template Selection) - 7 tasks
- **Team B**: Phase 4 (US2 - Agent Configuration) - 5 tasks
- **Team C**: Phase 5 (US3 - UUID Generation) - 4 tasks

**After P1 stories complete**, these can run in parallel:
- **Team A**: Phase 7 (US4 - Non-Interactive Mode) - 4 tasks
- **Team B**: Phase 8 (US5 - List Templates) - 2 tasks
- **Team C**: Phase 9 (US7 - Claude Settings) - 3 tasks

### MVP Recommendation

**Minimum Viable Product** = Phases 1, 2, 3, 4, 5, 6 (19 implementation tasks)

This delivers:
- ✅ Template selection (US1)
- ✅ Agent configuration (US2)
- ✅ UUID generation (US3)
- ✅ Backward compatibility (US6)
- ✅ Core value proposition functional
- ❌ Non-interactive mode (defer to post-MVP)
- ❌ Template listing (defer to post-MVP)
- ❌ Claude settings automation (defer to post-MVP)

**Estimated Effort**: ~2 weeks with 1 developer, ~1 week with 2-3 developers (parallel phases 3-5)

---

## Summary

**Total**: 48 tasks (47 implementation + 1 epic)
- Epic: 1
- Features (Phases): 10
- Tasks: 47

**By Priority**:
- Critical (P0): 6 tasks (foundational - MUST complete first)
- High (P1): 24 tasks (P1 user stories + critical path)
- Normal (P2): 12 tasks (P2 user stories + important polish)
- Low (P3): 5 tasks (documentation + nice-to-haves)

**By User Story**:
- Foundation: 9 tasks (setup + foundational)
- US1 (Template Selection): 7 tasks
- US2 (Agent Configuration): 5 tasks
- US3 (UUID Generation): 4 tasks
- US6 (Backward Compatibility): 3 tasks
- US4 (Non-Interactive Mode): 4 tasks
- US5 (List Templates): 2 tasks
- US7 (Claude Settings): 3 tasks
- Polish: 10 tasks

**Story Testability**: ✅ All user stories independently testable
- US1: Create project, select template, verify structure
- US2: Create project, select agent, verify config files
- US3: Create project, verify UUID in metadata
- US4: Run with flags, verify no prompts
- US5: Run --list-templates, verify output
- US6: Use defaults, compare to current sl new
- US7: Select Claude Code, verify settings.json

**Parallel Opportunities**:
- After foundational phase: US1, US2, US3 can proceed in parallel (3 teams)
- After P1 stories: US4, US5, US7 can proceed in parallel (3 teams)

**Suggested MVP**: Phases 1-6 (US1, US2, US3, US6) = 28 tasks delivering core template/agent selection

---

## Next Steps

1. ✅ **Review this tasks.md file** - Verify task breakdown, dependencies, acceptance criteria
2. ⏳ **Execute epic creation** - Run epic CLI command with appropriate labels
3. ⏳ **Execute feature creation** - Create 10 phase features under epic
4. ⏳ **Execute task creation** - Create all 47 tasks under appropriate features
5. ⏳ **Set up dependencies** - Run `sl issue link` commands for task dependencies
6. ⏳ **Begin implementation** - Start with Phase 1 (Setup), then Phase 2 (Foundational)

**To create all issues**, run the commands in the "CLI Commands Reference" section at the top of this file.

---

**Status**: READY FOR REVIEW - Please review task breakdown before executing creation commands.
