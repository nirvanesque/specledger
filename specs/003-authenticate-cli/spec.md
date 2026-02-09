# Feature Specification: Authenticate CLI

**Feature Branch**: `feat-003/authenticate-cli`
**Created**: 2026-02-09
**Status**: Implemented
**Input**: Browser-based user authentication and token management for SpecLedger CLI

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Browser-Based Login (Priority: P1)

A user wants to authenticate their SpecLedger CLI with their SpecLedger account by opening their browser, signing in, and having the credentials automatically captured and stored locally.

**Why this priority**: This is the primary authentication flow that most users will use. It provides a secure, familiar OAuth-style experience without requiring users to manually copy tokens.

**Independent Test**: Can be fully tested by running `sl auth login`, completing browser authentication, and verifying credentials are stored securely. Delivers authenticated access to protected CLI features.

**Acceptance Scenarios**:

1. **Given** a user is not authenticated, **When** they run `sl auth login`, **Then** their default browser opens to the SpecLedger sign-in page
2. **Given** a browser authentication flow is initiated, **When** the user completes sign-in in the browser, **Then** credentials are automatically captured via local callback server and stored securely
3. **Given** the browser cannot open automatically, **When** authentication is initiated, **Then** the CLI displays the authentication URL for manual navigation
4. **Given** a user is already authenticated, **When** they run `sl auth login` again, **Then** they are informed of existing session and can re-authenticate if desired

---

### User Story 2 - Token-Based Login for CI/Headless (Priority: P2)

A user operating in a CI/CD pipeline or headless environment wants to authenticate using pre-obtained tokens without requiring browser interaction.

**Why this priority**: Enables automation and CI/CD integration, which is critical for enterprise workflows and DevOps practices.

**Independent Test**: Can be tested by running `sl auth login --token <access_token>` or `sl auth login --refresh <refresh_token>` and verifying credentials are stored correctly.

**Acceptance Scenarios**:

1. **Given** a user has an access token, **When** they run `sl auth login --token <token>`, **Then** credentials are stored with the provided access token
2. **Given** a user has a refresh token, **When** they run `sl auth login --refresh <token>`, **Then** the system exchanges it for an access token via API and stores valid credentials
3. **Given** an invalid refresh token is provided, **When** authentication is attempted, **Then** a clear error message is displayed indicating the failure reason

---

### User Story 3 - Check Authentication Status (Priority: P2)

A user wants to verify their current authentication status and see their account information.

**Why this priority**: Essential for users to understand their current session state and troubleshoot authentication issues.

**Independent Test**: Can be tested by running `sl auth status` and verifying it shows correct status (signed in/out), email, and token expiration.

**Acceptance Scenarios**:

1. **Given** a user is authenticated with valid tokens, **When** they run `sl auth status`, **Then** they see "Signed in" status, their email, and token expiration time
2. **Given** a user is not authenticated, **When** they run `sl auth status`, **Then** they see "Not signed in" status and instructions to login
3. **Given** a user has expired tokens, **When** they run `sl auth status`, **Then** they see their status with a note that token will refresh on next request

---

### User Story 4 - Logout and Clear Credentials (Priority: P3)

A user wants to sign out and remove all stored credentials from their machine.

**Why this priority**: Security feature that allows users to clear sensitive data when needed.

**Independent Test**: Can be tested by running `sl auth logout` and verifying credentials file is removed.

**Acceptance Scenarios**:

1. **Given** a user is authenticated, **When** they run `sl auth logout`, **Then** credentials are deleted and confirmation message shows previous email
2. **Given** a user is not authenticated, **When** they run `sl auth logout`, **Then** they see a message that they are not currently signed in

---

### User Story 5 - Token Refresh (Priority: P3)

A user or system component wants to manually refresh an expired access token using the stored refresh token.

**Why this priority**: Supports long-running sessions and provides manual control over token refresh.

**Independent Test**: Can be tested by running `sl auth refresh` with valid refresh token and verifying new access token is obtained.

**Acceptance Scenarios**:

1. **Given** a user has valid credentials with refresh token, **When** they run `sl auth refresh`, **Then** a new access token is obtained and stored
2. **Given** access token is expired but refresh token is valid, **When** any authenticated CLI command runs, **Then** token is automatically refreshed

---

### Edge Cases

- What happens when the local callback server port (2026) is already in use?
- How does the system handle browser authentication timeout (5 minutes)?
- What happens when network connectivity is lost during authentication?
- How does the system handle corrupted credentials file?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide an `sl auth` command group for managing authentication
- **FR-002**: System MUST support browser-based OAuth-style authentication flow via `sl auth login`
- **FR-003**: System MUST start a local callback server on port 2026 to receive authentication tokens from browser
- **FR-004**: System MUST support direct token authentication via `--token` flag for CI/headless environments
- **FR-005**: System MUST support refresh token authentication via `--refresh` flag
- **FR-006**: System MUST store credentials securely at `~/.specledger/credentials.json` with 0600 permissions
- **FR-007**: System MUST support `sl auth logout` to sign out and delete credentials
- **FR-008**: System MUST support `sl auth status` to display current authentication state
- **FR-009**: System MUST support `sl auth refresh` to manually refresh access tokens
- **FR-010**: System MUST automatically refresh expired access tokens when refresh token is available
- **FR-011**: System MUST support cross-platform browser opening (macOS, Linux, Windows)
- **FR-012**: System MUST support environment-based configuration (production vs development URLs)
- **FR-013**: System MUST display authentication timeout countdown (5 minutes)
- **FR-014**: System MUST provide manual authentication instructions when browser callback fails

### Key Entities

- **Credentials**: Represents stored authentication data including access_token, refresh_token, expires_in, created_at, user_email, user_id
- **CallbackServer**: Local HTTP server that handles OAuth callback from browser on port 2026
- **CallbackResult**: Authentication result containing tokens and user info from browser callback

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can complete browser authentication in under 60 seconds after sign-in
- **SC-002**: Token-based authentication completes in under 5 seconds
- **SC-003**: Credentials are stored with restricted file permissions (0600)
- **SC-004**: Authentication status check completes in under 1 second
- **SC-005**: Token refresh operation completes in under 10 seconds
- **SC-006**: CLI provides clear feedback at each step of authentication flow
- **SC-007**: Manual authentication fallback is available when automatic callback fails

### Previous work

This is the first authentication feature for the SpecLedger CLI. No previous related work exists.
