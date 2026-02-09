# Tasks Index: CLI Authentication

Beads Issue Graph Index into the tasks and phases for this feature implementation.
This index does **not contain tasks directly**—those are fully managed through Beads CLI.

**Status**: All tasks completed - Feature fully implemented

## Feature Tracking

* **Beads Epic ID**: `sl-7x3`
* **User Stories Source**: `specs/003-authenticate-cli/spec.md`
* **Research Inputs**: `specs/003-authenticate-cli/research.md`
* **Planning Details**: `specs/003-authenticate-cli/plan.md`
* **Data Model**: `specs/003-authenticate-cli/data-model.md`

## Beads Query Hints

Use the `bd` CLI to query and manipulate the issue graph:

```bash
# Find all tasks for this feature
bd list --label spec:003-authenticate-cli --limit 20

# Find closed tasks
bd list --label spec:003-authenticate-cli --status closed --limit 20

# See dependencies for epic
bd dep tree sl-7x3

# View issues by component
bd list --label 'component:auth' --label 'spec:003-authenticate-cli' --limit 10

# Show all phases
bd list --type feature --label 'spec:003-authenticate-cli'
```

## Tasks and Phases Structure

This feature follows Beads' 2-level graph structure:

* **Epic**: sl-7x3 (CLI Authentication) - CLOSED
* **Phases**: Beads issues of type `feature`, child of the epic
  * sl-0cm: Setup Phase - CLOSED
  * sl-96q: Foundational: Core Auth Infrastructure - CLOSED
  * sl-wre: US1: Browser-Based Login - CLOSED
  * sl-19x: US2: Token-Based Login (CI/Headless) - CLOSED
  * sl-6vj: US3: Check Authentication Status - CLOSED
  * sl-78q: US4: Logout and Clear Credentials - CLOSED
  * sl-ufh: US5: Token Refresh - CLOSED
* **Tasks**: Issues of type `task`, children of each feature issue (phase) - ALL CLOSED

## Convention Summary

| Type    | Description                  | Labels                                            |
| ------- | ---------------------------- | ------------------------------------------------- |
| epic    | Full feature epic            | `spec:003-authenticate-cli`, `component:cli`      |
| feature | Implementation phase / story | `phase:[n]`, `story:[US#]`                        |
| task    | Implementation task          | `component:auth`, `fr:[FR-XXX]`                   |

## Phase Summary

### Phase 1: Setup (sl-0cm) - COMPLETE

| Task ID | Title | Status |
|---------|-------|--------|
| sl-ocx | Create pkg/cli/auth package structure | CLOSED |
| sl-76u | Add auth command group to CLI | CLOSED |

### Phase 2: Foundational (sl-96q) - COMPLETE

| Task ID | Title | Status |
|---------|-------|--------|
| sl-y1w | Implement Credentials struct and methods | CLOSED |
| sl-387 | Implement credential storage functions | CLOSED |
| sl-yw3 | Implement environment-based URL configuration | CLOSED |

### Phase 3: US1 - Browser-Based Login (sl-wre) - COMPLETE

| Task ID | Title | Status |
|---------|-------|--------|
| sl-ll7 | Implement CallbackServer for OAuth callback | CLOSED |
| sl-5qf | Implement cross-platform browser opening | CLOSED |
| sl-b6f | Implement sl auth login command | CLOSED |

### Phase 4: US2 - Token-Based Login (sl-19x) - COMPLETE

| Task ID | Title | Status |
|---------|-------|--------|
| sl-cl9 | Implement RefreshAccessToken HTTP client | CLOSED |
| sl-fps | Add --token flag for direct access token login | CLOSED |
| sl-afl | Add --refresh flag for refresh token login | CLOSED |

### Phase 5: US3 - Check Status (sl-6vj) - COMPLETE

| Task ID | Title | Status |
|---------|-------|--------|
| sl-idh | Implement sl auth status command | CLOSED |

### Phase 6: US4 - Logout (sl-78q) - COMPLETE

| Task ID | Title | Status |
|---------|-------|--------|
| sl-6mc | Implement sl auth logout command | CLOSED |

### Phase 7: US5 - Token Refresh (sl-ufh) - COMPLETE

| Task ID | Title | Status |
|---------|-------|--------|
| sl-idn | Implement GetValidAccessToken with auto-refresh | CLOSED |
| sl-b04 | Implement sl auth refresh command | CLOSED |

## Implementation Stats

- **Total Tasks**: 14
- **Completed**: 14 (100%)
- **User Stories**: 5 (all implemented)
- **MVP Delivered**: US1 (Browser-Based Login)
- **Full Feature**: All 5 user stories implemented

## Dependencies Flow

```
Setup (sl-0cm)
    └── Foundational (sl-96q)
            ├── US1: Browser Login (sl-wre)
            ├── US2: Token Login (sl-19x)
            ├── US3: Check Status (sl-6vj)
            ├── US4: Logout (sl-78q)
            └── US5: Token Refresh (sl-ufh)
```

## Implementation Files

| Component | File Path | User Stories |
|-----------|-----------|--------------|
| Credentials | pkg/cli/auth/credentials.go | All |
| Callback Server | pkg/cli/auth/server.go | US1 |
| Browser Opener | pkg/cli/auth/browser.go | US1 |
| HTTP Client | pkg/cli/auth/client.go | US2, US5 |
| Auth Commands | pkg/cli/commands/auth.go | All |

---

> This file is an index-only reference. All task data lives in Beads.
