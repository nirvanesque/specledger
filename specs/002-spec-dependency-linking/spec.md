# Feature Specification: Spec Dependency Linking

**Feature Branch**: `002-spec-dependency-linking`
**Created**: 2025-01-29
**Status**: Draft
**Input**: User description: "Implement golang-style dependency locking and linking for specifications, allowing the current spec to refer to specs from other repositories"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Declare External Spec Dependencies (Priority: P1)

Development teams working on services that share common specifications (e.g., authentication contracts, data models, API schemas) need to declare dependencies on specifications from other repositories. This enables specification reuse across service boundaries while ensuring version compatibility and preventing drift.

**Generated Artifacts**:
- `specs/<NNN>-<feature-name>/spec.mod` - Dependency manifest declaring external spec repositories and versions (similar to `go.mod`)
- `specs/<NNN>-<feature-name>/spec.sum` - Cryptographic hash lockfile verifying dependency integrity (similar to `go.sum`)

**Why this priority**: This is the foundation for cross-repository specification sharing. Without the ability to declare and lock dependencies, teams cannot reliably reference external specs, leading to duplication and inconsistency across services. This delivers immediate value by enabling specification reuse.

**Independent Test**: Can be fully tested by creating a `spec.mod` file declaring a dependency on an external repository's spec, running the dependency resolution command, and verifying that `spec.sum` is generated with correct cryptographic hashes. Delivers value by providing a reliable mechanism to reference external specifications.

**Acceptance Scenarios**:

1. **Given** a team wants to use a specification from another repository, **When** they add a dependency declaration to `spec.mod`, **Then** the system records the repository URL, branch/tag, and specification path
2. **Given** a `spec.mod` file exists with external dependencies, **When** the dependency resolution command is invoked, **Then** the system fetches the external specs and generates `spec.sum` with cryptographic hashes for verification
3. **Given** multiple specifications from the same external repository are needed, **When** dependencies are declared, **Then** the system fetches all referenced specs efficiently (single fetch per repository)
4. **Given** a dependency is declared with a branch reference (e.g., `main`), **When** the dependency is resolved, **Then** the system records the specific commit hash in `spec.sum` for reproducibility

---

### User Story 2 - Reference External Specs in Current Spec (Priority: P2)

Once dependencies are declared, teams need to reference specific sections, entities, or requirements from external specifications within their own specification. This enables composition and extension of shared specifications without duplication.

**Generated Artifacts**:
- Inline references in `spec.md` using markdown link syntax with validation
- Resolved dependency graph showing transitive dependencies

**Why this priority**: Referencing external specs provides the actual usage mechanism. Without the ability to reference specific sections, declaring dependencies has limited utility. This builds on P1 by adding the "usage" layer to the "declaration" layer.

**Independent Test**: Can be fully tested by adding reference links to `spec.md` that point to external specs, running validation, and verifying that references resolve correctly. Delivers value by enabling specification composition.

**Acceptance Scenarios**:

1. **Given** a dependency is declared for a shared data model spec, **When** writing the current spec, **Then** the user can reference external entities using `[External Spec](repo-url#spec-id#section)` syntax
2. **Given** a specification contains external references, **When** the spec is validated, **Then** the system verifies all references exist and are accessible
3. **Given** an external spec is updated, **When** the current spec references it, **Then** the system detects version mismatches and warns the user
4. **Given** a reference points to a non-existent external section, **When** validation runs, **Then** the system reports the specific reference that failed to resolve

---

### User Story 3 - Update and Pin Dependency Versions (Priority: P3)

As external specifications evolve, teams need to update their dependencies to get fixes and improvements while maintaining control over when updates occur. Teams must be able to pin specific versions and update intentionally.

**Generated Artifacts**:
- Updated `spec.mod` with new version constraints
- Updated `spec.sum` with new cryptographic hashes
- Migration notes showing differences between old and new dependency versions

**Why this priority**: Version management is critical for long-term maintenance but is not required for initial dependency usage. Teams can start with fixed versions and adopt update workflows as dependencies evolve.

**Independent Test**: Can be fully tested by running an update command on dependencies, verifying version changes, and checking that `spec.sum` hashes are updated. Delivers value by enabling controlled dependency evolution.

**Acceptance Scenarios**:

1. **Given** a dependency is pinned to a specific version, **When** the user runs the update command, **Then** the system fetches the latest version from the declared branch and updates the pin
2. **Given** an external spec has breaking changes, **When** the user updates to the new version, **Then** the system shows a diff of changes and prompts for confirmation
3. **Given** a dependency is updated, **When** the current spec references changed sections, **Then** the system identifies which references need review
4. **Given** a user wants to lock to an exact commit, **When** they specify a commit hash in `spec.mod`, **Then** the system uses that exact version regardless of branch updates

---

### User Story 4 - Detect and Resolve Dependency Conflicts (Priority: P4)

When multiple dependencies reference the same external specification with different version requirements, or when circular dependencies exist, the system must detect and help resolve these conflicts.

