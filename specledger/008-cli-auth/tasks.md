# Tasks Index: CLI Authentication

Beads Issue Graph Index into the tasks and phases for this feature implementation.
This index does **not contain tasks directly**—those are fully managed through Beads CLI.

## Feature Tracking

* **Beads Epic ID**: `SL-92m` (CLOSED - Implementation Complete)
* **User Stories Source**: `specledger/008-cli-auth/spec.md`
* **Research Inputs**: `specledger/008-cli-auth/research.md`
* **Planning Details**: `specledger/008-cli-auth/plan.md`
* **Data Model**: `specledger/008-cli-auth/data-model.md`
* **Contract Definitions**: `specledger/008-cli-auth/contracts/`

## Implementation Status: ✅ COMPLETE

All 13 tasks across 6 phases have been implemented and closed.

| Phase | Feature ID | Status | Tasks |
|-------|------------|--------|-------|
| Setup | SL-cph | ✅ Closed | 2/2 complete |
| US1: Browser-Based Login | SL-0vt | ✅ Closed | 4/4 complete |
| US2: Authentication Status | SL-2ka | ✅ Closed | 1/1 complete |
| US3: Logout | SL-a68 | ✅ Closed | 1/1 complete |
| US4: Token Refresh | SL-cxm | ✅ Closed | 2/2 complete |
| US5: CI/CD Token Auth | SL-li2 | ✅ Closed | 3/3 complete |

## Beads Query Hints

Use the `bd` CLI to query and manipulate the issue graph:

```bash
# Find all tasks for this feature (all closed)
bd list --label spec:008-cli-auth --limit 20

# See the complete dependency tree
bd dep tree --reverse SL-92m

# View by user story
bd list --label spec:008-cli-auth --label story:US1

# View by component
bd list --label spec:008-cli-auth --label component:auth
```

## Tasks and Phases Structure

This feature follows Beads' 2-level graph structure:

* **Epic**: SL-92m → CLI Authentication (CLOSED)
* **Phases**: Features as children of the epic
  * SL-cph: Setup Phase (CLOSED)
  * SL-0vt: US1 - Browser-Based Login (CLOSED)
  * SL-2ka: US2 - Authentication Status (CLOSED)
  * SL-a68: US3 - Logout (CLOSED)
  * SL-cxm: US4 - Token Refresh (CLOSED)
  * SL-li2: US5 - CI/CD Token Auth (CLOSED)
* **Tasks**: Issues of type `task`, children of each feature issue

## Convention Summary

| Type    | Description                  | Labels                                       |
| ------- | ---------------------------- | -------------------------------------------- |
| epic    | Full feature epic            | `spec:008-cli-auth`, `component:cli`         |
| feature | Implementation phase / story | `phase:*`, `story:US*`                       |
| task    | Implementation task          | `component:auth`, `fr:FR-*`                  |

---

## Phase 1: Setup ✅ COMPLETE

**Feature ID**: SL-cph (CLOSED)

| Task ID | Title | Status |
|---------|-------|--------|
| SL-7iv | Create auth package structure | ✅ Closed |
| SL-6v3 | Add auth command to main.go | ✅ Closed |

**Files Created**:
- `pkg/cli/auth/` directory
- `pkg/cli/commands/auth.go`
- Modified `cmd/sl/main.go`

---

## Phase 2: US1 - Browser-Based Login (P1) ✅ COMPLETE 🎯 MVP

**Feature ID**: SL-0vt (CLOSED)

**Goal**: Enable developers to authenticate via browser-based OAuth flow

**Implemented Files**:
- `pkg/cli/auth/credentials.go` - Credential storage
- `pkg/cli/auth/server.go` - Callback server
- `pkg/cli/auth/browser.go` - Cross-platform browser opening
- `pkg/cli/commands/auth.go` - Login command

| Task ID | Title | FR | Status |
|---------|-------|-----|--------|
| SL-s0f | Implement Credentials struct and storage | FR-004 | ✅ Closed |
| SL-ya8 | Implement CallbackServer | FR-003 | ✅ Closed |
| SL-ipy | Implement browser opening | FR-012 | ✅ Closed |
| SL-s6b | Implement sl auth login command | FR-002 | ✅ Closed |

