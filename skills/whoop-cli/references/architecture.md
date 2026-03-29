# Architecture Reference

## Execution Flow

```
cmd/whoop/main.go
  └─ commands.Execute()
       └─ NewRootCmd()
            ├─ PersistentPreRunE:
            │   ├─ config.NewManager(cfgDir)  →  load config
            │   ├─ output.New(format)         →  create formatter
            │   ├─ auth.NewFileTokenStore()   →  load token
            │   ├─ auth.NewTokenSource()      →  auto-refreshing token source
            │   ├─ api.NewClient(ctx, ts)     →  authenticated HTTP client
            │   └─ api.New*Service(client)    →  all service implementations
            └─ subcommand.RunE               →  uses state.svc.* + state.formatter
```

## Interface Inventory

| Interface | File | Implementations |
|---|---|---|
| `api.HTTPDoer` | `internal/api/client.go:23` | `*http.Client`, `*oauth2.Transport` |
| `api.UserService` | `internal/api/user.go` | `userService` (prod), `mockUserService` (test) |
| `api.CycleService` | `internal/api/cycle.go` | `cycleService` (prod), `mockCycleService` (test) |
| `api.SleepService` | `internal/api/sleep.go` | `sleepService` (prod), `mockSleepService` (test) |
| `api.RecoveryService` | `internal/api/recovery.go` | `recoveryService` (prod), `mockRecoveryService` (test) |
| `api.WorkoutService` | `internal/api/workout.go` | `workoutService` (prod), `mockWorkoutService` (test) |
| `output.Formatter` | `internal/output/formatter.go` | `TableFormatter`, `JSONFormatter` |
| `config.Manager` | `internal/config/config.go` | `fileManager` |

## API Client

**File**: `internal/api/client.go`

- Base URL: `https://api.prod.whoop.com/developer`
- User-Agent: `whoop-cli/1.0.0`
- Auth: OAuth2 token source auto-injects `Authorization: Bearer ...`
- Rate limiting: on HTTP 429, reads `X-RateLimit-Reset` header (unix epoch), sleeps until then, retries once. Fallback: 500ms fixed wait.
- Errors: non-2xx responses → `*APIError{StatusCode, Message}`
- JSON: responses decoded via `json.NewDecoder(resp.Body).Decode(dest)`

Testing entrypoint: `NewClientWithBaseURL(url, httpDoer)` bypasses OAuth2 setup.

## Authentication

**Files**: `internal/auth/`

- OAuth2 Authorization Code + PKCE (S256)
- Auth URL: `https://api.prod.whoop.com/oauth/oauth2/auth`
- Token URL: `https://api.prod.whoop.com/oauth/oauth2/token`
- Local callback server on `localhost:8484` at `/callback/whoop`
- PKCE: 32 random bytes → base64url verifier → SHA256 → base64url challenge
- State parameter validated in callback (CSRF protection)
- Scopes: `read:profile`, `read:body_measurement`, `read:cycles`, `read:sleep`, `read:recovery`, `read:workout`, `offline`

### Token Storage

**File**: `internal/auth/token_store.go`

- Path: `~/.config/whoop-cli/credentials.json`
- Permissions: `0600` (file), `0700` (directory)
- Stores: `oauth2.Token` (access token, refresh token, expiry, token type)
- `PersistentTokenSource` wraps `oauth2.TokenSource` and saves refreshed tokens to disk

## Configuration

**File**: `internal/config/config.go`

```go
type Config struct {
    OutputFormat string  // "table" or "json"
    ClientID     string  // OAuth2 client ID
    ClientSecret string  // OAuth2 client secret
    RedirectPort int     // callback server port (default: 8484)
}
```

Precedence (highest wins):
1. Environment variables: `WHOOP_CLIENT_ID`, `WHOOP_CLIENT_SECRET`, `WHOOP_OUTPUT`, `WHOOP_REDIRECT_PORT`
2. Config file: `~/.config/whoop-cli/config.yaml`
3. Defaults: `output_format="table"`, `redirect_port=8484`

## Output Pipeline

**Files**: `internal/output/`

- `formatter.go`: `Formatter` interface + `New(format string)` factory
- `table.go`: `TableFormatter` with type-switch dispatch to per-type render functions
- `json.go`: `JSONFormatter` uses `json.MarshalIndent` (2-space indent)

The table formatter falls through to JSON for unknown types.

## Test Infrastructure

**Files**: `internal/testutil/`

### APIServer (`server.go`)

Wraps `httptest.Server` with:
- Method+path routing via `Handle()` / `HandleFunc()`
- Call counting via `CallCount()` / `AssertCalled()`
- Thread-safe via `sync.Mutex`

### Fixtures (`fixtures.go`)

- `baseTime`: `2024-06-01T08:00:00Z` (anchor for all dates)
- `Sample*()` functions: return single fully-populated structs
- `Paginated*()` functions: return `PaginatedResponse[T]` with N items + optional next token

### Command test helpers (`commands/commands_test.go`)

- `setupState(t, &services{...})`: injects mock services, JSON formatter, temp config dir
- `runCmd(t, args)`: executes root command, captures stdout as string
- Mock services: one per data type, return canned results from struct fields
