// Package aggregator provides log line counting and grouping by regex capture groups.
package aggregator

import (
	"fmt"
	"regexp"
	"sort"
	"sync"
)

// Entry holds an aggregated key and its count.
type Entry struct {
	Key   string
	Count int
}

// Aggregator groups log lines by a named capture group from a regex pattern.
type Aggregator struct {
	mu      sync.Mutex
	pattern *regexp.Regexp
	group   string
	counts  map[string]int
}

// New creates an Aggregator that groups by the named capture group in pattern.
// Returns an error if the pattern is invalid or the group name is not present.
func New(pattern, group string) (*Aggregator, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("aggregator: invalid pattern: %w", err)
	}
	found := false
	for _, name := range re.SubexpNames() {
		if name == group {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("aggregator: group %q not found in pattern", group)
	}
	return &Aggregator{
		pattern: re,
		group:   group,
		counts:  make(map[string]int),
	}, nil
}

// Add processes a single log line, extracting the capture group value and
// incrementing its count. Lines that do not match are silently ignored.
func (a *Aggregator) Add(line string) {
	match := a.pattern.FindStringSubmatch(line)
	if match == nil {
		return
	}
	key := match[a.pattern.SubexpIndex(a.group)]
	a.mu.Lock()
	a.counts[key]++
	a.mu.Unlock()
}

// Results returns aggregated entries sorted by count descending, then key ascending.
func (a *Aggregator) Results() []Entry {
	a.mu.Lock()
	defer a.mu.Unlock()
	entries := make([]Entry, 0, len(a.counts))
	for k, v := range a.counts {
		entries = append(entries, Entry{Key: k, Count: v})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Key < entries[j].Key
	})
	return entries
}

// Reset clears all accumulated counts.
func (a *Aggregator) Reset() {
	a.mu.Lock()
	a.counts = make(map[string]int)
	a.mu.Unlock()
}
