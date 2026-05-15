// Package grok provides named-pattern parsing for common log formats.
// It maps symbolic pattern names (e.g. IP, TIMESTAMP) to regular expressions
// and compiles them into a single extractor that returns named fields.
package grok

import (
	"fmt"
	"regexp"
	"strings"
)

// built-in pattern library.
var builtins = map[string]string{
	"IP":        `(?P<IP>\d{1,3}(?:\.\d{1,3}){3})`,
	"TIMESTAMP": `(?P<TIMESTAMP>\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:?\d{2})?)`,
	"LOGLEVEL":  `(?P<LOGLEVEL>DEBUG|INFO|WARN|ERROR|FATAL)`,
	"WORD":      `(?P<WORD>\S+)`,
	"NUMBER":    `(?P<NUMBER>[+-]?(?:\d+\.?\d*|\.\d+))`,
	"GREEDYDATA": `(?P<GREEDYDATA>.*)`,
}

// Parser compiles a grok-style pattern string into a regexp extractor.
type Parser struct {
	re      *regexp.Regexp
	fields  []string
}

// New creates a Parser from a pattern string that may contain %{NAME} tokens.
// Unknown pattern names return an error.
func New(pattern string, extra map[string]string) (*Parser, error) {
	if pattern == "" {
		return nil, fmt.Errorf("grok: pattern must not be empty")
	}

	lib := make(map[string]string, len(builtins)+len(extra))
	for k, v := range builtins {
		lib[k] = v
	}
	for k, v := range extra {
		lib[k] = v
	}

	expanded, err := expand(pattern, lib)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(expanded)
	if err != nil {
		return nil, fmt.Errorf("grok: compiled pattern invalid: %w", err)
	}

	return &Parser{re: re, fields: re.SubexpNames()}, nil
}

// Parse extracts named fields from line. Returns nil if the pattern does not match.
func (p *Parser) Parse(line string) map[string]string {
	m := p.re.FindStringSubmatch(line)
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(p.fields))
	for i, name := range p.fields {
		if name != "" {
			out[name] = m[i]
		}
	}
	return out
}

// expand replaces %{NAME} tokens with their regexp equivalents.
func expand(pattern string, lib map[string]string) (string, error) {
	var err error
	result := regexp.MustCompile(`%\{(\w+)\}`).ReplaceAllStringFunc(pattern, func(tok string) string {
		name := tok[2 : len(tok)-1]
		v, ok := lib[name]
		if !ok {
			err = fmt.Errorf("grok: unknown pattern %q", name)
			return tok
		}
		return v
	})
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(result), nil
}
