package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/domo84/whoop-cli/internal/api"
	"github.com/domo84/whoop-cli/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_APIError_NonOK(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()
	srv.Handle("GET", "/v2/user/profile/basic", http.StatusUnauthorized, map[string]string{"error": "unauthorized"})

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewUserService(client)

	_, err := svc.GetProfile(context.Background())
	require.Error(t, err)

	var apiErr *api.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
}

func TestClient_RateLimitRetry(t *testing.T) {
	callCount := 0
	profile := testutil.SampleUserProfile()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if callCount == 1 {
			// Return 429 with a reset time slightly in the past (immediate retry).
			resetTime := time.Now().Add(-1 * time.Millisecond).Unix()
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(profile)
	}))
	defer srv.Close()

	client := api.NewClientWithBaseURL(srv.URL, srv.Client())
	svc := api.NewUserService(client)

	result, err := svc.GetProfile(context.Background())
	require.NoError(t, err)
	assert.Equal(t, profile.Email, result.Email)
	assert.Equal(t, 2, callCount, "should have retried once after 429")
}


func TestBuildQuery(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	p := api.ListParams{
		Limit:     25,
		NextToken: "tok123",
		Start:     &start,
		End:       &end,
	}

	q := api.BuildQueryForTest(p)
	assert.Equal(t, "25", q.Get("limit"))
	assert.Equal(t, "tok123", q.Get("nextToken"))
	assert.Equal(t, "2024-01-01T00:00:00Z", q.Get("start"))
	assert.Equal(t, "2024-12-31T23:59:59Z", q.Get("end"))
}