---

## Phase 3: US2 - Authentication Status (P2) ✅ COMPLETE

**Feature ID**: SL-2ka (CLOSED)

**Goal**: Allow users to check their current authentication state

| Task ID | Title | FR | Status |
|---------|-------|-----|--------|
| SL-21q | Implement sl auth status command | FR-005 | ✅ Closed |

---

## Phase 4: US3 - Logout (P2) ✅ COMPLETE

**Feature ID**: SL-a68 (CLOSED)

**Goal**: Enable users to sign out and remove stored credentials

| Task ID | Title | FR | Status |
|---------|-------|-----|--------|
| SL-4ab | Implement sl auth logout command | FR-006 | ✅ Closed |

---

## Phase 5: US4 - Token Refresh (P3) ✅ COMPLETE

**Feature ID**: SL-cxm (CLOSED)

**Goal**: Implement automatic and manual token refresh

**Implemented Files**:
- `pkg/cli/auth/client.go` - API client for token refresh

| Task ID | Title | FR | Status |
|---------|-------|-----|--------|
| SL-hla | Implement API client for token refresh | FR-008 | ✅ Closed |
| SL-dyz | Implement sl auth refresh command | FR-007 | ✅ Closed |

---

## Phase 6: US5 - CI/CD Token Auth (P3) ✅ COMPLETE

**Feature ID**: SL-li2 (CLOSED)

**Goal**: Enable headless authentication for CI/CD environments

| Task ID | Title | FR | Status |
|---------|-------|-----|--------|
| SL-k8d | Add --token flag for headless auth | FR-009 | ✅ Closed |
| SL-d5n | Add --refresh flag for headless auth | FR-009 | ✅ Closed |
| SL-gyk | Add --dev flag and environment support | FR-010, FR-011 | ✅ Closed |

---

## Implementation Summary

### Delivered Functionality

1. **Browser-Based Login** (`sl auth login`)
   - Opens browser to SpecLedger sign-in page
   - Local callback server on port 2026
   - Automatic credential storage

2. **Authentication Status** (`sl auth status`)
   - Shows signed in/out state
   - Displays email and token expiration

3. **Logout** (`sl auth logout`)
   - Removes stored credentials
   - Confirmation message

4. **Token Refresh** (`sl auth refresh`)
   - Manual token refresh
   - Automatic refresh on expiration

5. **CI/CD Support**
   - `--token` flag for access token auth
   - `--refresh` flag for refresh token auth
   - `--dev` flag for development mode
   - Environment variable support

### Files Implemented

```
cmd/sl/main.go                    # Modified - added VarAuthCmd
pkg/cli/auth/
├── browser.go                    # Cross-platform browser opening
├── client.go                     # API client for token refresh
├── credentials.go                # Credential storage management
└── server.go                     # OAuth callback server
pkg/cli/commands/auth.go          # Auth command definitions
```

### Functional Requirements Coverage

| FR | Description | Status |
|----|-------------|--------|
| FR-001 | Auth command group | ✅ |
| FR-002 | Browser-based login | ✅ |
| FR-003 | Callback server port 2026 | ✅ |
| FR-004 | Secure credential storage | ✅ |
| FR-005 | Auth status command | ✅ |
| FR-006 | Logout command | ✅ |
| FR-007 | Manual token refresh | ✅ |
| FR-008 | Automatic token refresh | ✅ |
| FR-009 | Headless auth flags | ✅ |
| FR-010 | Dev mode flag | ✅ |
| FR-011 | Custom auth URL env | ✅ |
| FR-012 | Cross-platform browser | ✅ |
| FR-013 | Manual URL fallback | ✅ |
| FR-014 | GET/POST callback | ✅ |
| FR-015 | CORS headers | ✅ |
| FR-016 | Frontend redirect | ✅ |

---

> This file is an index only. Implementation data lives in Beads. All tasks are complete.
