// Package fieldextract provides named-capture-group extraction from log lines,
// turning unstructured text into key/value maps for downstream processing.
package fieldextract

import (
	"fmt"
	"regexp"
)

// Extractor compiles a regex with named capture groups and extracts fields
// from each log line into a map[string]string.
type Extractor struct {
	re     *regexp.Regexp
	fields []string // ordered list of named groups
}

// New compiles pattern and returns an Extractor. pattern must contain at
// least one named capture group ((?P<name>...)); otherwise an error is
// returned.
func New(pattern string) (*Extractor, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("fieldextract: compile pattern: %w", err)
	}
	names := namedGroups(re)
	if len(names) == 0 {
		return nil, fmt.Errorf("fieldextract: pattern has no named capture groups")
	}
	return &Extractor{re: re, fields: names}, nil
}

// Extract returns a map of named capture group → matched substring for line.
// If the pattern does not match, the returned map is nil and ok is false.
func (e *Extractor) Extract(line string) (fields map[string]string, ok bool) {
	match := e.re.FindStringSubmatch(line)
	if match == nil {
		return nil, false
	}
	fields = make(map[string]string, len(e.fields))
	for _, name := range e.fields {
		idx := e.re.SubexpIndex(name)
		if idx >= 0 && idx < len(match) {
			fields[name] = match[idx]
		}
	}
	return fields, true
}

// Fields returns the ordered list of named capture groups defined in the
// compiled pattern.
func (e *Extractor) Fields() []string {
	out := make([]string, len(e.fields))
	copy(out, e.fields)
	return out
}

// namedGroups returns all named subexpression names in order.
func namedGroups(re *regexp.Regexp) []string {
	var names []string
	for _, n := range re.SubexpNames() {
		if n != "" {
			names = append(names, n)
		}
	}
	return names
}
