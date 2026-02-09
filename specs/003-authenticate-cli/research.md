# Research: Authenticate CLI

**Date**: 2026-02-09
**Feature**: [spec.md](./spec.md)

## Prior Work

No previous authentication features exist in the SpecLedger CLI. This is the foundational auth implementation.

## Technical Decisions

### 1. OAuth-Style Browser Authentication Flow

**Decision**: Implement local callback server pattern for browser-based authentication

**Rationale**:
- Familiar flow for developers (similar to GitHub CLI, Firebase CLI)
- No need to manually copy tokens
- Secure - tokens never exposed in URLs
- Automatic capture of credentials after browser sign-in

**Alternatives Considered**:
- Device code flow: Rejected - requires polling, more complex UX
- Manual token copy: Rejected - poor UX, error-prone
- API key only: Rejected - less secure, no refresh capability

### 2. Callback Server Port Selection

**Decision**: Use fixed port 2026 for callback server

**Rationale**:
- Port 2026 is unlikely to conflict with common services
- Fixed port simplifies frontend configuration
- Easy to remember (SpecLedger founding year reference)

**Alternatives Considered**:
- Random port: Rejected - complicates callback URL configuration
- Port 0 (OS-assigned): Rejected - requires dynamic frontend callback URL
- Well-known ports: Rejected - likely conflicts

### 3. Credential Storage Location

**Decision**: Store credentials at `~/.specledger/credentials.json` with 0600 permissions

**Rationale**:
- Follows Unix conventions for user config
- Permissions prevent other users from reading tokens
- JSON format is human-readable for debugging
- Single file simplifies management

**Alternatives Considered**:
- System keychain: Rejected - adds platform-specific complexity
- Environment variables only: Rejected - no persistence
- Encrypted file: Rejected - key management complexity, diminishing returns since file is already 0600

### 4. Token-Based Authentication for CI/Headless

**Decision**: Support `--token` and `--refresh` flags for non-interactive authentication

**Rationale**:
- Essential for CI/CD pipelines
- Allows pre-provisioned tokens from secure secret managers
- `--refresh` enables token rotation without re-authentication

**Alternatives Considered**:
- Environment variable only: Rejected - flags provide explicit control
- Service account keys: Rejected - over-engineering for current scope

### 5. Cross-Platform Browser Opening

**Decision**: Use OS-specific commands (open, xdg-open, cmd /c start)

**Rationale**:
- Simple, reliable, no dependencies
- Works on all major platforms
- Graceful fallback to manual URL display

**Alternatives Considered**:
- Third-party browser package: Rejected - unnecessary dependency
- Hardcoded browser paths: Rejected - fragile, platform-specific

### 6. Token Refresh Strategy

**Decision**: Automatic refresh on demand with manual refresh command

**Rationale**:
- Seamless UX - users don't need to think about token expiry
- Manual command for explicit control when needed
- Refresh token stored alongside access token

**Alternatives Considered**:
- Proactive refresh (background): Rejected - complexity, not needed for CLI
- No auto-refresh: Rejected - poor UX requiring frequent re-login

## API Endpoints

### Production
- Auth URL: `https://app.specledger.io/cli/auth`
- Refresh API: `https://app.specledger.io/api/cli/refresh`

### Development
- Auth URL: `http://localhost:3000/cli/auth`
- Refresh API: `http://localhost:3000/api/cli/refresh`

### Environment Configuration
- `SPECLEDGER_AUTH_URL`: Override auth URL
- `SPECLEDGER_API_URL`: Override API URL
- `SPECLEDGER_ENV=dev`: Use development URLs

## Security Considerations

1. **Credential File Permissions**: 0600 ensures only owner can read
2. **Local Callback**: 127.0.0.1 binding prevents network access
3. **CORS Headers**: Allow browser to complete callback
4. **Timeout**: 5-minute limit prevents abandoned auth sessions
5. **No Token Logging**: Tokens never printed to stdout/stderr (except truncated length)
