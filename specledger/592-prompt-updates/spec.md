# Feature Specification: Improve SpecLedger Command Prompts

**Feature Branch**: `592-prompt-updates`
**Created**: 2026-02-20
**Status**: Draft
**Input**: User description: "update specledger prompts (both in .claude and embedded) - specledger.specify: utilize dependency (sl deps) if referred by user, specledger.tasks: fix sl issue link/create errors, utilize definition of done, make issues more descriptive, specledger.implement: check definition of done and acceptance criteria"

## Clarifications

### Session 2026-02-20

- Q: How should the /specledger.specify command detect dependency references in user descriptions? → A: Explicit syntax only (e.g., `deps:alias-name` or `@dependency`)
- Q: How should the /specledger.implement command verify Definition of Done items before closing issues? → A: Automated where possible (verify programmatically), fall back to interactive for others
- Q: What level of error detail should be provided when sl issue create/link commands fail? → A: Automatically fix errors when possible, retry with corrected parameters
- Q: Should /specledger.tasks fill definition_of_done for generated issues? → A: Yes, populate DoD items from acceptance criteria in spec.md
- Q: Should /specledger.tasks also display/reference DoD items in tasks.md? → A: Yes, include DoD summary in tasks.md for visibility and reference during implementation

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Reference Dependencies During Specification (Priority: P1)

As a developer creating a feature specification, when I mention external specifications, APIs, or dependencies in my feature description, the `/specledger.specify` command should automatically recognize these references and either load context from `sl deps` or prompt me to add them as dependencies.

**Why this priority**: This is the starting point of the workflow. If dependencies aren't properly recognized during specification, downstream tasks will lack necessary context.

**Independent Test**: Can be fully tested by creating a spec that references an external API (e.g., "integrate with Stripe payment API") and verifying that the system either loads existing Stripe deps context or prompts to add them.

**Acceptance Scenarios**:

1. **Given** a user creates a spec with explicit syntax `deps:api-contracts`, **When** the dependency exists via `sl deps`, **Then** the spec generation should load and reference relevant context from that dependency
2. **Given** a user references a dependency using `@dependency-alias` syntax, **When** the dependency exists, **Then** load the content and include relevant context in the spec
3. **Given** a user explicitly references a dependency that doesn't exist, **When** no matching dep is found, **Then** the system should display: "Dependency 'X' not found. Use 'sl deps add --alias X <source>' to add it."

---

### User Story 2 - Generate Descriptive, Complete Issues (Priority: P1)

As a developer running `/specledger.tasks`, I want the generated issues to be descriptive, complete, and concise so that any developer can pick up a task and understand what needs to be done without additional context gathering.

**Why this priority**: Task quality directly impacts implementation efficiency. Poorly described tasks lead to context-switching and rework.

**Independent Test**: Can be tested by generating tasks from a completed plan and verifying that each issue: (a) has a clear problem statement, (b) describes inputs/outputs, (c) references relevant files, (d) has testable acceptance criteria.

**Acceptance Scenarios**:

1. **Given** a plan.md with clear components, **When** tasks are generated, **Then** each issue includes a concise title (under 80 chars), a problem statement explaining WHY, implementation details explaining HOW/WHERE, and acceptance criteria for WHAT success looks like
2. **Given** tasks are being created from spec.md with acceptance criteria, **When** the `sl issue create` command is called, **Then** each issue should have definition_of_done items populated from the relevant acceptance criteria
3. **Given** tasks.md is generated, **When** viewing the file, **Then** it should include a DoD summary section listing the DoD items for each issue (referenceable by issue ID)
4. **Given** tasks are being created, **When** the `sl issue create` command is called, **Then** it should succeed without errors (automatically fix special characters, retry on transient failures)
5. **Given** tasks need dependencies linked, **When** the `sl issue link` command is called, **Then** it should properly establish the relationship without errors

---

### User Story 3 - Utilize Definition of Done During Implementation (Priority: P1)

As a developer implementing tasks with `/specledger.implement`, I want the system to check each task's Definition of Done (DoD) and acceptance criteria before marking it complete, ensuring quality standards are met.

**Why this priority**: Without proper DoD verification, tasks may be marked complete prematurely, leading to technical debt and incomplete features.

