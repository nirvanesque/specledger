# Feature Specification: Built-In Issue Tracker

**Feature Branch**: `591-issue-tracking-upgrade`
**Created**: 2026-02-18
**Status**: Draft
**Input**: User description: "Beads is slow, need to kill daemon. We don't need beads issue tracking, can be replaced with the sl binary. Usage of beads couple specledger with beads. Solution: Removed the beads + perles usage. Create simple issue tracker within sl CLI. Opensource CLI -> json file artifacts. Use backward comp beads artifact. Don't keep issue list at repo root, keep at spec level so there's less cross branch conflicts."

## Clarifications

### Session 2026-02-18

- Q: How should merge conflicts on issues.jsonl be resolved when merging branches? → A: Auto-merge with dedup (automatically deduplicate by issue ID, keeping both versions' changes where possible)
- Q: What priority scale should issues use? → A: Numeric 0-5 (0 = highest, 5 = lowest; matches Beads format)
- Q: What should happen when issue commands run outside a feature branch context? → A: Fail with error (clear message: "Not on a feature branch. Use --spec flag or checkout a ###-branch.")

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create and Manage Issues (Priority: P1)

As a developer using SpecLedger, I want to create, list, and manage issues directly from the sl CLI so that I can track work without external dependencies or slow daemon processes.

**Why this priority**: Issue creation and management is the core functionality. Without this, no other features matter. This replaces the primary use case of Beads.

**Independent Test**: Can be fully tested by creating an issue with `sl issue create`, listing it with `sl issue list`, updating it with `sl issue update`, and closing it with `sl issue close`. Delivers immediate value as a standalone issue tracker.

**Acceptance Scenarios**:

1. **Given** a spec directory exists at `specledger/010-my-feature/`, **When** I run `sl issue create --title "Add validation" --type task`, **Then** a new issue is created with a unique ID and stored in `specledger/010-my-feature/.sl/issues.jsonl`
2. **Given** issues exist in a spec directory, **When** I run `sl issue list`, **Then** all issues for the current spec are displayed in a readable table format
3. **Given** an issue exists with status "open", **When** I run `sl issue close ISSUE-ID`, **Then** the issue status changes to "closed" with a closed_at timestamp

---

### User Story 2 - Migrate Existing Beads Data (Priority: P2)

As a developer with existing Beads issues, I want to migrate my `.beads/issues.jsonl` data to the new format so that I don't lose historical work tracking.

**Why this priority**: Existing users need continuity. Without migration, users would lose all their issue history when upgrading.

**Independent Test**: Can be tested by running `sl issue migrate` on a repository with existing Beads data and verifying the new format files contain all the original data.

**Acceptance Scenarios**:

1. **Given** a `.beads/issues.jsonl` file exists at repository root, **When** I run `sl issue migrate`, **Then** issues are distributed to their respective spec directories based on branch/feature association
2. **Given** migration completes successfully, **When** I compare original Beads data with migrated data, **Then** all issue fields (id, title, description, status, priority, type, timestamps, dependencies) are preserved
3. **Given** an issue cannot be mapped to a spec directory, **When** migration runs, **Then** the issue is placed in a fallback location with a warning message

---

### User Story 3 - Track Dependencies Between Issues (Priority: P3)

As a developer planning complex features, I want to define dependencies between issues so that I can understand the order of work and identify blocking relationships.

**Why this priority**: Dependency tracking adds significant value for complex projects but is not required for basic issue management functionality.

**Independent Test**: Can be tested by creating two issues, linking them with `sl issue link ISSUE-A blocks ISSUE-B`, and verifying the dependency is reflected in both `sl issue list --tree` and JSON exports.

**Acceptance Scenarios**:

1. **Given** two issues exist in the same spec, **When** I run `sl issue link ISSUE-A blocks ISSUE-B`, **Then** ISSUE-B shows ISSUE-A in its `blocked_by` list and ISSUE-A shows ISSUE-B in its `blocks` list
2. **Given** issues with dependencies exist, **When** I run `sl issue list --tree`, **Then** issues are displayed in a tree structure showing parent-child and blocking relationships
3. **Given** a dependency would create a cycle, **When** I attempt to create it, **Then** the operation fails with a clear error message explaining the cycle

---

### User Story 4 - Work Across Multiple Specs (Priority: P3)

As a developer working on multiple features, I want to list and filter issues across all specs so that I can see my complete workload in one view.

**Why this priority**: Cross-spec visibility is valuable for project management but individual spec tracking is the primary use case.

**Independent Test**: Can be tested by creating issues in multiple spec directories and running `sl issue list --all` to see them aggregated.

**Acceptance Scenarios**:

1. **Given** issues exist in multiple spec directories, **When** I run `sl issue list --all`, **Then** all issues from all specs are listed with their spec context
2. **Given** issues across specs, **When** I run `sl issue list --all --status open`, **Then** only open issues from all specs are shown
3. **Given** issues across specs, **When** I run `sl issue list --all --type epic`, **Then** only epic-type issues are shown with their spec prefixes

---

### Edge Cases

- What happens when the spec directory doesn't exist when creating an issue? The command should fail with a helpful error message suggesting to run the spec initialization first.
- How does the system handle concurrent writes to issues.jsonl? File locking prevents corruption; operations queue if lock is held.
- What happens if issues.jsonl becomes corrupted? A `sl issue repair` command attempts to recover valid JSON lines and reports any unrecoverable data.
- How are issue IDs generated to avoid conflicts? IDs use format `SL-###` where ### is auto-incremented per spec directory, ensuring uniqueness within each spec context.
- How are merge conflicts on issues.jsonl resolved? Auto-merge with deduplication by issue ID, preserving both branches' changes where possible.
- What happens if commands run outside a feature branch? Command fails with error: "Not on a feature branch. Use --spec flag or checkout a ###-branch."

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST store issues in JSONL format at `specledger/<spec-dir>/.sl/issues.jsonl` (per-spec storage)
- **FR-002**: System MUST support creating issues with fields: id, title, description, status, priority, issue_type, created_at, updated_at
- **FR-003**: System MUST generate unique issue IDs within each spec directory using format `SL-###`
- **FR-004**: System MUST support issue types: epic, feature, task, bug
- **FR-005**: System MUST support issue statuses: open, in_progress, closed
- **FR-006**: System MUST support issue priority as numeric 0-5 (0 = highest, 5 = lowest)
- **FR-007**: System MUST provide `sl issue create` command with flags for --title, --description, --type, --priority
- **FR-008**: System MUST provide `sl issue list` command with optional --status, --type, --all flags
- **FR-009**: System MUST provide `sl issue update` command to modify existing issue fields
- **FR-010**: System MUST provide `sl issue close` command that sets status to closed and records closed_at timestamp
- **FR-011**: System MUST provide `sl issue show <id>` command to display full issue details
- **FR-012**: System MUST maintain backward compatibility with Beads JSONL format for migration purposes
- **FR-013**: System MUST provide `sl issue migrate` command to convert `.beads/issues.jsonl` to per-spec format
- **FR-014**: System MUST support issue dependencies with fields: blocked_by, blocks
- **FR-015**: System MUST provide `sl issue link <id1> <relationship> <id2>` for dependency management
- **FR-016**: System MUST NOT require any daemon process or background service
- **FR-017**: System MUST complete all operations with file I/O only (no database required)
- **FR-018**: System MUST detect current spec context from git branch name (###-short-name pattern)
- **FR-019**: System MUST fail with error when no spec context detected and no --spec flag provided
- **FR-020**: System MUST auto-merge issues.jsonl conflicts by deduplicating on issue ID, preserving both branches' changes where possible

### Key Entities

- **Issue**: Core tracking unit with unique ID, title, description, type (epic/feature/task/bug), status (open/in_progress/closed), priority (0-5, where 0=highest), timestamps, and optional dependencies
- **IssueStore**: JSONL file at `specledger/<spec>/.sl/issues.jsonl` containing all issues for that spec
- **Dependency**: Relationship between issues with type (blocks/blocked_by/parent-child) and direction

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Issue creation and listing commands complete in under 100ms for files with up to 1000 issues (no daemon overhead)
- **SC-002**: Migration of existing Beads data preserves 100% of issue data without manual intervention
- **SC-003**: Zero daemon processes required - all operations are direct file I/O
- **SC-004**: Issue storage at spec level eliminates cross-branch merge conflicts for typical workflows
- **SC-005**: Users can perform all common issue operations (create, list, update, close) without consulting documentation (intuitive CLI design)
- **SC-006**: New users can start tracking issues within 30 seconds of installing sl CLI

### Previous work

- **010-checkpoint-session-capture**: Established pattern for session capture via hooks; this feature will integrate with session tracking for task completion events
- **009-command-system-enhancements**: Established CLI command patterns in Go/Cobra that this feature will follow

## Dependencies & Assumptions

### Dependencies

- Go 1.24+ and Cobra CLI framework (existing)
- Git for branch detection (existing)

### Assumptions

- Users work primarily within a single spec/feature at a time, making per-spec storage natural
- Issue volume per spec is typically under 1000 issues, making JSONL performant
- Users accept migrating from Beads once to gain simplified architecture
- Standard JSONL line-append pattern is sufficient for concurrent write safety (no high-frequency concurrent writes expected)
- Issue IDs can be spec-scoped rather than globally unique across the repository

### Out of Scope

- Real-time collaboration features (requires backend service)
- Issue synchronization across multiple machines (can be added later via git-based sync)
- Advanced querying beyond basic filters (can be added later)
- Integration with external issue trackers (Jira, Linear, GitHub Issues)
- Web UI for issue management
