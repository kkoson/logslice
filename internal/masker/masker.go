package masker

import (
	"fmt"
	"regexp"
)

// Masker replaces sensitive fields in log lines with a fixed mask string.
type Masker struct {
	re   *regexp.Regexp
	mask string
}

// New creates a Masker that replaces every match of pattern with mask.
// pattern must contain at least one named capture group; only the groups
// are replaced, the surrounding text is preserved.
//
// Returns an error if pattern is empty, invalid, or has no named groups.
func New(pattern, mask string) (*Masker, error) {
	if pattern == "" {
		return nil, fmt.Errorf("masker: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("masker: invalid pattern: %w", err)
	}
	if len(namedGroups(re)) == 0 {
		return nil, fmt.Errorf("masker: pattern must contain at least one named capture group")
	}
	if mask == "" {
		mask = "***"
	}
	return &Masker{re: re, mask: mask}, nil
}

// Apply returns a copy of line with every named capture group replaced by
// the configured mask. Non-capturing and unnamed groups are left intact.
func (m *Masker) Apply(line string) string {
	names := namedGroups(m.re)
	set := make(map[string]struct{}, len(names))
	for _, n := range names {
		set[n] = struct{}{}
	}

	return m.re.ReplaceAllStringFunc(line, func(match string) string {
		sub := m.re.FindStringSubmatchIndex(match)
		if sub == nil {
			return match
		}
		result := []byte(match)
		// Walk named groups in reverse so index offsets stay valid.
		for i := len(names) - 1; i >= 0; i-- {
			name := names[i]
			if _, ok := set[name]; !ok {
				continue
			}
			gIdx := m.re.SubexpIndex(name)
			start, end := sub[2*gIdx], sub[2*gIdx+1]
			if start < 0 {
				continue
			}
			replaced := make([]byte, 0, len(result)-(end-start)+len(m.mask))
			replaced = append(replaced, result[:start]...)
			replaced = append(replaced, m.mask...)
			replaced = append(replaced, result[end:]...)
			result = replaced
			// Recompute submatch indices on updated slice.
			sub = m.re.FindSubmatchIndex(result)
			if sub == nil {
				break
			}
		}
		return string(result)
	})
}

// namedGroups returns the ordered list of named subexpressions in re.
func namedGroups(re *regexp.Regexp) []string {
	var out []string
	for _, n := range re.SubexpNames() {
		if n != "" {
			out = append(out, n)
		}
	}
	return out
}
