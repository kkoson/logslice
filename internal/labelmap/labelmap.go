// Package labelmap provides key-value label enrichment for log lines.
// A Mapper matches a log line against a regular expression and attaches
// named capture groups as structured labels, enabling downstream formatters
// and aggregators to consume pre-parsed metadata without re-parsing.
package labelmap

import (
	"fmt"
	"regexp"
)

// Labels is an ordered map of label name → value extracted from a log line.
type Labels map[string]string

// Mapper extracts named labels from log lines via a compiled regular expression.
type Mapper struct {
	re     *regexp.Regexp
	groups []string // ordered named capture group names
}

// New compiles pattern and returns a Mapper ready to label log lines.
// pattern must contain at least one named capture group (e.g. (?P<level>\w+)).
// Returns an error if pattern is invalid or contains no named groups.
func New(pattern string) (*Mapper, error) {
	if pattern == "" {
		return nil, fmt.Errorf("labelmap: pattern must not be empty")
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("labelmap: invalid pattern: %w", err)
	}

	groups := namedGroups(re)
	if len(groups) == 0 {
		return nil, fmt.Errorf("labelmap: pattern contains no named capture groups")
	}

	return &Mapper{re: re, groups: groups}, nil
}

// Apply matches line against the compiled pattern and returns a Labels map
// populated with every named capture group that participated in the match.
// If the line does not match, Apply returns nil, nil — callers should treat
// a nil Labels as "no labels available" rather than an error condition.
func (m *Mapper) Apply(line string) (Labels, error) {
	match := m.re.FindStringSubmatch(line)
	if match == nil {
		return nil, nil
	}

	labels := make(Labels, len(m.groups))
	for _, name := range m.groups {
		idx := m.re.SubexpIndex(name)
		if idx >= 0 && idx < len(match) {
			labels[name] = match[idx]
		}
	}
	return labels, nil
}

// Groups returns the ordered list of named capture groups that this Mapper
// will populate. Useful for consumers that need to know the label schema
// before processing any lines (e.g. CSV header generation).
func (m *Mapper) Groups() []string {
	out := make([]string, len(m.groups))
	copy(out, m.groups)
	return out
}

// namedGroups returns the named subexpression names for re, preserving
// declaration order and skipping the implicit empty name at index 0.
func namedGroups(re *regexp.Regexp) []string {
	var groups []string
	for _, name := range re.SubexpNames() {
		if name != "" {
			groups = append(groups, name)
		}
	}
	return groups
}
