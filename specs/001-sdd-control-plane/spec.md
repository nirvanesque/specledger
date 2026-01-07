# Feature Specification: SpecLedger - SDD Control Plane

**Feature Branch**: `001-sdd-control-plane`
**Created**: 2025-12-22
**Status**: Draft
**Input**: User description: "Develop SpecLedger, an LLM-driven Specification Driven Development workflow platform for cross team collaboration through a control plane. The platform extends SDD beyond a single developer workflow inta a shared, auditable, and scalable collaboration model for humans and LLM Agent collaboration. At its core, the platform is driven through a remote control plane that captures and links: 1. specifications (User stories, Functional Requirements, Edge Case clarification questions and answers), 2. Implementation planning, technical research (alternatives, tradeoffs, tech stack decisions) and quickstart examples for user stories, 3. Generated task graphs organised by specification, broken down across phases with cross phase and task dependency and priority tracking. 4. Per task implementation session history logs from LLM and human interactions providing file edit information and user decision points, course adjustments or clarifications. Each workflow step (1-4) is executed through LLM-assisted commands, but every decision, clarification and alternative explored is checkpointed and versioned on a central platform. This enables branching, comparison of approaches and safe rollback while preserving the full reasoning trail behind every outcome."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Capture and Version Specifications (Priority: P1)

Development teams need to capture user requirements, functional specifications, and clarifications in a structured, version-controlled format. Team members and LLM agents collaborate using `/specledger.specify` to create specifications and `/specledger.clarify` to refine them through iterative Q&A sessions, ensuring all stakeholders have a shared understanding before implementation begins.

**Generated Artifacts**:
- `specs/<NNN>-<feature-name>/spec.md` - Feature specification with user stories, functional requirements, success criteria
- `specs/<NNN>-<feature-name>/checklists/requirements.md` - Quality validation checklist
- `## Clarifications` section with `### Session YYYY-MM-DD` subheadings tracking Q&A history

**Why this priority**: This is the foundation of the entire SDD workflow. Without structured specifications, all downstream activities (planning, task generation, implementation) lack a reliable source of truth. This delivers immediate value by centralizing requirement gathering and eliminating ambiguity.

**Independent Test**: Can be fully tested by invoking `/specledger.specify` to create a specification, then `/specledger.clarify` to add clarification questions and answers, and verifying the changes are versioned in the control plane. Delivers value by providing a single source of truth for requirements.

**Acceptance Scenarios**:

1. **Given** a team wants to start a new feature, **When** they invoke `/specledger.specify` with a feature description, **Then** the system creates a numbered feature branch (e.g., `001-feature-name`), generates `spec.md` with user stories and functional requirements, and stores it in the control plane
2. **Given** an existing specification has ambiguous requirements, **When** team members invoke `/specledger.clarify`, **Then** the system scans for ambiguities across taxonomy categories (Functional Scope, Domain Model, Non-Functional, etc.), asks up to 10 targeted questions, and records answers in a `## Clarifications` section with session timestamps
3. **Given** a specification has been modified multiple times, **When** a user requests the version history, **Then** they see all changes with timestamps, authors (human or LLM agent), and can view or restore any previous version
4. **Given** multiple team members are collaborating on a specification, **When** they make concurrent edits, **Then** the system tracks all changes using Git-style merge with conflict markers

**Follow up Questions**:

