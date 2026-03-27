package api

import "time"

// PaginatedResponse is a generic wrapper for all paginated API responses.
type PaginatedResponse[T any] struct {
	Records   []T    `json:"records"`
	NextToken string `json:"next_token,omitempty"`
}

// ListParams contains common query parameters for list endpoints.
type ListParams struct {
	Limit     int
	NextToken string
	Start     *time.Time
	End       *time.Time
}

// --- User ---

// UserProfile contains basic profile information for the authenticated user.
type UserProfile struct {
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// BodyMeasurement contains physical measurements for the authenticated user.
type BodyMeasurement struct {
	HeightMeter    float64 `json:"height_meter"`
	WeightKilogram float64 `json:"weight_kilogram"`
	MaxHeartRate   int     `json:"max_heart_rate"`
}

// --- Cycle ---

// Cycle represents a WHOOP physiological cycle.
type Cycle struct {
	ID             int        `json:"id"`
	UserID         int        `json:"user_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Start          time.Time  `json:"start"`
	End            *time.Time `json:"end,omitempty"`
	TimezoneOffset string     `json:"timezone_offset"`
	ScoreState     string     `json:"score_state"`
	Score          *CycleScore `json:"score,omitempty"`
}

// CycleScore holds metrics for a cycle.
type CycleScore struct {
	Strain           float64 `json:"strain"`
	Kilojoule        float64 `json:"kilojoule"`
	AverageHeartRate int     `json:"average_heart_rate"`
	MaxHeartRate     int     `json:"max_heart_rate"`
}

// --- Sleep ---

// Sleep represents a WHOOP sleep or nap activity.
type Sleep struct {
	ID             string     `json:"id"`
	UserID         int        `json:"user_id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Start          time.Time  `json:"start"`
	End            time.Time  `json:"end"`
	TimezoneOffset string     `json:"timezone_offset"`
	Nap            bool       `json:"nap"`
	ScoreState     string     `json:"score_state"`
	Score          *SleepScore `json:"score,omitempty"`
}

// SleepScore holds metrics for a sleep activity.
type SleepScore struct {
	StageSummary     SleepStageSummary `json:"stage_summary"`
	SleepNeeded      SleepNeeded       `json:"sleep_needed"`
	RespiratoryRate  float64           `json:"respiratory_rate"`
	SleepPerformance *float64          `json:"sleep_performance_percentage,omitempty"`
	SleepConsistency *float64          `json:"sleep_consistency_percentage,omitempty"`
	SleepEfficiency  *float64          `json:"sleep_efficiency_percentage,omitempty"`
}

// SleepStageSummary holds durations for each sleep stage.
type SleepStageSummary struct {
	TotalInBedTimeMilli         int `json:"total_in_bed_time_milli"`
	TotalAwakeTimeMilli         int `json:"total_awake_time_milli"`
	TotalNoDataTimeMilli        int `json:"total_no_data_time_milli"`
	TotalLightSleepTimeMilli    int `json:"total_light_sleep_time_milli"`
	TotalSlowWaveSleepTimeMilli int `json:"total_slow_wave_sleep_time_milli"`
	TotalRemSleepTimeMilli      int `json:"total_rem_sleep_time_milli"`
	SleepCycleCount             int `json:"sleep_cycle_count"`
	DisturbanceCount            int `json:"disturbance_count"`
}

// SleepNeeded holds sleep need breakdown.
type SleepNeeded struct {
	BaselineMilli             int `json:"baseline_milli"`
	NeedFromSleepDebtMilli    int `json:"need_from_sleep_debt_milli"`
	NeedFromRecentStrainMilli int `json:"need_from_recent_strain_milli"`
	NeedFromRecentNapMilli    int `json:"need_from_recent_nap_milli"`
}

// --- Recovery ---

// Recovery represents a WHOOP recovery score.
type Recovery struct {
	CycleID    int           `json:"cycle_id"`
	SleepID    string        `json:"sleep_id"`
	UserID     int           `json:"user_id"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	ScoreState string        `json:"score_state"`
	Score      *RecoveryScore `json:"score,omitempty"`
}

// RecoveryScore holds recovery metrics.
type RecoveryScore struct {
	UserCalibrating  bool     `json:"user_calibrating"`
	RecoveryScore    float64  `json:"recovery_score"`
	RestingHeartRate float64  `json:"resting_heart_rate"`
	HrvRmssdMilli    float64  `json:"hrv_rmssd_milli"`
	Spo2Percentage   *float64 `json:"spo2_percentage,omitempty"`
	SkinTempCelsius  *float64 `json:"skin_temp_celsius,omitempty"`
}

// --- Workout ---

// Workout represents a WHOOP workout activity.
type Workout struct {
	ID             string       `json:"id"`
	UserID         int          `json:"user_id"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Start          time.Time    `json:"start"`
	End            time.Time    `json:"end"`
	TimezoneOffset string       `json:"timezone_offset"`
	SportID        int          `json:"sport_id"`
	ScoreState     string       `json:"score_state"`
	Score          *WorkoutScore `json:"score,omitempty"`
}

// WorkoutScore holds workout metrics.
type WorkoutScore struct {
	Strain              float64      `json:"strain"`
	AverageHeartRate    int          `json:"average_heart_rate"`
	MaxHeartRate        int          `json:"max_heart_rate"`
	Kilojoule           float64      `json:"kilojoule"`
	PercentRecorded     float64      `json:"percent_recorded"`
	DistanceMeter       *float64     `json:"distance_meter,omitempty"`
	AltitudeGainMeter   *float64     `json:"altitude_gain_meter,omitempty"`
	AltitudeChangeMeter *float64     `json:"altitude_change_meter,omitempty"`
	ZoneDuration        ZoneDuration `json:"zone_duration"`
}

// ZoneDuration holds time spent in each heart rate zone.
type ZoneDuration struct {
	ZoneZeroMilli  int `json:"zone_zero_milli"`
	ZoneOneMilli   int `json:"zone_one_milli"`
	ZoneTwoMilli   int `json:"zone_two_milli"`
	ZoneThreeMilli int `json:"zone_three_milli"`
	ZoneFourMilli  int `json:"zone_four_milli"`
	ZoneFiveMilli  int `json:"zone_five_milli"`
}

// --- Errors ---

// APIError represents a non-2xx response from the Whoop API.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return "whoop API error " + itoa(e.StatusCode) + ": " + e.Message
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [10]byte
	pos := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		pos--
		buf[pos] = byte(n%10) + '0'
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
