// Package highlight provides ANSI terminal colour highlighting for
// log lines that match a given regular expression pattern.
package highlight

import (
	"fmt"
	"regexp"
)

// ANSI colour escape codes.
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

// validColours is the set of colour names accepted by New.
var validColours = map[string]string{
	"red":    Red,
	"green":  Green,
	"yellow": Yellow,
	"blue":   Blue,
	"cyan":   Cyan,
}

// Highlighter wraps matched substrings in a line with ANSI colour codes.
type Highlighter struct {
	re    *regexp.Regexp
	colour string
}

// New creates a Highlighter that colours matches of pattern using the named
// colour. Returns an error if the pattern is invalid or the colour is unknown.
func New(pattern, colour string) (*Highlighter, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("highlight: invalid pattern %q: %w", pattern, err)
	}
	code, ok := validColours[colour]
	if !ok {
		return nil, fmt.Errorf("highlight: unknown colour %q", colour)
	}
	return &Highlighter{re: re, colour: code}, nil
}

// Apply returns the line with every match of the compiled pattern wrapped in
// the configured ANSI colour sequence. Lines with no match are returned as-is.
func (h *Highlighter) Apply(line string) string {
	return h.re.ReplaceAllStringFunc(line, func(match string) string {
		return h.colour + match + Reset
	})
}
