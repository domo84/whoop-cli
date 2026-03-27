package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/magnusmv/whoop-cli/internal/api"
	"github.com/magnusmv/whoop-cli/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCycleService_List(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()

	expected := testutil.PaginatedCycles(3, "")
	srv.Handle("GET", "/v2/cycle", http.StatusOK, expected)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewCycleService(client)

	result, err := svc.List(context.Background(), api.ListParams{Limit: 3})
	require.NoError(t, err)
	assert.Len(t, result.Records, 3)
	assert.Empty(t, result.NextToken)
	srv.AssertCalled(t, "GET", "/v2/cycle", 1)
}

func TestCycleService_ListAll_Pagination(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(testutil.PaginatedCycles(5, "token-abc"))
		} else {
			json.NewEncoder(w).Encode(testutil.PaginatedCycles(3, ""))
		}
	}))
	defer srv.Close()

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewCycleService(client)

	all, err := svc.ListAll(context.Background(), api.ListParams{})
	require.NoError(t, err)
	assert.Len(t, all, 8)
	assert.Equal(t, 2, callCount, "should have made exactly 2 requests")
}

func TestCycleService_Get(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()

	expected := testutil.SampleCycle()
	srv.Handle("GET", "/v2/cycle/1001", http.StatusOK, expected)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewCycleService(client)

	result, err := svc.Get(context.Background(), 1001)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, "SCORED", result.ScoreState)
	assert.NotNil(t, result.Score)
}

func TestCycleService_Get_NotFound(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()
	srv.Handle("GET", "/v2/cycle/9999", http.StatusNotFound, nil)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewCycleService(client)

	_, err := svc.Get(context.Background(), 9999)
	require.Error(t, err)
	var apiErr *api.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
}

func TestCycleService_GetSleep(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()

	expected := testutil.SampleSleep()
	srv.Handle("GET", "/v2/cycle/1001/sleep", http.StatusOK, expected)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewCycleService(client)

	result, err := svc.GetSleep(context.Background(), 1001)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, "SCORED", result.ScoreState)
}

func TestCycleService_List_QueryParams(t *testing.T) {
	var capturedQuery string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(testutil.PaginatedCycles(1, ""))
	}))
	defer srv.Close()

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewCycleService(client)

	_, err := svc.List(context.Background(), api.ListParams{Limit: 10, NextToken: "tok"})
	require.NoError(t, err)
	assert.Contains(t, capturedQuery, "limit=10")
	assert.Contains(t, capturedQuery, "nextToken=tok")
}
