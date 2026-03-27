package testutil

import (
	"fmt"
	"time"

	"github.com/magnusmv/whoop-cli/internal/api"
)

var baseTime = time.Date(2024, 6, 1, 8, 0, 0, 0, time.UTC)

func SampleUserProfile() api.UserProfile {
	return api.UserProfile{
		UserID:    12345,
		Email:     "test@example.com",
		FirstName: "Jane",
		LastName:  "Doe",
	}
}

func SampleBodyMeasurement() api.BodyMeasurement {
	return api.BodyMeasurement{
		HeightMeter:    1.75,
		WeightKilogram: 70.5,
		MaxHeartRate:   195,
	}
}

func SampleCycle() api.Cycle {
	end := baseTime.Add(24 * time.Hour)
	strain := 8.5
	kilojoule := 2500.0
	avgHR := 72
	maxHR := 145
	return api.Cycle{
		ID:             1001,
		UserID:         12345,
		CreatedAt:      baseTime,
		UpdatedAt:      baseTime.Add(1 * time.Hour),
		Start:          baseTime,
		End:            &end,
		TimezoneOffset: "+00:00",
		ScoreState:     "SCORED",
		Score: &api.CycleScore{
			Strain:           strain,
			Kilojoule:        kilojoule,
			AverageHeartRate: avgHR,
			MaxHeartRate:     maxHR,
		},
	}
}

func SampleSleep() api.Sleep {
	perf := 88.0
	cons := 92.0
	eff := 90.0
	return api.Sleep{
		ID:             "sleep-2001-uuid",
		UserID:         12345,
		CreatedAt:      baseTime,
		UpdatedAt:      baseTime.Add(30 * time.Minute),
		Start:          baseTime.Add(-8 * time.Hour),
		End:            baseTime,
		TimezoneOffset: "+00:00",
		Nap:            false,
		ScoreState:     "SCORED",
		Score: &api.SleepScore{
			RespiratoryRate:  14.2,
			SleepPerformance: &perf,
			SleepConsistency: &cons,
			SleepEfficiency:  &eff,
			StageSummary: api.SleepStageSummary{
				TotalInBedTimeMilli:         28800000,
				TotalAwakeTimeMilli:         1800000,
				TotalLightSleepTimeMilli:    9000000,
				TotalSlowWaveSleepTimeMilli: 8000000,
				TotalRemSleepTimeMilli:      7000000,
				SleepCycleCount:             5,
				DisturbanceCount:            2,
			},
			SleepNeeded: api.SleepNeeded{
				BaselineMilli:             28800000,
				NeedFromSleepDebtMilli:    1800000,
				NeedFromRecentStrainMilli: 900000,
				NeedFromRecentNapMilli:    0,
			},
		},
	}
}

func SampleRecovery() api.Recovery {
	spo2 := 98.5
	skinTemp := 34.2
	return api.Recovery{
		CycleID:    1001,
		SleepID:    "sleep-2001-uuid",
		UserID:     12345,
		CreatedAt:  baseTime.Add(30 * time.Minute),
		UpdatedAt:  baseTime.Add(1 * time.Hour),
		ScoreState: "SCORED",
		Score: &api.RecoveryScore{
			UserCalibrating:  false,
			RecoveryScore:    78.5,
			RestingHeartRate: 52.0,
			HrvRmssdMilli:    48.3,
			Spo2Percentage:   &spo2,
			SkinTempCelsius:  &skinTemp,
		},
	}
}

func SampleWorkout() api.Workout {
	distM := 8000.0
	return api.Workout{
		ID:             "workout-3001-uuid",
		UserID:         12345,
		CreatedAt:      baseTime.Add(2 * time.Hour),
		UpdatedAt:      baseTime.Add(3 * time.Hour),
		Start:          baseTime.Add(2 * time.Hour),
		End:            baseTime.Add(3 * time.Hour),
		TimezoneOffset: "+00:00",
		SportID:        0,
		ScoreState:     "SCORED",
		Score: &api.WorkoutScore{
			Strain:           9.1,
			AverageHeartRate: 155,
			MaxHeartRate:     182,
			Kilojoule:        1800.0,
			PercentRecorded:  100.0,
			DistanceMeter:    &distM,
			ZoneDuration: api.ZoneDuration{
				ZoneZeroMilli:  300000,
				ZoneOneMilli:   600000,
				ZoneTwoMilli:   1200000,
				ZoneThreeMilli: 900000,
				ZoneFourMilli:  600000,
				ZoneFiveMilli:  300000,
			},
		},
	}
}

// PaginatedCycles returns a paginated response containing n copies of SampleCycle.
func PaginatedCycles(n int, nextToken string) api.PaginatedResponse[api.Cycle] {
	records := make([]api.Cycle, n)
	for i := range records {
		c := SampleCycle()
		c.ID = 1001 + i
		records[i] = c
	}
	return api.PaginatedResponse[api.Cycle]{
		Records:   records,
		NextToken: nextToken,
	}
}

// PaginatedSleeps returns a paginated response containing n copies of SampleSleep.
func PaginatedSleeps(n int, nextToken string) api.PaginatedResponse[api.Sleep] {
	records := make([]api.Sleep, n)
	for i := range records {
		s := SampleSleep()
		s.ID = fmt.Sprintf("sleep-%d-uuid", 2001+i)
		records[i] = s
	}
	return api.PaginatedResponse[api.Sleep]{
		Records:   records,
		NextToken: nextToken,
	}
}

// PaginatedRecoveries returns a paginated response containing n copies of SampleRecovery.
func PaginatedRecoveries(n int, nextToken string) api.PaginatedResponse[api.Recovery] {
	records := make([]api.Recovery, n)
	for i := range records {
		r := SampleRecovery()
		r.CycleID = 1001 + i
		records[i] = r
	}
	return api.PaginatedResponse[api.Recovery]{
		Records:   records,
		NextToken: nextToken,
	}
}

// PaginatedWorkouts returns a paginated response containing n copies of SampleWorkout.
func PaginatedWorkouts(n int, nextToken string) api.PaginatedResponse[api.Workout] {
	records := make([]api.Workout, n)
	for i := range records {
		w := SampleWorkout()
		w.ID = fmt.Sprintf("workout-%d-uuid", 3001+i)
		records[i] = w
	}
	return api.PaginatedResponse[api.Workout]{
		Records:   records,
		NextToken: nextToken,
	}
}