**Independent Test**: Can be tested by implementing a task that has DoD items, then verifying the system checks each DoD item before allowing closure.

**Acceptance Scenarios**:

1. **Given** a task with DoD items, **When** implementation is complete, **Then** the system should attempt automated verification (e.g., file exists, tests pass, syntax valid) for applicable items
2. **Given** DoD items that cannot be verified automatically, **When** implementation is complete, **Then** the system should prompt the user for interactive confirmation of remaining items
3. **Given** a task where automated DoD items fail verification, **When** attempting to close the issue, **Then** the system should display which items failed and why, requiring explicit confirmation to proceed
4. **Given** all DoD items pass (automated or confirmed), **When** closing the issue, **Then** proceed without additional prompts

---

### Edge Cases

- **Missing dependency cache**: When `sl deps` references a dependency that no longer exists in the cache, display error: "Dependency 'X' not found in cache. Run 'sl deps add --alias X <source>' to restore it."
- **File system errors**: When `sl issue create` fails due to file system errors (e.g., permissions), display clear error message with the specific path and permission issue, then suggest remediation steps
- **Ambiguous dependency reference**: When using explicit syntax, ambiguity is eliminated - exact alias match required
- **Special characters in DoD items**: DoD items should be sanitized when parsed; special characters are escaped or quoted in the issue content

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The `/specledger.specify` command MUST detect dependency references using explicit syntax (`deps:alias` or `@alias`)
- **FR-002**: The `/specledger.specify` command MUST load existing dependency content from `sl deps list` cache when referenced dependencies exist
- **FR-003**: The `/specledger.specify` command MUST prompt users to add missing dependencies when explicit references don't resolve
- **FR-004**: The `/specledger.tasks` command MUST generate issues with structured content: title, problem statement (WHY), implementation details (HOW/WHERE), and acceptance criteria (WHAT)
- **FR-005**: The `/specledger.tasks` command MUST populate definition_of_done items in each generated issue from acceptance criteria in spec.md
- **FR-006**: The `/specledger.tasks` command MUST include a DoD summary section in tasks.md referencing the DoD items for each generated issue
- **FR-007**: The `/specledger.tasks` command MUST handle `sl issue create` and `sl issue link` errors by automatically retrying with corrected parameters
- **FR-008**: The `/specledger.implement` command MUST attempt automated verification of DoD items where possible (file existence, test results, syntax validation)
- **FR-009**: The `/specledger.implement` command MUST fall back to interactive confirmation for DoD items that cannot be verified automatically
- **FR-010**: The `/specledger.implement` command MUST display failed verification results with clear explanations before requiring explicit confirmation
- **FR-011**: Both `.claude/commands/` and `pkg/embedded/skills/commands/` prompt files MUST be updated with the same changes

### Key Entities

- **Dependency Reference**: A mention of an external specification, API, or system within a feature description that should be resolved via `sl deps`
- **Issue Content Structure**: A structured format containing: title, problem statement, implementation details, acceptance criteria, and definition_of_done items
- **Definition of Done (DoD)**: A checklist of verifiable items that must be satisfied before a task is considered complete

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can reference dependencies in spec descriptions and have relevant context automatically loaded 95% of the time when dependencies exist
- **SC-002**: Generated issues have all required fields (title, problem statement, implementation details, acceptance criteria) 100% of the time
- **SC-003**: `sl issue create` and `sl issue link` commands succeed on first attempt 99% of the time (error handling for edge cases)
- **SC-004**: Tasks are never closed without DoD verification unless explicitly forced
- **SC-005**: Developers can understand what a task requires without additional context gathering in 90% of cases

### Previous work

- **591-issue-tracking-upgrade**: Built-in issue tracking system with `sl issue` commands
- **008-fix-sl-deps**: SpecLedger dependencies management system

### Epic: 592 - Improve SpecLedger Command Prompts

- **Dependency-aware specification**: Recognize and load external dependencies during spec creation
- **Quality task generation**: Generate descriptive, complete, concise issues with DoD items, include DoD summary in tasks.md
- **DoD-enforced implementation**: Verify definition of done and acceptance criteria during implementation
