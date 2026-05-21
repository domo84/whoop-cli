package commands

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/domo84/whoop-cli/internal/api"
	"github.com/domo84/whoop-cli/internal/config"
	"github.com/domo84/whoop-cli/internal/output"
	"github.com/domo84/whoop-cli/internal/testutil"
)

// --- Mock services ---

type mockUserService struct {
	profile     *api.UserProfile
	measurement *api.BodyMeasurement
	err         error
}

func (m *mockUserService) GetProfile(_ context.Context) (*api.UserProfile, error) {
	return m.profile, m.err
}

func (m *mockUserService) GetBodyMeasurement(_ context.Context) (*api.BodyMeasurement, error) {
	return m.measurement, m.err
}

type mockCycleService struct {
	listResult    *api.PaginatedResponse[api.Cycle]
	listAllResult []api.Cycle
	getResult     *api.Cycle
	sleepResult   *api.Sleep
	err           error
}

func (m *mockCycleService) List(_ context.Context, _ api.ListParams) (*api.PaginatedResponse[api.Cycle], error) {
	return m.listResult, m.err
}

func (m *mockCycleService) ListAll(_ context.Context, _ api.ListParams) ([]api.Cycle, error) {
	return m.listAllResult, m.err
}

func (m *mockCycleService) Get(_ context.Context, _ int) (*api.Cycle, error) {
	return m.getResult, m.err
}

func (m *mockCycleService) GetSleep(_ context.Context, _ int) (*api.Sleep, error) {
	return m.sleepResult, m.err
}

type mockSleepService struct {
	listResult    *api.PaginatedResponse[api.Sleep]
	listAllResult []api.Sleep
	getResult     *api.Sleep
	err           error
}

func (m *mockSleepService) List(_ context.Context, _ api.ListParams) (*api.PaginatedResponse[api.Sleep], error) {
	return m.listResult, m.err
}

func (m *mockSleepService) ListAll(_ context.Context, _ api.ListParams) ([]api.Sleep, error) {
	return m.listAllResult, m.err
}

func (m *mockSleepService) Get(_ context.Context, _ string) (*api.Sleep, error) {
	return m.getResult, m.err
}

type mockRecoveryService struct {
	listResult    *api.PaginatedResponse[api.Recovery]
	listAllResult []api.Recovery
	getResult     *api.Recovery
	err           error
}

func (m *mockRecoveryService) List(_ context.Context, _ api.ListParams) (*api.PaginatedResponse[api.Recovery], error) {
	return m.listResult, m.err
}

func (m *mockRecoveryService) ListAll(_ context.Context, _ api.ListParams) ([]api.Recovery, error) {
	return m.listAllResult, m.err
}

func (m *mockRecoveryService) Get(_ context.Context, _ int) (*api.Recovery, error) {
	return m.getResult, m.err
}

type mockWorkoutService struct {
	listResult    *api.PaginatedResponse[api.Workout]
	listAllResult []api.Workout
	getResult     *api.Workout
	err           error
}

func (m *mockWorkoutService) List(_ context.Context, _ api.ListParams) (*api.PaginatedResponse[api.Workout], error) {
	return m.listResult, m.err
}

func (m *mockWorkoutService) ListAll(_ context.Context, _ api.ListParams) ([]api.Workout, error) {
	return m.listAllResult, m.err
}

func (m *mockWorkoutService) Get(_ context.Context, _ string) (*api.Workout, error) {
	return m.getResult, m.err
}

// setupState injects mock state before running a command and restores it after.
func setupState(t *testing.T, svc *services) {
	t.Helper()
	dir := t.TempDir()
	mgr, err := config.NewManager(dir)
	require.NoError(t, err)
	state.cfgMgr = mgr
	state.cfg = &config.Config{OutputFormat: "json", RedirectURI: config.DefaultRedirectURI}
	state.svc = svc
	state.formatter = output.New("json")
	t.Cleanup(func() {
		state.svc = nil
		state.cfg = nil
		state.cfgMgr = nil
		state.formatter = nil
	})
}

// runCmd executes the root command with args and captures stdout.
func runCmd(t *testing.T, args []string) (string, error) {
	t.Helper()
	// Disable color for deterministic output in tests.
	os.Setenv("NO_COLOR", "1")
	root := NewRootCmd()
	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&bytes.Buffer{})
	root.SetArgs(args)
	err := root.Execute()
	return out.String(), err
}

