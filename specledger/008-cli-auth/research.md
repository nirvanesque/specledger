# Research: CLI Authentication

**Feature**: 008-cli-auth
**Date**: 2026-02-10

## Prior Work

No previous authentication features exist in the SpecLedger CLI. This is the first implementation of user authentication.

## Technical Decisions

### 1. Authentication Flow

**Decision**: Browser-based OAuth with local callback server

**Rationale**:
- Standard pattern for CLI tools (GitHub CLI, Vercel CLI, etc.)
- No need to handle password input in terminal
- Delegates authentication complexity to web service
- Better security - credentials never pass through CLI

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|-----------------|
| Password input | Security concerns, need to handle hashing |
| API key only | Less user-friendly, no refresh capability |
| Device code flow | More complex, requires polling |

### 2. Callback Server Port

**Decision**: Port 2026 (DefaultCallbackPort)

**Rationale**:
- Memorable year-based port
- Above privileged port range (>1024)
- Unlikely to conflict with common services

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|-----------------|
| Random port | Harder to configure firewall, less predictable |
| Port 8080 | Common port, likely conflicts |
| Port 3000 | Often used by dev servers |

### 3. Credential Storage

**Decision**: JSON file at `~/.specledger/credentials.json` with 0600 permissions

**Rationale**:
- Follows XDG-like conventions (home directory config)
- JSON is human-readable for debugging
- 0600 ensures only owner can read/write
- Consistent with other CLI tools (gh, docker, etc.)

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|-----------------|
| System keychain | Platform-specific complexity (macOS Keychain, Windows Credential Manager, Linux Secret Service) |
| Environment variables | Not persistent, security concerns |
| Encrypted file | Added complexity, key management issues |

### 4. Token Refresh Strategy

**Decision**: Automatic refresh on access with manual `sl auth refresh` option

**Rationale**:
- Seamless user experience for most cases
- Manual option for debugging/testing
- Refresh token stored alongside access token

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|-----------------|
| Manual-only refresh | Poor UX, requires user intervention |
| Background refresh daemon | Over-engineered for CLI tool |
| Proactive refresh | Complexity, race conditions |

### 5. CI/CD Support

**Decision**: `--token` and `--refresh` flags for headless authentication

**Rationale**:
- Common pattern for CI/CD integration
- Allows automation without browser
- Compatible with secret management systems

**Alternatives Considered**:
| Alternative | Rejected Because |
|-------------|-----------------|
| Environment variable only | Less flexible, harder to debug |
| Config file in repo | Security risk if committed |
| SSH-like key pairs | Different auth paradigm, more complex |

### 6. Cross-Platform Browser Opening

**Decision**: Platform-specific commands (open, xdg-open, cmd /c start)

**Rationale**:
- Uses native browser opening on each platform
- No external dependencies required
- Fallback to manual URL if browser fails

**Implementation**:
- macOS: `open <url>`
- Linux: `xdg-open <url>`
- Windows: `cmd /c start <url>`

## API Contract Decisions

### Callback Server Endpoints

**POST/GET /callback**: Receives authentication result
- Supports both query params (GET) and JSON body (POST)
- Sets CORS headers for browser compatibility
- Redirects to frontend with status

### Refresh Endpoint

**POST /api/cli/refresh**: Exchange refresh token for new access token
- Request: `{"refresh_token": "..."}`
- Response: `{"access_token": "...", "refresh_token": "...", "expires_in": 3600, ...}`

## Security Considerations

1. **Credential file permissions**: 0600 (owner read/write only)
2. **Callback server**: Bound to 127.0.0.1 only (localhost)
3. **Token expiration**: 1 hour default, with refresh token for renewal
4. **No secrets in logs**: Token values not logged, only lengths shown

## Unresolved Questions

None - all technical decisions have been made and implemented.
