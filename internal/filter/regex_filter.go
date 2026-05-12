package filter

import (
	"fmt"
	"regexp"
)

// Filter represents a compiled regex filter that can include or exclude log lines.
type Filter struct {
	pattern *regexp.Regexp
	invert  bool
}

// NewFilter creates a new Filter from a regex pattern string.
// If invert is true, lines matching the pattern are excluded instead of included.
func NewFilter(pattern string, invert bool) (*Filter, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
	}
	return &Filter{
		pattern: re,
		invert:  invert,
	}, nil
}

// Match reports whether the line passes the filter.
func (f *Filter) Match(line string) bool {
	matched := f.pattern.MatchString(line)
	if f.invert {
		return !matched
	}
	return matched
}

// Pipeline holds an ordered sequence of filters applied to each log line.
type Pipeline struct {
	filters []*Filter
}

// NewPipeline creates an empty filter pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// Add appends a filter to the pipeline.
func (p *Pipeline) Add(f *Filter) {
	p.filters = append(p.filters, f)
}

// Match returns true only if the line passes every filter in the pipeline.
func (p *Pipeline) Match(line string) bool {
	for _, f := range p.filters {
		if !f.Match(line) {
			return false
		}
	}
	return true
}

// Len returns the number of filters in the pipeline.
func (p *Pipeline) Len() int {
	return len(p.filters)
}
