package output

import (
	"encoding/json"
	"io"
)

// JSONFormatter renders data as indented JSON.
type JSONFormatter struct{}

func (f *JSONFormatter) Format(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
