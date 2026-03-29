# Adding a New API Endpoint

Use `internal/api/cycle.go` as the canonical template.

## Step 1: Add model structs

In `internal/api/models.go`, add the data types. Follow existing conventions:

```go
// <Type> represents a WHOOP <type> record.
type <Type> struct {
    ID             int        `json:"id"`
    UserID         int        `json:"user_id"`
    CreatedAt      time.Time  `json:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at"`
    Start          time.Time  `json:"start"`
    End            *time.Time `json:"end,omitempty"`      // nullable = pointer
    TimezoneOffset string     `json:"timezone_offset"`
    ScoreState     string     `json:"score_state"`
    Score          *<Type>Score `json:"score,omitempty"`   // nullable = pointer
}

type <Type>Score struct {
    // Use *float64 for optional numeric fields
    SomeMetric *float64 `json:"some_metric,omitempty"`
}
```

Conventions:
- `time.Time` for required timestamps, `*time.Time` for nullable
- `*float64` for optional numeric scores
- Nested score struct named `<Type>Score`
- JSON tags match API response field names exactly
- String IDs use `string`, numeric IDs use `int`

## Step 2: Create the service file

Create `internal/api/<name>.go`:

```go
package api

import (
    "context"
    "fmt"
)

// <Type>Service defines operations for the /v2/<path>/ endpoints.
type <Type>Service interface {
    List(ctx context.Context, params ListParams) (*PaginatedResponse[<Type>], error)
    ListAll(ctx context.Context, params ListParams) ([]<Type>, error)
    Get(ctx context.Context, id int) (*<Type>, error)
}

type <name>Service struct{ client *Client }

// New<Type>Service creates a <Type>Service backed by client.
func New<Type>Service(client *Client) <Type>Service {
    return &<name>Service{client: client}
}

func (s *<name>Service) List(ctx context.Context, params ListParams) (*PaginatedResponse[<Type>], error) {
    var out PaginatedResponse[<Type>]
    if err := s.client.get(ctx, "/v2/<path>", buildQuery(params), &out); err != nil {
        return nil, err
    }
    return &out, nil
}

func (s *<name>Service) ListAll(ctx context.Context, params ListParams) ([]<Type>, error) {
    var all []<Type>
    for {
        page, err := s.List(ctx, params)
        if err != nil {
            return nil, err
        }
        all = append(all, page.Records...)
        if page.NextToken == "" {
            break
        }
        params.NextToken = page.NextToken
    }
    return all, nil
}

func (s *<name>Service) Get(ctx context.Context, id int) (*<Type>, error) {
    var out <Type>
    if err := s.client.get(ctx, fmt.Sprintf("/v2/<path>/%d", id), nil, &out); err != nil {
        return nil, err
    }
    return &out, nil
}
```

For string IDs, change `Get(ctx context.Context, id string)` and use `"/v2/<path>/"+id`.

## Step 3: Wire the service

In `commands/root.go`:

1. Add field to `services` struct:
   ```go
   <Name> api.<Type>Service
   ```

2. Initialize in `PersistentPreRunE`:
   ```go
   <Name>: api.New<Type>Service(client),
   ```

## API client internals

The `Client` (`internal/api/client.go`) handles:
- **Base URL**: `https://api.prod.whoop.com/developer`
- **Auth**: OAuth2 token source adds `Authorization: Bearer ...` automatically
- **Rate limiting**: `doWithRetry` detects 429, waits for `X-RateLimit-Reset`, retries once
- **Error handling**: Non-2xx responses return `*APIError{StatusCode, Message}`
- **Decoding**: JSON response decoded into the `dest` parameter

Available client methods:
- `s.client.get(ctx, path, query, &dest)` -- GET with query params
- `s.client.post(ctx, path, bodyReader)` -- POST (used for token revocation)
- `buildQuery(params)` -- converts `ListParams` to `url.Values`
