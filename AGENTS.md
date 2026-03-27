# AGENTS.md

AI agent guidance for whoop-cli — a Go CLI for the WHOOP fitness wearable API v2.

## Tech Stack

- Go 1.25, module `github.com/domo84/whoop-cli`
- cobra (CLI), viper (config), golang.org/x/oauth2, tablewriter, fatih/color
- testify for assertions

## Build & Test

```sh
# Build
go build -o whoop ./cmd/whoop

# Install
go build -o whoop ./cmd/whoop && sudo mv whoop /usr/local/bin/

# All tests
go test ./...

# Single package
go test ./internal/api/...
go test ./commands/...

# Single test
go test ./commands/... -run TestCycleList_Command

# Vet & format
go vet ./...
go fmt ./...
```

## Project Structure

```
cmd/whoop/main.go          # Entry point — calls commands.Execute()
commands/
  root.go                  # rootState struct, PersistentPreRunE wiring, authExemptCommands map
  auth.go                  # login / logout / setup / status
  cycle.go / sleep.go / recovery.go / workout.go / user.go
internal/
  api/
    client.go              # HTTP client; retries once on 429, reads X-RateLimit-Reset
    models.go              # All API types (Cycle, Sleep, Recovery, Workout, etc.)
    cycle.go / sleep.go …  # Service interfaces + implementations
  auth/
    pkce.go                # OAuth2 PKCE flow (S256, 32-byte verifier)
    server.go              # Local callback server on redirect_port
    token_store.go         # credentials.json persistence (0600)
    auth.go                # Auto-refreshing TokenSource
  config/
    config.go              # Manager interface; reads ~/.config/whoop-cli/config.yaml via viper
  output/
    formatter.go           # Formatter interface
    table.go / json.go     # table and json implementations
  testutil/
    server.go              # httptest wrapper with call counting (APIServer)
    fixtures.go            # Shared sample data factories
```

## Architecture

Request flow:
```
cobra command → service interface → api.Client.get() → WHOOP API v2
```

`commands/root.go` `PersistentPreRunE` runs before every command: loads config via `config.Manager`, loads token from `credentials.json`, wraps it in an auto-refreshing `auth.TokenSource`, builds `api.Client`, populates `state` (package-level `rootState`). All subcommands read from `state`.

Adding a new auth-exempt command: add its full command path string to the `authExemptCommands` map in `root.go`.

## WHOOP API Integration

- Base URL: `https://api.prod.whoop.com/developer`
- Auth endpoints: `https://api.prod.whoop.com/oauth/oauth2/{auth,token}`
- Paginated responses: `{ "records": [...], "next_token": "..." }`; max 25 per page
- Rate limit: automatic single retry on HTTP 429, sleeps until `X-RateLimit-Reset`
- Scopes: `read:profile read:body_measurement read:cycles read:sleep read:recovery read:workout offline`

## Config & Credentials

| File | Path | Notes |
|---|---|---|
| Config | `~/.config/whoop-cli/config.yaml` | Fields: `client_id`, `client_secret`, `output_format`, `redirect_port` |
| Token | `~/.config/whoop-cli/credentials.json` | OAuth2 token, 0600, auto-refreshed |

Env vars: `WHOOP_CLIENT_ID`, `WHOOP_CLIENT_SECRET`, `WHOOP_OUTPUT`, `WHOOP_REDIRECT_PORT` — override config file.

## Code Conventions

**Testing — two patterns:**
1. Command tests (`commands/commands_test.go`): inject mock services via `setupState()`, no network, `runCmd()` captures stdout
2. API service tests (`internal/api/*_test.go`): use `testutil.NewAPIServer` (real httptest server) + `api.NewClientWithBaseURL`

**Adding a command:** create a file in `commands/`, follow the pattern in `cycle.go`. Register in `root.go` `NewRootCmd()`. If auth-exempt, add to `authExemptCommands`.

**Adding an API resource:** add types to `internal/api/models.go`, create service interface + implementation file, add `_test.go` using `testutil.APIServer`.

## Security Boundaries

**Always safe (no confirmation needed):**
- Read files, grep, glob
- `go vet ./...`, `go fmt ./...`
- `go test ./internal/...` or `go test ./commands/...` (single package)
- `whoop auth status`

**Ask first:**
- `go build` and installing binary
- `go test ./...` (full suite)
- `go get` / modifying `go.mod`
- Any file writes outside the repo

**Never:**
- Commit `client_secret`, tokens, or `credentials.json`
- Hardcode credentials in source
- Modify OAuth endpoints or disable token validation
- Push directly to `trunk` without review
