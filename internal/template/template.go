// Package template provides line-level text templating using named capture
// groups from a regex pattern. Each matched group is substituted into a
// Go text/template string, enabling structured reformatting of log lines.
package template

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"
)

// Renderer compiles a regex and a Go template once and applies them to
// individual log lines. Named capture groups from the regex become the
// template data map.
type Renderer struct {
	re   *regexp.Regexp
	tmpl *template.Template
}

// New returns a Renderer that matches lines with pattern and formats them
// with tmplStr. pattern must contain at least one named capture group.
// Lines that do not match are returned unchanged.
func New(pattern, tmplStr string) (*Renderer, error) {
	if pattern == "" {
		return nil, fmt.Errorf("template: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("template: invalid pattern: %w", err)
	}
	if len(namedGroups(re)) == 0 {
		return nil, fmt.Errorf("template: pattern must contain at least one named capture group")
	}
	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("template: invalid template: %w", err)
	}
	return &Renderer{re: re, tmpl: tmpl}, nil
}

// Apply executes the template against named groups captured from line.
// If line does not match the pattern, line is returned unchanged.
func (r *Renderer) Apply(line string) (string, error) {
	match := r.re.FindStringSubmatch(line)
	if match == nil {
		return line, nil
	}
	data := make(map[string]string, len(r.re.SubexpNames()))
	for i, name := range r.re.SubexpNames() {
		if name != "" {
			data[name] = match[i]
		}
	}
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template: execute: %w", err)
	}
	return buf.String(), nil
}

// namedGroups returns the named subexpression names for re.
func namedGroups(re *regexp.Regexp) []string {
	var names []string
	for _, n := range re.SubexpNames() {
		if n != "" {
			names = append(names, n)
		}
	}
	return names
}