**Generated Artifacts**:
- Dependency conflict report showing incompatible version requirements
- Resolution suggestions (update all, use compatible version, create local copy)

**Why this priority**: Conflict detection is important for complex dependency graphs but simple projects may not encounter conflicts initially. This becomes valuable as the number of dependencies grows.

**Independent Test**: Can be fully tested by creating conflicting dependencies and running validation, which should report the conflict with resolution options. Delivers value by preventing subtle bugs from version mismatches.

**Acceptance Scenarios**:

1. **Given** two dependencies require different versions of the same external spec, **When** validation runs, **Then** the system reports the conflict with version details
2. **Given** a circular dependency is detected (A depends on B, B depends on A), **When** validation runs, **Then** the system reports the cycle and suggests breaking it
3. **Given** a conflict is detected, **When** the user chooses to resolve by updating, **Then** the system updates all dependencies to a compatible version
4. **Given** transitive dependencies exist, **When** the dependency graph is displayed, **Then** it shows the full tree including indirect dependencies

---

### User Story 5 - Vendor Dependencies for Offline Use (Priority: P5)

Teams working in restricted environments or needing to ensure build reproducibility may need to vendor (copy) external specifications into their repository for offline access and version control.

**Generated Artifacts**:
- `specs/vendor/` directory containing copies of all external dependencies
- Updated `spec.sum` referencing vendored copies instead of remote URLs

**Why this priority**: Vendoring is important for certain environments but is not required for basic dependency usage. Teams can adopt this when offline access or strict reproducibility is needed.

**Independent Test**: Can be fully tested by running the vendor command, verifying copies exist in the vendor directory, and confirming that specs work offline. Delivers value by enabling air-gapped workflows.

**Acceptance Scenarios**:

1. **Given** external dependencies are declared, **When** the vendor command is invoked, **Then** all external specs are copied into `specs/vendor/` with preserved directory structure
2. **Given** dependencies are vendored, **When** working offline, **Then** all spec references resolve from the local vendor directory
3. **Given** vendored dependencies exist, **When** a new dependency is added, **Then** the vendor command can be run again to include the new dependency
4. **Given** a vendored dependency is modified locally, **When** validation runs, **Then** the system warns that the vendored copy differs from the remote version

---

### Edge Cases

- **External repository unavailable**: System uses cached version with warning if available; otherwise returns error blocking operation (FR-033)
- **Private repository authentication**: System authenticates using SSH keys or tokens from environment variables or config files (FR-026, FR-028)
- **Referenced external spec section deleted or renamed**: Validation reports the specific reference that failed to resolve (US2 Acceptance Scenario 4)
- **Very large dependency graphs (20+ external specs)**: Dependency graph visualization displays for up to 50 transitive dependencies (SC-010); resolution completes within 30 seconds for 10 repos (SC-003)
- **Hash mismatch in spec.sum**: System detects mismatch and prompts user to regenerate lockfile or abort (FR-027)
- **Merge conflicts when multiple developers update dependencies**: System tracks changes via Git; conflicts resolved using standard Git conflict resolution (inherited from 001 branching/versioning)
- **Dependency references non-existent commit hash**: System validates commit existence during resolution; returns error if invalid
- **Cross-references between vendored and non-vendored specs**: System warns when vendored copy differs from remote version (FR-025)

## Requirements *(mandatory)*

### Functional Requirements

**Dependency Declaration**:
- **FR-001**: System MUST support declaring external spec dependencies in a `spec.mod` file with repository URL, branch/tag, and spec path
- **FR-002**: System MUST support version pinning using commit hashes, branch names, or semantic version tags
- **FR-003**: System MUST validate `spec.mod` syntax and report clear errors for malformed declarations
- **FR-004**: System MUST support indirect (transitive) dependencies discovered from external specs' own `spec.mod` files

**Dependency Resolution and Locking**:
- **FR-005**: System MUST fetch external specifications from Git repositories accessible via HTTP(S) or SSH
- **FR-006**: System MUST generate a `spec.sum` file containing cryptographic hashes (SHA-256) of all resolved dependencies
- **FR-007**: System MUST verify fetched content against `spec.sum` hashes before use
- **FR-008**: System MUST record the exact commit hash for each resolved dependency
- **FR-009**: System MUST support fetching from the `main` or `master` branch by default when no branch is specified
- **FR-033**: When an external repository is unavailable, system MUST use the last successfully cached version if available, otherwise return an error

**Specification References**:
- **FR-010**: System MUST support referencing external spec sections using `[Name](repo-url#spec-id#section-id)` syntax
- **FR-011**: System MUST validate that all external references resolve to existing sections
- **FR-012**: System MUST support referencing specific entities, requirements, or user stories from external specs
- **FR-013**: System MUST provide a command to list all external references in the current spec

**Version Management**:
- **FR-014**: System MUST support updating dependencies to their latest compatible versions
- **FR-015**: System MUST show diffs between current and updated dependency versions
- **FR-016**: System MUST warn when an update would break existing references
- **FR-017**: System MUST support semantic version constraints (e.g., `^1.2.0`, `~2.0.0`)