// --- User command tests ---

func TestUserProfile_Command(t *testing.T) {
	profile := testutil.SampleUserProfile()
	setupState(t, &services{
		User: &mockUserService{profile: &profile},
	})

	out, err := runCmd(t, []string{"user", "profile"})
	require.NoError(t, err)
	assert.Contains(t, out, "test@example.com")
}

func TestUserMeasurements_Command(t *testing.T) {
	m := testutil.SampleBodyMeasurement()
	setupState(t, &services{
		User: &mockUserService{measurement: &m},
	})

	out, err := runCmd(t, []string{"user", "measurements"})
	require.NoError(t, err)
	assert.Contains(t, out, "height_meter")
}

// --- Cycle command tests ---

func TestCycleList_Command(t *testing.T) {
	page := testutil.PaginatedCycles(2, "")
	setupState(t, &services{
		Cycle: &mockCycleService{listResult: &page},
	})

	out, err := runCmd(t, []string{"cycle", "list", "--limit", "2"})
	require.NoError(t, err)
	assert.Contains(t, out, "SCORED")
}

func TestCycleList_AllFlag(t *testing.T) {
	records := []api.Cycle{testutil.SampleCycle(), testutil.SampleCycle()}
	setupState(t, &services{
		Cycle: &mockCycleService{listAllResult: records},
	})

	out, err := runCmd(t, []string{"cycle", "list", "--all"})
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestCycleGet_Command(t *testing.T) {
	c := testutil.SampleCycle()
	setupState(t, &services{
		Cycle: &mockCycleService{getResult: &c},
	})

	out, err := runCmd(t, []string{"cycle", "get", "1001"})
	require.NoError(t, err)
	assert.Contains(t, out, "SCORED")
}

func TestCycleGet_InvalidID(t *testing.T) {
	setupState(t, &services{
		Cycle: &mockCycleService{},
	})

	_, err := runCmd(t, []string{"cycle", "get", "notanumber"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid cycle ID")
}

func TestCycleSleep_Command(t *testing.T) {
	s := testutil.SampleSleep()
	setupState(t, &services{
		Cycle: &mockCycleService{sleepResult: &s},
	})

	out, err := runCmd(t, []string{"cycle", "sleep", "1001"})
	require.NoError(t, err)
	assert.Contains(t, out, "SCORED")
}

// --- Sleep command tests ---

func TestSleepList_Command(t *testing.T) {
	page := testutil.PaginatedSleeps(2, "")
	setupState(t, &services{
		Sleep: &mockSleepService{listResult: &page},
	})

	out, err := runCmd(t, []string{"sleep", "list"})
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestSleepGet_Command(t *testing.T) {
	s := testutil.SampleSleep()
	setupState(t, &services{
		Sleep: &mockSleepService{getResult: &s},
	})

	out, err := runCmd(t, []string{"sleep", "get", "2001"})
	require.NoError(t, err)
	assert.Contains(t, out, "SCORED")
}

// --- Recovery command tests ---

func TestRecoveryList_Command(t *testing.T) {
	page := testutil.PaginatedRecoveries(2, "")
	setupState(t, &services{
		Recovery: &mockRecoveryService{listResult: &page},
	})

	out, err := runCmd(t, []string{"recovery", "list"})
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestRecoveryGet_Command(t *testing.T) {
	r := testutil.SampleRecovery()
	setupState(t, &services{
		Recovery: &mockRecoveryService{getResult: &r},
	})

	out, err := runCmd(t, []string{"recovery", "get", "1001"})
	require.NoError(t, err)
	assert.Contains(t, out, "SCORED")
}

// --- Workout command tests ---

func TestWorkoutList_Command(t *testing.T) {
	page := testutil.PaginatedWorkouts(3, "")
	setupState(t, &services{
		Workout: &mockWorkoutService{listResult: &page},
	})

	out, err := runCmd(t, []string{"workout", "list"})
	require.NoError(t, err)
	assert.NotEmpty(t, out)
}

func TestWorkoutGet_Command(t *testing.T) {
	w := testutil.SampleWorkout()
	setupState(t, &services{
		Workout: &mockWorkoutService{getResult: &w},
	})

	out, err := runCmd(t, []string{"workout", "get", "3001"})
	require.NoError(t, err)
	assert.Contains(t, out, "SCORED")
}

