package labelmap

import (
	"errors"
	"regexp"
)

// Mapper extracts named capture groups from a log line and returns them as a
// label map. Static overrides are merged with lower priority than captured
// values so that dynamic data always wins.
type Mapper struct {
	re        *regexp.Regexp
	groups    []string
	overrides map[string]string
}

// New compiles pattern and returns a Mapper. pattern must contain at least one
// named capture group. overrides is an optional set of static key=value pairs
// that are added to every successful match result.
func New(pattern string, overrides map[string]string) (*Mapper, error) {
	if pattern == "" {
		return nil, errors.New("labelmap: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	groups := namedGroups(re)
	if len(groups) == 0 {
		return nil, errors.New("labelmap: pattern must contain at least one named capture group")
	}
	ovrCopy := make(map[string]string, len(overrides))
	for k, v := range overrides {
		ovrCopy[k] = v
	}
	return &Mapper{re: re, groups: groups, overrides: ovrCopy}, nil
}

// Map attempts to match line against the compiled pattern. On success it
// returns a label map populated with static overrides and captured values
// (captured values take precedence). Returns false when the line does not
// match.
func (m *Mapper) Map(line string) (map[string]string, bool) {
	match := m.re.FindStringSubmatch(line)
	if match == nil {
		return nil, false
	}
	labels := make(map[string]string, len(m.overrides)+len(m.groups))
	for k, v := range m.overrides {
		labels[k] = v
	}
	for i, name := range m.re.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		labels[name] = match[i]
	}
	return labels, true
}

// namedGroups returns the list of named subexpression names for re.
func namedGroups(re *regexp.Regexp) []string {
	var out []string
	for _, name := range re.SubexpNames() {
		if name != "" {
			out = append(out, name)
		}
	}
	return out
}
