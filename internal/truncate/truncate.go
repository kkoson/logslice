package truncate

import "fmt"

// Truncator shortens log lines that exceed a maximum byte length,
// optionally appending a configurable suffix to indicate truncation.
type Truncator struct {
	maxLen int
	suffix string
}

// New creates a Truncator with the given maximum line length and suffix.
// maxLen must be greater than zero and greater than len(suffix).
func New(maxLen int, suffix string) (*Truncator, error) {
	if maxLen <= 0 {
		return nil, fmt.Errorf("truncate: maxLen must be > 0, got %d", maxLen)
	}
	if len(suffix) >= maxLen {
		return nil, fmt.Errorf("truncate: suffix length %d must be less than maxLen %d", len(suffix), maxLen)
	}
	return &Truncator{maxLen: maxLen, suffix: suffix}, nil
}

// Apply returns the line unchanged if it fits within maxLen bytes.
// Otherwise it trims the line to (maxLen - len(suffix)) bytes and appends the suffix.
func (t *Truncator) Apply(line string) string {
	if len(line) <= t.maxLen {
		return line
	}
	cutAt := t.maxLen - len(t.suffix)
	// Ensure we do not cut in the middle of a multi-byte rune.
	for cutAt > 0 && !isRuneBoundary(line, cutAt) {
		cutAt--
	}
	return line[:cutAt] + t.suffix
}

// isRuneBoundary reports whether position i is a valid UTF-8 rune boundary in s.
func isRuneBoundary(s string, i int) bool {
	if i == 0 || i == len(s) {
		return true
	}
	// A byte is a rune boundary if it is not a UTF-8 continuation byte (10xxxxxx).
	return s[i]&0xC0 != 0x80
}
