package multiline

import (
	"regexp"
	"strings"
)

// Folder accumulates log lines that belong to a single logical event.
// A new event begins whenever a line matches the start pattern.
// Lines that do not match the start pattern are folded into the
// current event (e.g. Java stack traces, multi-line JSON blocks).
type Folder struct {
	start   *regexp.Regexp
	join    string
	buf     []string
}

// New returns a Folder that treats any line matching startPattern as
// the beginning of a new event. join is the string used to concatenate
// continuation lines (typically "\n" or " ").
func New(startPattern, join string) (*Folder, error) {
	if startPattern == "" {
		return nil, ErrEmptyPattern
	}
	re, err := regexp.Compile(startPattern)
	if err != nil {
		return nil, err
	}
	return &Folder{start: re, join: join}, nil
}

// Add feeds the next raw line to the folder.
// If the line starts a new event and a previous event was buffered,
// the completed event is returned together with ok=true.
// Otherwise ok is false and the caller should continue feeding lines.
func (f *Folder) Add(line string) (event string, ok bool) {
	if f.start.MatchString(line) {
		if len(f.buf) > 0 {
			event = strings.Join(f.buf, f.join)
			f.buf = []string{line}
			return event, true
		}
		f.buf = []string{line}
		return "", false
	}
	f.buf = append(f.buf, line)
	return "", false
}

// Flush returns any buffered event that has not yet been emitted.
// It must be called after the input stream is exhausted.
func (f *Folder) Flush() (event string, ok bool) {
	if len(f.buf) == 0 {
		return "", false
	}
	event = strings.Join(f.buf, f.join)
	f.buf = nil
	return event, true
}
