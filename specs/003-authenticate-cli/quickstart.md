# Quickstart: CLI Authentication

**Date**: 2026-02-09
**Feature**: [spec.md](./spec.md)

## Overview

The SpecLedger CLI provides authentication commands under `sl auth` for managing user sessions.

## Commands

### Browser Login (Interactive)

```bash
# Standard login - opens browser
sl auth login

# Development mode (localhost:3000)
sl auth login --dev
```

**Flow**:
1. CLI starts local callback server on port 2026
2. Opens browser to SpecLedger sign-in page
3. User completes authentication in browser
4. Browser redirects to localhost callback
5. CLI captures tokens and stores credentials

### Token Login (CI/Headless)

```bash
# Direct access token (for CI environments)
sl auth login --token "your-access-token"

# Refresh token (exchanges for access token)
sl auth login --refresh "your-refresh-token"
```

Use these for:
- CI/CD pipelines
- Docker containers
- Headless servers
- Automated scripts

### Check Status

```bash
sl auth status
```

Output:
```
Status: Signed in
Email:  user@example.com
Token:  Valid (expires in 45m)
Credentials: /Users/you/.specledger/credentials.json
```

### Logout

```bash
sl auth logout
```

Removes `~/.specledger/credentials.json`.

### Refresh Token

```bash
sl auth refresh
```

Manually refreshes the access token using the stored refresh token.

## Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `SPECLEDGER_ENV` | Set to `dev` for development URLs | `dev` |
| `SPECLEDGER_AUTH_URL` | Override authentication URL | `https://custom.auth.example.com` |
| `SPECLEDGER_API_URL` | Override API base URL | `https://api.example.com` |

## For Developers

### Code Structure

```
pkg/cli/auth/
├── browser.go       # OpenBrowser() - cross-platform browser launch
├── client.go        # RefreshAccessToken(), GetValidAccessToken()
├── credentials.go   # Credential storage/loading
└── server.go        # CallbackServer - local OAuth receiver

pkg/cli/commands/
└── auth.go          # Command implementations (login, logout, status, refresh)
```

### Adding Auth to Other Commands

To require authentication in a new command:

```go
import "specledger/pkg/cli/auth"

func myCommand(cmd *cobra.Command, args []string) error {
    // Get valid access token (auto-refreshes if expired)
    token, err := auth.GetValidAccessToken()
    if err != nil {
        return fmt.Errorf("authentication required: %w", err)
    }

    // Use token for API calls
    req.Header.Set("Authorization", "Bearer " + token)
    // ...
}
```

### Testing Authentication Flow

```bash
# Start in dev mode
SPECLEDGER_ENV=dev sl auth login

# Verify credentials stored
cat ~/.specledger/credentials.json | jq .

# Check expiration
sl auth status

# Test refresh
sl auth refresh

# Clean up
sl auth logout
```

### Manual Authentication Fallback

If browser callback fails, users can authenticate manually:

```bash
# Copy access token from browser and use directly
sl auth login --token "paste-token-here"

# Or copy refresh token for longer-lived auth
sl auth login --refresh "paste-refresh-token-here"
```

## Troubleshooting

### Port 2026 Already in Use

The callback server uses port 2026. If it's occupied:
1. Find and stop the process using the port: `lsof -i :2026`
2. Or use manual token authentication: `sl auth login --token <token>`

### Browser Doesn't Open

If the browser doesn't open automatically:
1. Copy the URL shown in terminal
2. Paste in your browser manually
3. Complete authentication
4. The CLI will capture the callback

### Token Expired

Access tokens expire after 1 hour by default. The CLI automatically refreshes when possible. If refresh fails:
```bash
sl auth login  # Re-authenticate via browser
```

### Credentials File Permissions

Credentials are stored with 0600 permissions. If you see permission errors:
```bash
chmod 600 ~/.specledger/credentials.json
```
