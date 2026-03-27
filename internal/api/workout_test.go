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

func TestWorkoutService_List(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()

	expected := testutil.PaginatedWorkouts(3, "")
	srv.Handle("GET", "/v2/activity/workout", http.StatusOK, expected)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewWorkoutService(client)

	result, err := svc.List(context.Background(), api.ListParams{Limit: 3})
	require.NoError(t, err)
	assert.Len(t, result.Records, 3)
	assert.Empty(t, result.NextToken)
}

func TestWorkoutService_ListAll_Pagination(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		if callCount == 1 {
			json.NewEncoder(w).Encode(testutil.PaginatedWorkouts(5, "page2"))
		} else {
			json.NewEncoder(w).Encode(testutil.PaginatedWorkouts(2, ""))
		}
	}))
	defer srv.Close()

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewWorkoutService(client)

	all, err := svc.ListAll(context.Background(), api.ListParams{})
	require.NoError(t, err)
	assert.Len(t, all, 7)
	assert.Equal(t, 2, callCount)
}

func TestWorkoutService_Get(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()

	expected := testutil.SampleWorkout()
	srv.Handle("GET", "/v2/activity/workout/workout-3001-uuid", http.StatusOK, expected)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewWorkoutService(client)

	result, err := svc.Get(context.Background(), "workout-3001-uuid")
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, "SCORED", result.ScoreState)
	assert.NotNil(t, result.Score)
	assert.InDelta(t, 9.1, result.Score.Strain, 0.01)
	assert.Equal(t, 155, result.Score.AverageHeartRate)
	assert.NotNil(t, result.Score.DistanceMeter)
	assert.InDelta(t, 8000.0, *result.Score.DistanceMeter, 0.01)
}

func TestWorkoutService_Get_NotFound(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()
	srv.Handle("GET", "/v2/activity/workout/nonexistent-uuid", http.StatusNotFound, nil)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewWorkoutService(client)

	_, err := svc.Get(context.Background(), "nonexistent-uuid")
	require.Error(t, err)
	var apiErr *api.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
}
