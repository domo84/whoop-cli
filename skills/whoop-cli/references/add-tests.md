# Writing Tests

Three test layers, each with established patterns.

## 1. Test Fixtures (`internal/testutil/fixtures.go`)

Add factory functions returning fully populated structs.

### Sample function

```go
func Sample<Type>() api.<Type> {
    return api.<Type>{
        ID:         1001,
        UserID:     12345,
        CreatedAt:  baseTime,
        UpdatedAt:  baseTime.Add(1 * time.Hour),
        Start:      baseTime,
        ScoreState: "SCORED",
        Score: &api.<Type>Score{
            // fill all fields with realistic values
        },
    }
}
```

Use `baseTime` (`2024-06-01T08:00:00Z`) as the anchor for all dates.

### Paginated helper

```go
func Paginated<Types>(n int, nextToken string) api.PaginatedResponse[api.<Type>] {
    records := make([]api.<Type>, n)
    for i := range records {
        item := Sample<Type>()
        item.ID = 1001 + i                    // int IDs
        // item.ID = fmt.Sprintf("<type>-%d-uuid", 1001+i)  // string IDs
        records[i] = item
    }
    return api.PaginatedResponse[api.<Type>]{
        Records:   records,
        NextToken: nextToken,
    }
}
```

## 2. API Service Tests (`internal/api/<name>_test.go`)

Use `testutil.NewAPIServer(t)` as the mock HTTP backend.

```go
package api

import (
    "context"
    "net/http"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/domo84/whoop-cli/internal/testutil"
)

func Test<Type>Service_List(t *testing.T) {
    expected := testutil.Paginated<Types>(2, "")
    srv := testutil.NewAPIServer(t)
    defer srv.Close()
    srv.Handle("GET", "/v2/<path>", http.StatusOK, expected)

    client := NewClientWithBaseURL(srv.URL, http.DefaultClient)
    svc := New<Type>Service(client)

    result, err := svc.List(context.Background(), ListParams{Limit: 25})
    require.NoError(t, err)
    assert.Len(t, result.Records, 2)
    srv.AssertCalled(t, "GET", "/v2/<path>", 1)
}

func Test<Type>Service_Get(t *testing.T) {
    expected := testutil.Sample<Type>()
    srv := testutil.NewAPIServer(t)
    defer srv.Close()
    srv.Handle("GET", "/v2/<path>/1001", http.StatusOK, expected)

    client := NewClientWithBaseURL(srv.URL, http.DefaultClient)
    svc := New<Type>Service(client)

    result, err := svc.Get(context.Background(), 1001)
    require.NoError(t, err)
    assert.Equal(t, 1001, result.ID)
}

func Test<Type>Service_Get_NotFound(t *testing.T) {
    srv := testutil.NewAPIServer(t)
    defer srv.Close()
    srv.Handle("GET", "/v2/<path>/9999", http.StatusNotFound, nil)

    client := NewClientWithBaseURL(srv.URL, http.DefaultClient)
    svc := New<Type>Service(client)

    _, err := svc.Get(context.Background(), 9999)
    require.Error(t, err)

    var apiErr *APIError
    require.ErrorAs(t, err, &apiErr)
    assert.Equal(t, 404, apiErr.StatusCode)
}
```

### APIServer methods

- `srv.Handle(method, path, status, body)` -- register canned JSON response
- `srv.HandleFunc(method, path, handler)` -- register custom handler
- `srv.HandleAny(path, handler)` -- any method
- `srv.CallCount(method, path)` -- get request count
- `srv.AssertCalled(t, method, path, times)` -- assert call count

## 3. Command Tests (`commands/commands_test.go`)

Add mock service and test functions to the existing test file.

### Mock service

```go
type mock<Type>Service struct {
    listResult    *api.PaginatedResponse[api.<Type>]
    listAllResult []api.<Type>
    getResult     *api.<Type>
    err           error
}

func (m *mock<Type>Service) List(_ context.Context, _ api.ListParams) (*api.PaginatedResponse[api.<Type>], error) {
    return m.listResult, m.err
}

func (m *mock<Type>Service) ListAll(_ context.Context, _ api.ListParams) ([]api.<Type>, error) {
    return m.listAllResult, m.err
}

func (m *mock<Type>Service) Get(_ context.Context, _ int) (*api.<Type>, error) {
    return m.getResult, m.err
}
```

### Command tests

```go
func Test<Type>List_Command(t *testing.T) {
    page := testutil.Paginated<Types>(2, "")
    setupState(t, &services{
        <Name>: &mock<Type>Service{listResult: &page},
    })

    out, err := runCmd(t, []string{"<name>", "list"})
    require.NoError(t, err)
    assert.NotEmpty(t, out)
}

func Test<Type>Get_Command(t *testing.T) {
    item := testutil.Sample<Type>()
    setupState(t, &services{
        <Name>: &mock<Type>Service{getResult: &item},
    })

    out, err := runCmd(t, []string{"<name>", "get", "1001"})
    require.NoError(t, err)
    assert.Contains(t, out, "SCORED")
}

func Test<Type>Get_InvalidID(t *testing.T) {
    setupState(t, &services{
        <Name>: &mock<Type>Service{},
    })

    _, err := runCmd(t, []string{"<name>", "get", "notanumber"})
    require.Error(t, err)
    assert.Contains(t, err.Error(), "invalid <name> ID")
}
```

### Key helpers

- `setupState(t, svc)` -- injects mock services, sets JSON formatter, auto-cleans up
- `runCmd(t, args)` -- runs root command with args, returns captured stdout

## Running Tests

```bash
go test ./...                          # all tests
go test ./commands/ -run Test<Type>    # specific command tests
go test ./internal/api/ -v             # verbose API tests
go test -race ./...                    # with race detector
```
