# Feature Specification: CLI Authentication

**Feature Branch**: `008-cli-auth`
**Created**: 2026-02-10
**Status**: Draft
**Input**: User description: "document lại những gì tôi thay đổi" (Document the changes I made)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Browser-Based Login (Priority: P1)

A developer needs to authenticate their SpecLedger CLI to access protected features like private specifications and remote synchronization. They run `sl auth login` and their default browser opens to the SpecLedger sign-in page. After completing authentication in the browser, the CLI automatically receives the credentials via a local callback server and stores them securely.

**Why this priority**: Authentication is the gateway to all protected features. Without this capability, users cannot access private specs or sync with remote repositories.

**Independent Test**: Can be fully tested by running `sl auth login`, completing the browser flow, and verifying credentials are stored at `~/.specledger/credentials.json`.

**Acceptance Scenarios**:

1. **Given** a user is not authenticated, **When** they run `sl auth login`, **Then** their default browser opens to the SpecLedger authentication page
2. **Given** a user completes browser authentication, **When** the callback is received, **Then** credentials are stored with 0600 permissions in the user's home directory
3. **Given** a user is already authenticated, **When** they run `sl auth login`, **Then** they are informed of current session and can re-authenticate if desired

---

### User Story 2 - Authentication Status Check (Priority: P2)

A developer wants to verify their current authentication status before running protected commands. They run `sl auth status` to see if they are signed in, who they are signed in as, and when their token expires.

**Why this priority**: Users need visibility into their auth state to troubleshoot issues and understand their current access level.

**Independent Test**: Can be tested by running `sl auth status` with and without valid credentials, verifying appropriate output in each case.

**Acceptance Scenarios**:

1. **Given** a user is authenticated with valid credentials, **When** they run `sl auth status`, **Then** they see "Signed in" status, their email, and token expiration time
2. **Given** a user is not authenticated, **When** they run `sl auth status`, **Then** they see "Not signed in" status with instructions to authenticate
3. **Given** a user has expired credentials, **When** they run `sl auth status`, **Then** they see token expiration warning

---

### User Story 3 - Logout (Priority: P2)

A developer wants to sign out of their SpecLedger CLI session, perhaps because they are switching accounts or leaving a shared workstation. They run `sl auth logout` to remove their stored credentials.

**Why this priority**: Users need the ability to securely sign out and clear credentials from their machine.

**Independent Test**: Can be tested by running `sl auth logout` after being authenticated, then verifying credentials file is removed.

**Acceptance Scenarios**:

1. **Given** a user is authenticated, **When** they run `sl auth logout`, **Then** credentials are deleted and confirmation is shown
2. **Given** a user is not authenticated, **When** they run `sl auth logout`, **Then** they see a message indicating they are not currently signed in

---

### User Story 4 - Token Refresh (Priority: P3)

A developer's access token has expired but they have a valid refresh token. The system automatically refreshes the token when needed, or they can manually trigger a refresh with `sl auth refresh`.

**Why this priority**: Seamless token refresh prevents interruption during long sessions and is essential for a good developer experience.

**Independent Test**: Can be tested by running `sl auth refresh` with valid credentials, verifying new token is obtained and stored.

**Acceptance Scenarios**:

1. **Given** a user has expired access token but valid refresh token, **When** they run `sl auth refresh`, **Then** a new access token is obtained and stored
2. **Given** a user has no refresh token, **When** they attempt to refresh, **Then** they see an error message instructing them to login again

---

### User Story 5 - Token-Based Authentication for CI/CD (Priority: P3)

A developer or CI system needs to authenticate without browser interaction using pre-obtained tokens. They can use `sl auth login --token <access_token>` or `sl auth login --refresh <refresh_token>` for headless authentication.

**Why this priority**: CI/CD integration is essential for automated workflows but secondary to interactive developer experience.

**Independent Test**: Can be tested by running `sl auth login --token <token>` and verifying credentials are stored.

**Acceptance Scenarios**:

1. **Given** a valid access token, **When** user runs `sl auth login --token <token>`, **Then** credentials are stored successfully
2. **Given** a valid refresh token, **When** user runs `sl auth login --refresh <token>`, **Then** token is exchanged and credentials are stored

---

### Edge Cases

- What happens when the callback server port (2026) is already in use?
- How does the system handle network timeouts during authentication?
- What happens when the browser cannot be opened automatically?
- How does the system handle corrupted credentials file?
- What happens when the authentication times out (after 5 minutes)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide an `sl auth` command group for managing authentication
- **FR-002**: System MUST support browser-based OAuth login via `sl auth login`
- **FR-003**: System MUST start a local callback server on port 2026 to receive authentication callbacks
- **FR-004**: System MUST store credentials securely at `~/.specledger/credentials.json` with 0600 permissions
- **FR-005**: System MUST support checking authentication status via `sl auth status`
- **FR-006**: System MUST support signing out via `sl auth logout`
- **FR-007**: System MUST support manual token refresh via `sl auth refresh`
- **FR-008**: System MUST automatically refresh expired access tokens using refresh tokens when available
- **FR-009**: System MUST support headless authentication via `--token` and `--refresh` flags for CI/CD environments
- **FR-010**: System MUST support development mode via `--dev` flag or `SPECLEDGER_ENV` environment variable
- **FR-011**: System MUST support custom auth URL via `SPECLEDGER_AUTH_URL` environment variable
- **FR-012**: System MUST open browser automatically on macOS, Linux (xdg-open), and Windows
- **FR-013**: System MUST display manual authentication URL if browser cannot be opened
- **FR-014**: System MUST handle callback via both GET (query params) and POST (JSON body) methods
- **FR-015**: System MUST set appropriate CORS headers for browser-based callbacks
- **FR-016**: System MUST redirect browser back to frontend after callback with success/error status

### Key Entities

- **Credentials**: Stored authentication data including access token, refresh token, expiration time, creation timestamp, user email, and user ID
- **CallbackResult**: Authentication result from browser including tokens, expiration, user info, and potential error
- **CallbackServer**: Local HTTP server handling OAuth callback with configurable port and frontend redirect URL

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can complete browser-based authentication in under 30 seconds (excluding identity provider time)
- **SC-002**: Authentication credentials persist across CLI sessions until logout or expiration
- **SC-003**: Token refresh happens automatically without user intervention when refresh token is valid
- **SC-004**: CI/CD systems can authenticate without browser using provided tokens
- **SC-005**: Authentication status is clearly displayed with relevant information (email, expiration)
- **SC-006**: Authentication works across macOS, Linux, and Windows platforms

### Previous work

No previous related features or tasks found in this repository.
