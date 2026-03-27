package commands

import (
	"time"
)

// parseDate tries to parse s as RFC3339, then as YYYY-MM-DD (midnight UTC).
// Returns nil if neither format matches.
func parseDate(s string) *time.Time {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return &t
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		utc := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
		return &utc
	}
	return nil
}
