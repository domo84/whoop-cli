# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
# Build
go build -o whoop ./cmd/whoop

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/api/...
go test ./commands/...

# Run a single test
go test ./commands/... -run TestCycleList_Command

# Install binary to PATH
go build -o whoop ./cmd/whoop && sudo mv whoop /usr/local/bin/
```

## Architecture

The entry point is `cmd/whoop/main.go`, which calls `commands.Execute()`. All CLI logic lives in the `commands/` package, with business logic split into `internal/` packages.

### Request flow

```
cobra command → service interface → api.Client → WHOOP API v2
```

`commands/root.go` wires everything together in `PersistentPreRunE`: it loads config, loads the OAuth2 token from disk, wraps it in an auto-refreshing `TokenSource`, and builds an `api.Client`. The result is stored in a package-level `state` struct that all subcommands share.

### Key packages

- **`commands/`** — Cobra commands. Each resource (user, cycle, sleep, recovery, workout, auth) has its own file. `root.go` defines the shared `rootState` struct and the `authExemptCommands` map. New auth-exempt commands must be added to that map.
- **`internal/config/`** — `Manager` interface backed by Viper. Reads `~/.config/whoop-cli/config.yaml` and binds env vars (`WHOOP_CLIENT_ID`, `WHOOP_CLIENT_SECRET`, `WHOOP_OUTPUT`, `WHOOP_REDIRECT_PORT`). `Save()` writes config back to disk.
- **`internal/auth/`** — OAuth2 PKCE flow (`pkce.go`), local callback server (`server.go`), token persistence (`token_store.go` → `~/.config/whoop-cli/credentials.json`), and auto-refreshing `TokenSource` (`auth.go`).
- **`internal/api/`** — Typed API client. `client.go` handles HTTP with one retry on 429 (reads `X-RateLimit-Reset`). Each resource has its own file with a service interface and implementation. `models.go` contains all API types.
- **`internal/output/`** — `Formatter` interface with `table` and `json` implementations.
- **`internal/testutil/`** — Shared test helpers: `APIServer` (httptest wrapper with call counting) and fixture factories (`fixtures.go`).

### Testing patterns

Command tests in `commands/commands_test.go` inject mock services via `setupState()` — they never hit the network. API service tests in `internal/api/*_test.go` use `testutil.NewAPIServer` (a real `httptest.Server`) and `api.NewClientWithBaseURL` to point at it.