**Conflict Detection**:
- **FR-018**: System MUST detect when multiple dependencies require different versions of the same external spec
- **FR-019**: System MUST detect circular dependencies in the dependency graph
- **FR-020**: System MUST provide resolution suggestions for detected conflicts
- **FR-021**: System MUST prevent dependency resolution when conflicts exist without explicit user override

**Vendoring**:
- **FR-022**: System MUST support copying external dependencies into a `specs/vendor/` directory
- **FR-023**: System MUST preserve directory structure and metadata when vendoring
- **FR-024**: System MUST support using vendored dependencies when remote repositories are unavailable
- **FR-025**: System MUST warn when vendored copies differ from remote versions

**Security and Integrity**:
- **FR-026**: System MUST authenticate private repositories using SSH keys or tokens
- **FR-027**: System MUST detect when `spec.sum` hash doesn't match fetched content and prompt user to regenerate the lockfile or abort
- **FR-028**: System MUST support authentication tokens via environment variables or config files
- **FR-029**: System MUST never cache credentials in plaintext

**Integration**:
- **FR-030**: System MUST integrate with existing SpecLedger commands (`/specledger.specify`, `/specledger.plan`, etc.)
- **FR-031**: System MUST include dependency information in generated documentation
- **FR-032**: System MUST support querying the dependency graph via CLI commands

### Key Entities

**spec.mod** (Dependency Manifest):
- Represents the declaration of external specification dependencies
- Contains repository URLs, version constraints, spec paths, and a unique spec ID for external references
- Format: Text file with `require <repo-url> <version> <spec-path>` syntax and optional `id <spec-id>` declaration
- Example: `require github.com/example/common-specs v1.2.0 specs/auth.md` with `id common-auth`

**spec.sum** (Lockfile):
- Represents the locked version of each dependency with cryptographic verification
- Contains repository URLs, commit hashes, SHA-256 content hashes, and spec paths
- Format: Text file with `<repo-url> <commit-hash> <sha256-hash> <spec-path>` entries
- Example: `github.com/example/common-specs abc123def456 sha256:... specs/auth.md`

**External Reference**:
- Represents a link from the current spec to a section in an external spec
- Attributes: source location, target spec identifier, target section ID, display text
- Target spec identifier is read from the external repository's `spec.mod` metadata file
- Validated against resolved dependencies to ensure targets exist

**Dependency Graph**:
- Represents all direct and transitive dependencies for a specification
- Nodes: External specifications
- Edges: "depends on" relationships
- Attributes: Version constraints, resolved versions, conflicts

**Vendored Spec**:
- Represents a local copy of an external specification stored in `specs/vendor/`
- Attributes: original repository URL, commit hash, local path, content hash
- Used for offline operation and reproducibility

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can declare and resolve a single external spec dependency in under 10 seconds
- **SC-002**: System validates all external references in a specification in under 5 seconds
- **SC-003**: Dependency resolution completes in under 30 seconds for 10 external repositories
- **SC-004**: 100% of dependency conflicts are detected before spec is used for planning or implementation
- **SC-005**: Users can update all dependencies and review changes in under 2 minutes
- **SC-006**: Vendoring completes in under 60 seconds for 20 external specifications
- **SC-007**: Cryptographic hash verification detects 100% of content modifications
- **SC-008**: Users can successfully add their first external spec reference on first attempt 90% of the time
- **SC-009**: System handles private repository authentication with cached credentials 100% of the time
- **SC-010**: Dependency graph visualization displays for specs with up to 50 transitive dependencies

### Previous work

**Epic: 001-sdd-control-plane - SpecLedger SDD Control Plane**

- **Specification Management (001)**: Establishes the foundation for specification creation, versioning, and collaboration. Dependency linking builds on this by extending specifications to reference external content.
- **Branching and Versioning (001)**: The existing Git-based versioning for specifications will be leveraged to track dependency versions via commit hashes.

**Notes**: This feature extends the core SpecLedger platform to support multi-repository specification sharing, enabling teams to maintain separate codebases while sharing common specification artifacts.

## Clarifications

### Session 2025-01-29

- Q: For the external spec reference syntax `[Name](repo-url#spec-id#section-id)`, how should the `spec-id` component be determined for external repositories? → A: Metadata file - External repositories must include a `spec.mod` or metadata file that declares a unique spec ID. This enables flexible identification without relying on directory structure conventions.
- Q: When an external repository becomes unavailable (edge case #1), what should the system behavior be? → A: Return error but use local cache if available first - The system should attempt to use the last successfully fetched cached version with a warning that the remote is unavailable. If no cache exists, return an error and block the operation.
- Q: When `spec.sum` hash doesn't match fetched content (tampering detection, edge case #5), what action should the system take? → A: Ask user to regenerate the sum - The system should detect the mismatch, warn the user, and prompt them to regenerate the `spec.sum` file with the new hash (or abort if they don't trust the change).
