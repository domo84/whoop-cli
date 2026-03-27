package api_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/domo84/whoop-cli/internal/api"
	"github.com/domo84/whoop-cli/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSleepService_List(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()

	expected := testutil.PaginatedSleeps(2, "")
	srv.Handle("GET", "/v2/activity/sleep", http.StatusOK, expected)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewSleepService(client)

	result, err := svc.List(context.Background(), api.ListParams{Limit: 2})
	require.NoError(t, err)
	assert.Len(t, result.Records, 2)
	assert.Empty(t, result.NextToken)
}

func TestSleepService_ListAll_Pagination(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(testutil.PaginatedSleeps(4, "next-page"))
		} else {
			json.NewEncoder(w).Encode(testutil.PaginatedSleeps(2, ""))
		}
	}))
	defer srv.Close()

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewSleepService(client)

	all, err := svc.ListAll(context.Background(), api.ListParams{})
	require.NoError(t, err)
	assert.Len(t, all, 6)
	assert.Equal(t, 2, callCount)
}

func TestSleepService_Get(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()

	expected := testutil.SampleSleep()
	srv.Handle("GET", "/v2/activity/sleep/sleep-2001-uuid", http.StatusOK, expected)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewSleepService(client)

	result, err := svc.Get(context.Background(), "sleep-2001-uuid")
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.False(t, result.Nap)
	assert.NotNil(t, result.Score)
	assert.InDelta(t, 14.2, result.Score.RespiratoryRate, 0.01)
}

func TestSleepService_Get_NotFound(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()
	srv.Handle("GET", "/v2/activity/sleep/nonexistent-uuid", http.StatusNotFound, nil)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewSleepService(client)

	_, err := svc.Get(context.Background(), "nonexistent-uuid")
	require.Error(t, err)
	var apiErr *api.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
}
