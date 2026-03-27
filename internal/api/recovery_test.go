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

func TestRecoveryService_List(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()

	expected := testutil.PaginatedRecoveries(3, "")
	srv.Handle("GET", "/v2/recovery", http.StatusOK, expected)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewRecoveryService(client)

	result, err := svc.List(context.Background(), api.ListParams{Limit: 3})
	require.NoError(t, err)
	assert.Len(t, result.Records, 3)
	assert.Empty(t, result.NextToken)
}

func TestRecoveryService_ListAll_Pagination(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(testutil.PaginatedRecoveries(3, "page2"))
		} else {
			json.NewEncoder(w).Encode(testutil.PaginatedRecoveries(2, ""))
		}
	}))
	defer srv.Close()

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewRecoveryService(client)

	all, err := svc.ListAll(context.Background(), api.ListParams{})
	require.NoError(t, err)
	assert.Len(t, all, 5)
	assert.Equal(t, 2, callCount)
}

func TestRecoveryService_Get(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()

	expected := testutil.SampleRecovery()
	srv.Handle("GET", "/v2/cycle/1001/recovery", http.StatusOK, expected)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewRecoveryService(client)

	result, err := svc.Get(context.Background(), 1001)
	require.NoError(t, err)
	assert.Equal(t, "SCORED", result.ScoreState)
	assert.NotNil(t, result.Score)
	assert.InDelta(t, 78.5, result.Score.RecoveryScore, 0.01)
	assert.InDelta(t, 52.0, result.Score.RestingHeartRate, 0.01)
	assert.InDelta(t, 48.3, result.Score.HrvRmssdMilli, 0.01)
	assert.NotNil(t, result.Score.Spo2Percentage)
	assert.InDelta(t, 98.5, *result.Score.Spo2Percentage, 0.01)
}

func TestRecoveryService_Get_NotFound(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()
	srv.Handle("GET", "/v2/cycle/9999/recovery", http.StatusNotFound, nil)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewRecoveryService(client)

	_, err := svc.Get(context.Background(), 9999)
	require.Error(t, err)
	var apiErr *api.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
}
