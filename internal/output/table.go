package output

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"

	"github.com/magnusmv/whoop-cli/internal/api"
)

// TableFormatter renders data as an ASCII table.
type TableFormatter struct{}

func (f *TableFormatter) Format(w io.Writer, data any) error {
	switch v := data.(type) {
	case *api.UserProfile:
		return renderUserProfile(w, v)
	case *api.BodyMeasurement:
		return renderBodyMeasurement(w, v)

	case []api.Cycle:
		return renderCycles(w, v)
	case *api.Cycle:
		return renderCycles(w, []api.Cycle{*v})
	case *api.PaginatedResponse[api.Cycle]:
		return renderCycles(w, v.Records)

	case []api.Sleep:
		return renderSleeps(w, v)
	case *api.Sleep:
		return renderSleeps(w, []api.Sleep{*v})
	case *api.PaginatedResponse[api.Sleep]:
		return renderSleeps(w, v.Records)

	case []api.Recovery:
		return renderRecoveries(w, v)
	case *api.Recovery:
		return renderRecoveries(w, []api.Recovery{*v})
	case *api.PaginatedResponse[api.Recovery]:
		return renderRecoveries(w, v.Records)

	case []api.Workout:
		return renderWorkouts(w, v)
	case *api.Workout:
		return renderWorkouts(w, []api.Workout{*v})
	case *api.PaginatedResponse[api.Workout]:
		return renderWorkouts(w, v.Records)

	default:
		return (&JSONFormatter{}).Format(w, data)
	}
}

// noBorderTable creates a Table with no outer borders.
func noBorderTable(w io.Writer, headers []string) *tablewriter.Table {
	table := tablewriter.NewTable(w,
		tablewriter.WithBorders(tw.Border{
			Left:   tw.Off,
			Right:  tw.Off,
			Top:    tw.Off,
			Bottom: tw.Off,
		}),
		tablewriter.WithHeader(headers),
		tablewriter.WithHeaderAlignment(tw.AlignLeft),
		tablewriter.WithAlignment(tw.Alignment{tw.AlignLeft}),
	)
	return table
}

func renderUserProfile(w io.Writer, p *api.UserProfile) error {
	table := noBorderTable(w, []string{"FIELD", "VALUE"})
	table.Append([]string{"User ID", strconv.Itoa(p.UserID)})
	table.Append([]string{"Email", p.Email})
	table.Append([]string{"First Name", p.FirstName})
	table.Append([]string{"Last Name", p.LastName})
	return table.Render()
}

func renderBodyMeasurement(w io.Writer, m *api.BodyMeasurement) error {
	table := noBorderTable(w, []string{"FIELD", "VALUE"})
	table.Append([]string{"Height (m)", fmt.Sprintf("%.2f", m.HeightMeter)})
	table.Append([]string{"Weight (kg)", fmt.Sprintf("%.2f", m.WeightKilogram)})
	table.Append([]string{"Max Heart Rate", strconv.Itoa(m.MaxHeartRate) + " bpm"})
	return table.Render()
}

func renderCycles(w io.Writer, cycles []api.Cycle) error {
	table := noBorderTable(w, []string{"ID", "START", "END", "STRAIN", "AVG HR", "MAX HR", "STATE"})
	for _, c := range cycles {
		strain, avgHR, maxHR := "N/A", "N/A", "N/A"
		if c.Score != nil {
			strain = fmt.Sprintf("%.1f", c.Score.Strain)
			avgHR = strconv.Itoa(c.Score.AverageHeartRate)
			maxHR = strconv.Itoa(c.Score.MaxHeartRate)
		}
		end := "active"
		if c.End != nil {
			end = c.End.Format("2006-01-02 15:04")
		}
		table.Append([]string{
			strconv.Itoa(c.ID),
			c.Start.Format("2006-01-02 15:04"),
			end,
			strain,
			avgHR,
			maxHR,
			c.ScoreState,
		})
	}
	return table.Render()
}

func renderSleeps(w io.Writer, sleeps []api.Sleep) error {
	table := noBorderTable(w, []string{"ID", "START", "END", "NAP", "PERF%", "RESP RATE", "STATE"})
	for _, s := range sleeps {
		perf, respRate := "N/A", "N/A"
		if s.Score != nil {
			respRate = fmt.Sprintf("%.1f", s.Score.RespiratoryRate)
			if s.Score.SleepPerformance != nil {
				perf = fmt.Sprintf("%.0f%%", *s.Score.SleepPerformance)
			}
		}
		nap := "No"
		if s.Nap {
			nap = "Yes"
		}
		table.Append([]string{
			s.ID,
			s.Start.Format("2006-01-02 15:04"),
			s.End.Format("2006-01-02 15:04"),
			nap,
			perf,
			respRate,
			s.ScoreState,
		})
	}
	return table.Render()
}

func renderRecoveries(w io.Writer, recoveries []api.Recovery) error {
	table := noBorderTable(w, []string{"CYCLE ID", "RECOVERY%", "RHR", "HRV", "SPO2%", "SKIN TEMP", "STATE"})
	for _, r := range recoveries {
		rec, rhr, hrv, spo2, skinTemp := "N/A", "N/A", "N/A", "N/A", "N/A"
		if r.Score != nil {
			rec = fmt.Sprintf("%.0f%%", r.Score.RecoveryScore)
			rhr = fmt.Sprintf("%.0f", r.Score.RestingHeartRate)
			hrv = fmt.Sprintf("%.1f ms", r.Score.HrvRmssdMilli)
			if r.Score.Spo2Percentage != nil {
				spo2 = fmt.Sprintf("%.1f%%", *r.Score.Spo2Percentage)
			}
			if r.Score.SkinTempCelsius != nil {
				skinTemp = fmt.Sprintf("%.1f°C", *r.Score.SkinTempCelsius)
			}
		}
		table.Append([]string{
			strconv.Itoa(r.CycleID),
			rec,
			rhr,
			hrv,
			spo2,
			skinTemp,
			r.ScoreState,
		})
	}
	return table.Render()
}

func renderWorkouts(w io.Writer, workouts []api.Workout) error {
	table := noBorderTable(w, []string{"ID", "SPORT", "START", "DURATION", "STRAIN", "AVG HR", "STATE"})
	for _, wo := range workouts {
		strain, avgHR := "N/A", "N/A"
		duration := "N/A"
		if wo.Score != nil {
			strain = fmt.Sprintf("%.1f", wo.Score.Strain)
			avgHR = strconv.Itoa(wo.Score.AverageHeartRate)
		}
		if !wo.End.IsZero() {
			dur := wo.End.Sub(wo.Start)
			duration = formatDuration(dur)
		}
		table.Append([]string{
			wo.ID,
			strconv.Itoa(wo.SportID),
			wo.Start.Format("2006-01-02 15:04"),
			duration,
			strain,
			avgHR,
			wo.ScoreState,
		})
	}
	return table.Render()
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}
