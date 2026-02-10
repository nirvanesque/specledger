# Quickstart: CLI Authentication

## Overview

SpecLedger CLI authentication enables access to protected features like private specifications and remote synchronization.

## Prerequisites

- SpecLedger CLI installed (`sl` command available)
- Modern web browser for authentication
- Network access to app.specledger.io

## Basic Usage

### Sign In

```bash
# Open browser for authentication
sl auth login

# The CLI will:
# 1. Start a local callback server on port 2026
# 2. Open your default browser to the sign-in page
# 3. Wait for you to complete authentication
# 4. Store credentials securely
```

### Check Status

```bash
sl auth status
# Output:
# Status: Signed in
# Email:  user@example.com
# Token:  Valid (expires in 58m)
```

### Sign Out

```bash
sl auth logout
# Output:
# Signed out successfully. (was: user@example.com)
```

### Refresh Token

```bash
sl auth refresh
# Output:
# Token refreshed successfully!
# Email:   user@example.com
# Expires: in 1h
```

## CI/CD Integration

For automated environments without browser access:

```bash
# Using access token directly
sl auth login --token "$SPECLEDGER_ACCESS_TOKEN"

# Using refresh token (exchanges for access token)
sl auth login --refresh "$SPECLEDGER_REFRESH_TOKEN"
```

## Development Mode

For local development against localhost:

```bash
# Use --dev flag
sl auth login --dev

# Or set environment variable
SPECLEDGER_ENV=dev sl auth login

# Or use custom URL
SPECLEDGER_AUTH_URL=http://localhost:3000/cli/auth sl auth login
```

## Troubleshooting

### Browser doesn't open

If the browser doesn't open automatically, manually navigate to the URL shown in the terminal.

### Callback fails

If the callback fails, use manual token authentication:
```bash
# Copy the token from the browser after authentication
sl auth login --token <your_access_token>
```

### Port 2026 in use

The CLI uses port 2026 for the callback server. If this port is in use, close the conflicting process.

### Token expired

Tokens automatically refresh when expired. To manually refresh:
```bash
sl auth refresh
```

## Credential Storage

Credentials are stored at `~/.specledger/credentials.json` with restricted permissions (0600 - owner read/write only).

```bash
# View stored credentials location
sl auth status
# Credentials: /Users/you/.specledger/credentials.json
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `SPECLEDGER_ENV` | Set to `dev` or `development` for localhost auth |
| `SPECLEDGER_AUTH_URL` | Custom authentication URL |
| `SPECLEDGER_API_URL` | Custom API URL for token refresh |
