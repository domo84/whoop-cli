package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/domo84/whoop-cli/internal/api"
	"github.com/domo84/whoop-cli/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetProfile(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()

	expected := testutil.SampleUserProfile()
	srv.Handle("GET", "/v2/user/profile/basic", http.StatusOK, expected)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewUserService(client)

	result, err := svc.GetProfile(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expected.UserID, result.UserID)
	assert.Equal(t, expected.Email, result.Email)
	assert.Equal(t, expected.FirstName, result.FirstName)
	assert.Equal(t, expected.LastName, result.LastName)
}

func TestUserService_GetProfile_Error(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()
	srv.Handle("GET", "/v2/user/profile/basic", http.StatusUnauthorized, nil)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewUserService(client)

	_, err := svc.GetProfile(context.Background())
	require.Error(t, err)
	var apiErr *api.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
}

func TestUserService_GetBodyMeasurement(t *testing.T) {
	srv := testutil.NewAPIServer(t)
	defer srv.Close()

	expected := testutil.SampleBodyMeasurement()
	srv.Handle("GET", "/v2/user/measurement/body", http.StatusOK, expected)

	client := api.NewClientWithBaseURL(srv.URL, http.DefaultClient)
	svc := api.NewUserService(client)

	result, err := svc.GetBodyMeasurement(context.Background())
	require.NoError(t, err)
	assert.InDelta(t, expected.HeightMeter, result.HeightMeter, 0.001)
	assert.InDelta(t, expected.WeightKilogram, result.WeightKilogram, 0.001)
	assert.Equal(t, expected.MaxHeartRate, result.MaxHeartRate)
}
