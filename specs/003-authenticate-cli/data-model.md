# Data Model: Authenticate CLI

**Date**: 2026-02-09
**Feature**: [spec.md](./spec.md)

## Entities

### 1. Credentials

Represents stored authentication tokens and user information.

**Location**: `~/.specledger/credentials.json`
**Permissions**: 0600 (owner read/write only)

| Field | Type | Description |
|-------|------|-------------|
| access_token | string | JWT access token for API authentication |
| refresh_token | string | Long-lived token for obtaining new access tokens |
| expires_in | int64 | Token validity duration in seconds |
| created_at | int64 | Unix timestamp when credentials were saved |
| user_email | string | User's email address |
| user_id | string | User's unique identifier |

**Validation Rules**:
- `access_token` required for valid credentials
- `refresh_token` required for auto-refresh capability
- `expires_in` defaults to 3600 (1 hour) if not provided

**State Transitions**:
- `Not Authenticated` → `Authenticated`: After successful login
- `Authenticated` → `Expired`: When current time > created_at + expires_in
- `Expired` → `Authenticated`: After successful token refresh
- `Authenticated` → `Not Authenticated`: After logout

### 2. CallbackResult

Contains authentication result data received from browser callback.

**Usage**: Transient object used during login flow

| Field | Type | Description |
|-------|------|-------------|
| access_token | string | JWT access token |
| refresh_token | string | Refresh token |
| expires_in | int64 | Token validity in seconds |
| user_email | string | User's email |
| user_id | string | User's ID |
| error | string | Error message if authentication failed |

### 3. CallbackServer

Local HTTP server for receiving OAuth callback.

**Configuration**:
- Port: 2026 (default)
- Binding: 127.0.0.1 (localhost only)
- Endpoint: `/callback`

| Field | Type | Description |
|-------|------|-------------|
| port | int | Port number server listens on |
| server | *http.Server | HTTP server instance |
| listener | net.Listener | Network listener |
| result | chan CallbackResult | Channel for passing auth result |
| frontendURL | string | Base URL for redirect after callback |

**Supported Callback Methods**:
- GET: Query parameters (access_token, refresh_token, email, user_id, expires_in, error)
- POST: JSON body with CallbackResult fields
- OPTIONS: CORS preflight support

### 4. RefreshTokenResponse

Response from the token refresh API endpoint.

| Field | Type | Description |
|-------|------|-------------|
| access_token | string | New access token |
| refresh_token | string | Optionally rotated refresh token |
| expires_in | int64 | New token validity in seconds |
| email | string | User's email |
| user_id | string | User's ID |
| error | string | Error code if refresh failed |
| error_description | string | Human-readable error message |

## File Storage Schema

### credentials.json

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "dGhpcyBpcyBhIHJlZnJlc2ggdG9rZW4...",
  "expires_in": 3600,
  "created_at": 1707494400,
  "user_email": "user@example.com",
  "user_id": "usr_abc123"
}
```

## Relationships

```
┌─────────────────┐
│  CallbackServer │
│                 │
│  port: 2026     │
│  frontendURL    │
└────────┬────────┘
         │ receives
         ▼
┌─────────────────┐      ┌─────────────────┐
│ CallbackResult  │      │ RefreshToken    │
│                 │      │ Response        │
│ access_token    │      │                 │
│ refresh_token   │      │ access_token    │
│ user_email      │      │ refresh_token   │
└────────┬────────┘      └────────┬────────┘
         │ transforms             │ transforms
         └──────────┬─────────────┘
                    ▼
         ┌─────────────────┐
         │   Credentials   │
         │                 │
         │  stored at      │
         │  ~/.specledger/ │
         │  credentials.json│
         └─────────────────┘
```

## Helper Functions

### Credentials Methods

| Method | Returns | Description |
|--------|---------|-------------|
| IsExpired() | bool | True if current time > expires_at |
| ExpiresAt() | time.Time | Calculated expiration timestamp |
| IsValid() | bool | True if both access_token and refresh_token are non-empty |

### Storage Functions

| Function | Description |
|----------|-------------|
| GetCredentialsPath() | Returns `~/.specledger/credentials.json` |
| LoadCredentials() | Reads and parses credentials file |
| SaveCredentials() | Writes credentials with 0600 permissions |
| DeleteCredentials() | Removes credentials file |
| GetValidAccessToken() | Returns valid token, refreshing if needed |
| RefreshAccessToken() | Exchanges refresh token for new access token |
