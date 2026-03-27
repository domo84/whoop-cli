package output

import "io"

// Format specifies the output format.
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
)

// Formatter renders data to the given io.Writer.
type Formatter interface {
	Format(w io.Writer, data any) error
}

// New returns the Formatter for the given format string.
// Defaults to table if the format is unrecognised.
func New(f string) Formatter {
	switch Format(f) {
	case FormatJSON:
		return &JSONFormatter{}
	default:
		return &TableFormatter{}
	}
}
