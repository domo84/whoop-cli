# whoop-cli

A command-line tool for accessing your [WHOOP](https://www.whoop.com) fitness data from the terminal via the WHOOP API v2.

## Prerequisites

- Go 1.25+
- A WHOOP developer application — register at [developer.whoop.com](https://developer.whoop.com) to obtain a **Client ID** and **Client Secret**

## Installation

```sh
go install github.com/domo84/whoop-cli/cmd/whoop@latest
```

The `whoop` binary is installed to `$GOPATH/bin` (typically `~/go/bin`) — make sure that's on your `$PATH`.

**Build from source:**

```sh
git clone https://github.com/domo84/whoop-cli.git
cd whoop-cli
go build -o whoop ./cmd/whoop
sudo mv whoop /usr/local/bin/
```

## Quick Start

```sh
# 1. Configure credentials interactively
whoop auth setup

# 2. Authenticate via browser (OAuth2)
whoop auth login

# 3. Fetch your data
whoop recovery list
```

## Commands

### `whoop auth`

| Subcommand | Description |
|---|---|
| `setup` | Interactively configure `config.yaml` (client ID, secret, output format, port) |
| `login` | Open browser to complete OAuth2 authentication |
| `logout` | Revoke token remotely and delete local credentials |
| `status` | Show whether the stored token is valid or expired |

### `whoop user`

| Subcommand | Description |
|---|---|
| `profile` | Name, email, user ID |
| `measurements` | Height, weight, max heart rate |

### `whoop cycle`

| Subcommand | Description |
|---|---|
| `list` | List physiological cycles |
| `get <id>` | Get a single cycle by numeric ID |
| `sleep <cycleId>` | Get the sleep record associated with a cycle |

### `whoop sleep`

| Subcommand | Description |
|---|---|
| `list` | List sleep records (includes naps) |
| `get <id>` | Get a single sleep record by ID |

### `whoop recovery`

| Subcommand | Description |
|---|---|
| `list` | List recovery scores |
| `get <cycleId>` | Get the recovery score for a specific cycle |

### `whoop workout`

| Subcommand | Description |
|---|---|
| `list` | List workouts |
| `get <id>` | Get a single workout by ID |

**List command flags** (available on all `list` subcommands):

| Flag | Default | Description |
|---|---|---|
| `--limit` | 25 | Records per page (max 25) |
| `--start` | — | Start date filter (`YYYY-MM-DD` or RFC3339) |
| `--end` | — | End date filter (`YYYY-MM-DD` or RFC3339) |
| `--all` | false | Fetch all pages |

## Global Flags

| Flag | Default | Description |
|---|---|---|
| `--config-dir` | `~/.config/whoop-cli` | Override config directory |
| `-o, --output` | `table` | Output format: `table` or `json` |
| `--all` | false | Fetch all pages (list commands) |

## Configuration

`~/.config/whoop-cli/config.yaml`:

```yaml
client_id: <your-client-id>
client_secret: <your-client-secret>
output_format: table        # table or json
redirect_port: 8282         # local OAuth2 callback port
```

Environment variables override config file values:

| Variable | Config field |
|---|---|
| `WHOOP_CLIENT_ID` | `client_id` |
| `WHOOP_CLIENT_SECRET` | `client_secret` |
| `WHOOP_OUTPUT` | `output_format` |
| `WHOOP_REDIRECT_PORT` | `redirect_port` |

## Credentials

OAuth2 tokens are stored at `~/.config/whoop-cli/credentials.json` (permissions `0600`). Tokens are refreshed automatically on expiry. Run `whoop auth logout` to revoke and delete them.
