package transform

import (
	"fmt"
	"regexp"
	"strings"
)

// Transformer replaces regex matches within a log line with a substitution
// string. Named capture groups from the pattern are available as $name
// references in the replacement template.
type Transformer struct {
	re          *regexp.Regexp
	replacement string
}

// New compiles pattern and returns a Transformer that replaces every
// non-overlapping match with replacement. Named back-references in
// replacement use the syntax $name or ${name}.
func New(pattern, replacement string) (*Transformer, error) {
	if pattern == "" {
		return nil, fmt.Errorf("transform: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("transform: invalid pattern %q: %w", pattern, err)
	}
	return &Transformer{re: re, replacement: replacement}, nil
}

// Apply returns a copy of line with all regex matches replaced according to
// the replacement template. If the pattern does not match, the original line
// is returned unchanged.
func (t *Transformer) Apply(line string) string {
	if !t.re.MatchString(line) {
		return line
	}
	return t.re.ReplaceAllString(line, t.replacement)
}

// Chain applies a slice of Transformers in order, feeding each output into
// the next transformer.
func Chain(transformers []*Transformer, line string) string {
	result := line
	for _, tr := range transformers {
		result = tr.Apply(result)
	}
	return result
}

// MustNew is like New but panics on error. Intended for use in tests or
// package-level variable initialisation where the pattern is a constant.
func MustNew(pattern, replacement string) *Transformer {
	tr, err := New(pattern, replacement)
	if err != nil {
		panic(fmt.Sprintf("transform.MustNew: %v", err))
	}
	_ = strings.Contains // keep import for potential future helpers
	return tr
}
