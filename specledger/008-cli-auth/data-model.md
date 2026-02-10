# Data Model: CLI Authentication

**Feature**: 008-cli-auth
**Date**: 2026-02-10

## Entities

### Credentials

Represents stored authentication tokens for a user session.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| access_token | string | JWT access token for API authentication | Required, non-empty |
| refresh_token | string | Token used to obtain new access tokens | Required for full auth |
| expires_in | int64 | Seconds until access token expires | Positive integer |
| created_at | int64 | Unix timestamp when credentials were saved | Valid timestamp |
| user_email | string | Email of the authenticated user | Valid email format |
| user_id | string | Unique identifier for the user | UUID format |

**Storage Location**: `~/.specledger/credentials.json`
**File Permissions**: 0600 (owner read/write only)

**State Transitions**:
```
[No Credentials] --login--> [Valid Credentials]
[Valid Credentials] --token expires--> [Expired Credentials]
[Expired Credentials] --refresh--> [Valid Credentials]
[Any State] --logout--> [No Credentials]
```

**Derived Properties**:
- `IsExpired()`: `time.Now() > CreatedAt + ExpiresIn`
- `IsValid()`: `AccessToken != "" && RefreshToken != ""`
- `ExpiresAt()`: `CreatedAt + ExpiresIn` as Time

---

### CallbackResult

Represents the authentication result received from the browser OAuth flow.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| access_token | string | JWT access token | Required on success |
| refresh_token | string | Refresh token | Required on success |
| expires_in | int64 | Token expiration in seconds | Default: 3600 |
| user_email | string | Email of authenticated user | Present on success |
| user_id | string | User identifier | Present on success |
| error | string | Error message if auth failed | Empty on success |

**Transport**: Received via HTTP callback (GET query params or POST JSON body)

---

### CallbackServer

Represents the local HTTP server that handles OAuth callbacks.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| port | int | Port number server listens on | Default: 2026 |
| listener | net.Listener | TCP listener | Bound to 127.0.0.1 |
| server | *http.Server | HTTP server instance | - |
| result | chan CallbackResult | Channel for callback result | Buffered (1) |
| frontendURL | string | Base URL for redirect after callback | Valid URL |

**Lifecycle**:
1. Created with `NewCallbackServer(frontendURL)`
2. Started with `Start()` (background goroutine)
3. Waits with `WaitForCallback(timeout)`
4. Shutdown with `Shutdown()`

---

### RefreshTokenResponse

Represents the API response when refreshing an access token.

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| access_token | string | New access token | Required on success |
| refresh_token | string | New or same refresh token | Optional |
| expires_in | int64 | Token expiration in seconds | Default: 3600 |
| email | string | User email | Optional |
| user_id | string | User identifier | Optional |
| error | string | Error code | Empty on success |
| error_description | string | Human-readable error | Empty on success |

**Transport**: Received as JSON from POST `/api/cli/refresh`

## Relationships

```
┌─────────────────────┐
│   CallbackServer    │
│  (temporary, local) │
└─────────┬───────────┘
          │ receives
          ▼
┌─────────────────────┐
│   CallbackResult    │
│    (in memory)      │
└─────────┬───────────┘
          │ converts to
          ▼
┌─────────────────────┐     refresh     ┌──────────────────────┐
│    Credentials      │◄───────────────►│ RefreshTokenResponse │
│ (~/.specledger/     │                 │    (API response)    │
│  credentials.json)  │                 └──────────────────────┘
└─────────────────────┘
```

## Storage Format

### credentials.json

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "dGhpcyBpcyBhIHJlZnJl...",
  "expires_in": 3600,
  "created_at": 1707580800,
  "user_email": "user@example.com",
  "user_id": "usr_123abc"
}
```