- [vincent]: Do specs live in-repo on filesystem per feature-name directories still or do we leverage the control plane to keep track and link one feature to many branches?
- [vincent]: What is the checkpoint system, does it just capture document versions out of git? (consider how a CMS has a version history per document and how it relates to git versioning? Do we snapshot a document on every LLM file-edit decoupled from git commits completely? this seems to be what Kiro does? See also the Claude's `/home/vincent/.claude/file-history` )

---

### User Story 2 - Track Implementation Planning and Research (Priority: P2)

Once specifications are defined, technical leads and LLM agents invoke `/specledger.plan` to explore implementation approaches, document technical alternatives, evaluate tradeoffs, and capture decisions about technology choices and architectural patterns. The planning workflow executes in phases, generating design artifacts linked to the originating specification.

**Generated Artifacts**:
- `specs/<NNN>-<feature-name>/plan.md` - Implementation plan with technical context, constitution checks, and phase gates
- `specs/<NNN>-<feature-name>/research.md` - Research findings with Decision/Rationale/Alternatives format
- `specs/<NNN>-<feature-name>/data-model.md` - Entity definitions, fields, relationships, validation rules
- `specs/<NNN>-<feature-name>/contracts/` - API contracts (Proto/OpenAPI schemas)
- `specs/<NNN>-<feature-name>/quickstart.md` - Test scenarios and getting-started examples

**Why this priority**: Planning bridges the gap between requirements and execution. Without documented alternatives and tradeoffs, teams lose the reasoning behind technical decisions, making future maintenance and adaptation difficult. This builds on P1 by adding the "how" layer to the "what" layer.

**Independent Test**: Can be fully tested by invoking `/specledger.plan` on a completed specification, verifying that research.md resolves all unknowns, data-model.md captures entities, and contracts/ defines API schemas. Delivers value by preserving technical decision-making rationale.

**Acceptance Scenarios**:

1. **Given** a completed specification, **When** a technical lead invokes `/specledger.plan`, **Then** the system executes Phase 0 (research) to resolve all "NEEDS CLARIFICATION" items and generates `research.md` with Decision/Rationale/Alternatives for each unknown
2. **Given** Phase 0 research is complete, **When** the system executes Phase 1 (design), **Then** it generates `data-model.md` (entities from spec), `contracts/` (API schemas from functional requirements), and `quickstart.md` (test scenarios)
3. **Given** a project constitution exists, **When** the plan is generated, **Then** the system validates against constitution gates and errors if violations are unjustified
4. **Given** planning includes external research, **When** LLM agents search for best practices, **Then** findings are consolidated in `research.md` with references to documentation and patterns

**Follow up Questions**:

- [vincent]: What does it mean to "invoke" the commands if there's a remote component? Currently speckit is just prepared prompts for the user preferred CLI and LLM Model, the output is files on disks or issues created in the issue tracker... This means that SpecLedger provides a "bootstrap" system to set up the filesystem for the preferred Agent Shell - see [specify-cli](https://github.com/github/spec-kit/blob/v0.0.90/src/specify_cli/__init__.py#L637-L749). SpecLedger pulls templates from the control plane, customer can easily fork and customise these templates.. the system must provide a way for customers to keep their own customizations in-line with upstream prompt improvements... - prompts must evolve to be LLM / Agent Shell specific (see Conductor prompts vs SpecKit prompts)... - this is a separate future feature of prompt management, first version can just pull default set of prompts for Claude Code only.
- [vincent]: SpecLedger should come with a CLI similar to [steveyegge/beads](https://github.com/steveyegge/beads) - perhaps backed by SQLite with background sync to the platform? Including "SKILLS" or "POWERS" tuned for the target Agent shell to allow them to use these custom CLI to the best of their capabilities (Does this include Agent Shell hooks such as deterministic issue platform commands? Does this include pulling down prompt templates and create issue tracking tasks?)

---

### User Story 3 - Generate and Manage Task Dependency Graphs (Priority: P3)

Based on specifications and plans, teams invoke `/specledger.tasks` to break down work into granular tasks organized by user story phases. The system uses the Beads issue tracker (`bd` CLI) to create epics, features, and tasks with dependency relationships, labels for traceability, and priority ordering.

**Generated Artifacts**:
- `specs/<NNN>-<feature-name>/tasks.md` - Task index with Beads queries, MVP scope, and phase structure
- Beads issues: Epic (top-level), Features (per phase), Tasks (granular work units)
- Labels: `spec:<slug>`, `phase:<name>`, `story:US1`, `component:<area>`, `requirement:FR-001`

**Issue Hierarchy**:
- **Epic**: Top-level feature container (e.g., `sl-0001`)
- **Feature**: Phase grouping (Setup, Foundational, US1, US2, Polish)
- **Task**: Individual work unit with description, design notes, acceptance criteria

**Why this priority**: Task graphs translate plans into executable work units. This enables parallel work streams and helps teams understand critical paths. While important, teams can manually create task lists initially, making this lower priority than core specification and planning capabilities.

**Independent Test**: Can be fully tested by invoking `/specledger.tasks` on a completed plan, verifying Beads issues are created with correct parent-child relationships, and using `bd ready` to find unblocked tasks. Delivers value by providing a clear execution roadmap.

**Acceptance Scenarios**:

1. **Given** a completed implementation plan, **When** a team invokes `/specledger.tasks`, **Then** the system creates a Beads epic with features per phase (Setup, Foundational, US1, US2, etc.) and tasks with `--deps parent-child:<id>` relationships
2. **Given** tasks have dependencies, **When** a task is marked complete via `bd close`, **Then** `bd ready` shows dependent tasks as unblocked and available for work
3. **Given** tasks span multiple user stories, **When** querying with `bd list --label "story:US1"`, **Then** only tasks for that story are returned with their dependency context preserved
4. **Given** a complex feature with many tasks, **When** filtering with `bd list --label "phase:setup" --limit 10`, **Then** relevant tasks are returned with priority ordering and file paths for implementation

**Follow up Questions**:

- [vincent]: Beads is just an example tracker customized for LLM Agent task management, several UI exist to visualize the interface - recommended to use [zjrosen/perles](https://github.com/zjrosen/perles) viz. does the CLI keep SQLite and background sync to the control plane or simpler just directly interact with the control plan? What about integration with Org preferred issue tracker (Redmine / JIRA / ...)


---

### User Story 4 - Capture Implementation Session History (Priority: P4)

During implementation via `/specledger.implement`, every interaction between developers, LLM agents, and the codebase is captured in JSONL session files. These files record tool calls, file edits, user decisions, and course corrections, creating an audit trail showing how decisions evolved during execution.

**Session Data Sources**:
- Claude Code session files: `~/.claude/projects/<project-hash>/<session-uuid>.jsonl`
- Each JSONL entry contains: timestamp, message type, tool calls, file operations, user responses
- Sessions are linked to tasks via Beads comments (`bd comments add`)

**Captured Events**:
- Tool invocations (Read, Write, Edit, Bash, etc.)
- File modifications with before/after content
- User clarification questions and answers (via AskUserQuestion)
- Course corrections when implementation deviates from plan

**Why this priority**: Session history provides accountability and learning opportunities, but the core SDD workflow can function without detailed implementation logs initially. Teams can adopt this as they mature their processes.

**Independent Test**: Can be fully tested by executing a task via `/specledger.implement`, making file edits, and verifying the session JSONL captures all tool calls with timestamps. Delivers value by enabling post-implementation review and knowledge sharing.

**Acceptance Scenarios**:

1. **Given** a developer starts working on a task, **When** they invoke `/specledger.implement` and make file edits, **Then** each tool call is logged to the session JSONL with file path, operation type, and timestamp
2. **Given** an implementation session encounters uncertainty, **When** the LLM uses AskUserQuestion tool, **Then** the question and user's answer are captured in the session JSONL
3. **Given** implementation deviates from the original plan, **When** course corrections are made, **Then** the session includes the deviation context and can be linked to the Beads task via `bd comments add`
4. **Given** a completed task, **When** querying session history, **Then** the complete JSONL timeline of tool calls, decisions, and interactions is available for audit or replay

---

### User Story 5 - Branch and Compare Approaches (Priority: P5)

Teams need to explore multiple implementation approaches in parallel by creating specification or planning branches, evolving them independently, and comparing outcomes before committing to one approach. This enables safe experimentation without losing work.

**Why this priority**: Branching enables advanced workflows but is not essential for basic SDD adoption. Teams should establish core workflows first before adding this complexity.

**Independent Test**: Can be tested by creating a branch from an existing specification, making divergent changes, comparing the branches side-by-side, and either merging or discarding the branch. Delivers value by supporting experimentation and risk reduction.

**Acceptance Scenarios**:

1. **Given** an existing specification, **When** a team wants to explore an alternative approach, **Then** they can create a branch with a descriptive name
2. **Given** multiple branches exist, **When** viewing branches, **Then** teams see branch names, creation timestamps, and divergence points from the main line
3. **Given** two branches with different approaches, **When** comparing them, **Then** the system highlights differences in specifications, plans, and task graphs
4. **Given** a branch has proven successful, **When** merging back to main, **Then** all versioned artifacts (specs, plans, tasks, sessions) are integrated with preserved history

---

### User Story 6 - Rollback to Previous States (Priority: P6)

When teams discover issues or want to revisit earlier decisions, they need to rollback specifications, plans, or task definitions to previous versions while preserving the complete history trail including the rollback action itself.

**Why this priority**: Rollback is a safety net that becomes valuable as complexity grows. Early adopters can work with manual versioning before investing in automated rollback capabilities.

**Independent Test**: Can be tested by making changes to a specification, rolling back to a previous version, verifying the content is restored, and confirming the rollback action is logged in the history. Delivers value by reducing risk of irreversible mistakes.

**Acceptance Scenarios**:

1. **Given** a specification has been modified several times, **When** a team decides to rollback to version 3, **Then** the specification content is restored to version 3 state
2. **Given** a rollback has occurred, **When** viewing version history, **Then** the rollback action appears as a new entry showing what was restored
3. **Given** a task graph has been regenerated, **When** rolling back to a previous task graph version, **Then** all task states and dependencies are restored
4. **Given** a planning document has diverged significantly, **When** performing a rollback, **Then** all linked artifacts (specs, tasks) maintain referential integrity

---

### Edge Cases

- What happens when multiple users edit the same specification section simultaneously?
- How does the system handle very large task graphs (500+ tasks) with complex cross-phase dependencies?
- What happens when a rollback conflicts with work already in progress on dependent tasks?
- How does the system handle LLM agent sessions that timeout or fail mid-execution?
- What happens when a specification branch is deleted but tasks referencing it are still in progress?
- How does the system handle circular dependencies in task graphs?
- What happens when session history grows very large (10,000+ interactions for a single task)?
- How does the system handle offline work that needs to be synchronized later?

## Requirements *(mandatory)*

### Functional Requirements

**Specification Management**:
- **FR-001**: System MUST allow users to create specifications with title, description, user stories, and functional requirements
- **FR-002**: System MUST assign unique identifiers to each specification
- **FR-003**: System MUST track complete version history for every specification change
- **FR-004**: System MUST link clarification questions and answers to specific requirements within a specification
- **FR-005**: System MUST support retrieving any previous version of a specification
- **FR-006**: System MUST track authorship (human user or LLM agent identifier) for all specification changes
- **FR-007**: System MUST timestamp all specification operations, storing timestamps in UTC and displaying them in the user's local timezone

**Planning and Research**:
- **FR-008**: System MUST allow creation of implementation plans linked to specifications
- **FR-009**: System MUST support documenting multiple technical alternatives with pros and cons
- **FR-010**: System MUST record technology stack decisions with justification
- **FR-011**: System MUST link external references and documentation to planning sections
- **FR-012**: System MUST version planning documents independently from specifications
- **FR-013**: System MUST preserve rejected alternatives for future reference

**Task Management**:
- **FR-014**: System MUST generate task graphs from implementation plans
- **FR-015**: System MUST support task dependencies within phases and across phases
- **FR-016**: System MUST track task status (pending, in-progress, blocked, completed)
- **FR-017**: System MUST automatically identify tasks blocked by dependencies
- **FR-018**: System MUST assign priority levels to tasks
- **FR-019**: System MUST support filtering and querying tasks by phase, status, and priority
- **FR-020**: System MUST detect circular dependencies in task graphs and prevent their creation

**Session Tracking**:
- **FR-021**: System MUST capture all file edits made during task implementation
- **FR-022**: System MUST log user decisions and clarifications during implementation sessions
- **FR-023**: System MUST record course corrections with rationale
- **FR-024**: System MUST link session history to specific tasks
- **FR-025**: System MUST timestamp all session events
- **FR-026**: System MUST distinguish between human and LLM agent actions in session logs

**Branching and Versioning**:
- **FR-027**: System MUST support creating branches from any specification or plan version
- **FR-028**: System MUST track branch genealogy (parent/child relationships)
- **FR-029**: System MUST support comparing two branches to highlight differences
- **FR-030**: System MUST allow merging branches with conflict detection
- **FR-031**: System MUST preserve complete history across branch operations

**Rollback Capabilities**:
- **FR-032**: System MUST support rollback of specifications to any previous version
- **FR-033**: System MUST support rollback of plans to any previous version
- **FR-034**: System MUST support rollback of task graphs to any previous version
- **FR-035**: System MUST record rollback operations in version history
- **FR-036**: System MUST maintain referential integrity when rolling back linked artifacts

**Multi-User Collaboration**:
- **FR-037**: System MUST support concurrent access by multiple users and LLM agents
- **FR-038**: System MUST detect conflicting edits to the same artifact
- **FR-039**: System MUST provide conflict resolution using automatic merge with conflict markers (similar to Git), requiring users to manually resolve conflicts when detected
- **FR-040**: System MUST allow users to view changes to shared artifacts via polling (manual refresh); real-time notifications are out of scope for initial release

**Audit and Compliance**:
- **FR-041**: System MUST provide complete audit trail for all operations
- **FR-042**: System MUST support querying history by user, time range, or artifact
- **FR-043**: System MUST preserve data integrity across all operations (no data loss)
- **FR-044**: System MUST enforce access control to artifacts

### Key Entities

**Specification Artifacts** (from `/specledger.specify` and `/specledger.clarify`):
- **Specification (spec.md)**: Feature requirements including user stories (with priorities P1-P6), functional requirements (FR-001, FR-002...), edge cases, and success criteria. Stored in `specs/<NNN>-<feature-name>/spec.md`
- **Clarification**: Question-answer pair recorded in `## Clarifications` section with `### Session YYYY-MM-DD` subheadings. Linked to specific taxonomy categories (Functional Scope, Domain Model, Non-Functional, etc.)
- **Requirements Checklist**: Quality validation stored in `specs/<NNN>-<feature-name>/checklists/requirements.md`

**Planning Artifacts** (from `/specledger.plan`):
- **Plan (plan.md)**: Implementation approach with technical context, constitution checks, and phase gates
- **Research (research.md)**: Technical decisions with Decision/Rationale/Alternatives format for each unknown
- **Data Model (data-model.md)**: Entity definitions, fields, relationships, validation rules extracted from spec
- **Contracts (contracts/)**: API schemas (Proto/OpenAPI) generated from functional requirements
- **Quickstart (quickstart.md)**: Test scenarios and getting-started examples

**Task Management** (from `/specledger.tasks` via Beads):
- **Epic**: Top-level Beads issue (type: epic) representing the entire feature. Labels: `spec:<slug>`, `component:<area>`
- **Feature**: Phase grouping (type: feature) with `--deps parent-child:<epic-id>`. Labels: `phase:setup`, `phase:US1`, etc.
- **Task**: Individual work unit (type: task) with description, design notes, acceptance criteria, and `--deps parent-child:<feature-id>`. Labels: `story:US1`, `requirement:FR-001`, `component:<area>`

**Session Tracking** (from `/specledger.implement`):
- **Session (JSONL)**: Claude Code session file at `~/.claude/projects/<hash>/<uuid>.jsonl` containing tool calls, file operations, and user responses
- **SessionEvent**: Individual JSONL entry with timestamp, message type (user/assistant/tool), and content
- **TaskComment**: Link between session and Beads task via `bd comments add <task-id>`

**Versioning**:
- **Branch**: Git branch (e.g., `001-feature-name`) containing all feature artifacts
- **Version**: Git commit representing a snapshot of all artifacts at a point in time

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Teams can create and version a complete specification (with 5+ user stories and 10+ functional requirements) in under 30 minutes
- **SC-002**: Users can retrieve complete version history for any artifact (specification, plan, or task graph) in under 5 seconds
- **SC-003**: System supports at least 50 concurrent users across 20 active specifications without performance degradation
- **SC-004**: 95% of dependency conflicts in task graphs are automatically detected before task execution begins
- **SC-005**: Teams can compare two specification or planning branches and identify all differences in under 10 seconds
- **SC-006**: Complete audit trail for any artifact is available within 3 seconds of request
- **SC-007**: Rollback operations complete in under 5 seconds and maintain 100% referential integrity
- **SC-008**: Session history captures 100% of file edits and decision points during implementation
- **SC-009**: Users successfully complete specification creation, planning, and task generation workflow on first attempt 80% of the time
- **SC-010**: System prevents 100% of circular dependencies from being created in task graphs

### Previous work

No previous related work found in issue tracker.
