// Package output provides structured output formatters for logslice.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// Format represents the output format type.
type Format string

const (
	FormatPlain Format = "plain"
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
)

// LogEntry represents a single parsed log line with optional metadata.
type LogEntry struct {
	Line      string    `json:"line"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Index     int       `json:"index"`
}

// Formatter writes log entries to an io.Writer in a specific format.
type Formatter interface {
	Write(entry LogEntry) error
	Flush() error
}

// NewFormatter returns a Formatter for the given format and writer.
// Returns an error if the format is unsupported.
func NewFormatter(format Format, w io.Writer) (Formatter, error) {
	switch format {
	case FormatPlain:
		return &plainFormatter{w: w}, nil
	case FormatJSON:
		return &jsonFormatter{w: w}, nil
	case FormatCSV:
		return &csvFormatter{w: w, headerWritten: false}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %q", format)
	}
}

// plainFormatter writes each log line as-is.
type plainFormatter struct {
	w io.Writer
}

func (f *plainFormatter) Write(entry LogEntry) error {
	_, err := fmt.Fprintln(f.w, entry.Line)
	return err
}

func (f *plainFormatter) Flush() error { return nil }

// jsonFormatter writes each log entry as a JSON object.
type jsonFormatter struct {
	w io.Writer
}

func (f *jsonFormatter) Write(entry LogEntry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}
	_, err = fmt.Fprintln(f.w, string(data))
	return err
}

func (f *jsonFormatter) Flush() error { return nil }

// csvFormatter writes log entries as CSV rows.
type csvFormatter struct {
	w             io.Writer
	headerWritten bool
}

func (f *csvFormatter) Write(entry LogEntry) error {
	if !f.headerWritten {
		if _, err := fmt.Fprintln(f.w, "index,timestamp,line"); err != nil {
			return err
		}
		f.headerWritten = true
	}
	ts := entry.Timestamp.Format(time.RFC3339)
	escaped := strings.ReplaceAll(entry.Line, `"`, `""`)
	_, err := fmt.Fprintf(f.w, "%d,%s,\"%s\"\n", entry.Index, ts, escaped)
	return err
}

func (f *csvFormatter) Flush() error { return nil }
