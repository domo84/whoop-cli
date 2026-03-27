package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/domo84/whoop-cli/internal/api"
	"github.com/domo84/whoop-cli/internal/output"
	"github.com/domo84/whoop-cli/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONFormatter_UserProfile(t *testing.T) {
	f := output.New("json")
	var buf bytes.Buffer
	profile := testutil.SampleUserProfile()
	require.NoError(t, f.Format(&buf, &profile))

	var decoded api.UserProfile
	require.NoError(t, json.NewDecoder(&buf).Decode(&decoded))
	assert.Equal(t, profile.Email, decoded.Email)
	assert.Equal(t, profile.UserID, decoded.UserID)
}

func TestJSONFormatter_Cycles(t *testing.T) {
	f := output.New("json")
	var buf bytes.Buffer
	cycles := []api.Cycle{testutil.SampleCycle()}
	require.NoError(t, f.Format(&buf, cycles))

	out := buf.String()
	assert.Contains(t, out, `"score_state"`)
	assert.Contains(t, out, `"strain"`)
}

func TestTableFormatter_UserProfile(t *testing.T) {
	f := output.New("table")
	var buf bytes.Buffer
	profile := testutil.SampleUserProfile()
	require.NoError(t, f.Format(&buf, &profile))

	out := buf.String()
	assert.Contains(t, out, "User ID")
	assert.Contains(t, out, "Email")
	assert.Contains(t, strings.ToLower(out), strings.ToLower(profile.Email))
}

func TestTableFormatter_BodyMeasurement(t *testing.T) {
	f := output.New("table")
	var buf bytes.Buffer
	m := testutil.SampleBodyMeasurement()
	require.NoError(t, f.Format(&buf, &m))

	out := buf.String()
	assert.Contains(t, out, "Height")
	assert.Contains(t, out, "Weight")
	assert.Contains(t, out, "Max Heart Rate")
}

func TestTableFormatter_Cycles(t *testing.T) {
	f := output.New("table")
	var buf bytes.Buffer
	cycles := []api.Cycle{testutil.SampleCycle()}
	require.NoError(t, f.Format(&buf, cycles))

	out := strings.ToUpper(buf.String())
	assert.Contains(t, out, "STRAIN")
	assert.Contains(t, out, "AVG HR")
	assert.Contains(t, out, "STATE")
}

func TestTableFormatter_PaginatedCycles(t *testing.T) {
	f := output.New("table")
	var buf bytes.Buffer
	page := testutil.PaginatedCycles(2, "")
	require.NoError(t, f.Format(&buf, &page))

	out := buf.String()
	assert.NotEmpty(t, out)
}

func TestTableFormatter_Sleeps(t *testing.T) {
	f := output.New("table")
	var buf bytes.Buffer
	sleeps := []api.Sleep{testutil.SampleSleep()}
	require.NoError(t, f.Format(&buf, sleeps))

	out := strings.ToUpper(buf.String())
	assert.Contains(t, out, "PERF")
	assert.Contains(t, out, "RESP RATE")
}

func TestTableFormatter_Recoveries(t *testing.T) {
	f := output.New("table")
	var buf bytes.Buffer
	recoveries := []api.Recovery{testutil.SampleRecovery()}
	require.NoError(t, f.Format(&buf, recoveries))

	out := strings.ToUpper(buf.String())
	assert.Contains(t, out, "RECOVERY")
	assert.Contains(t, out, "HRV")
}

func TestTableFormatter_Workouts(t *testing.T) {
	f := output.New("table")
	var buf bytes.Buffer
	workouts := []api.Workout{testutil.SampleWorkout()}
	require.NoError(t, f.Format(&buf, workouts))

	out := strings.ToUpper(buf.String())
	assert.Contains(t, out, "STRAIN")
	assert.Contains(t, out, "DURATION")
}

func TestTableFormatter_UnknownType_FallsBackToJSON(t *testing.T) {
	f := output.New("table")
	var buf bytes.Buffer
	data := map[string]string{"key": "value"}
	require.NoError(t, f.Format(&buf, data))
	assert.Contains(t, buf.String(), `"key"`)
}

func TestNew_DefaultsToTable(t *testing.T) {
	f := output.New("unknown-format")
	assert.IsType(t, &output.TableFormatter{}, f)
}

func TestNew_JSON(t *testing.T) {
	f := output.New("json")
	assert.IsType(t, &output.JSONFormatter{}, f)
}
