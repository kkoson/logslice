// Package dedup provides log line deduplication by tracking seen lines
// within a configurable sliding window or fixed-size cache.
package dedup

import "sync"

// Deduplicator tracks seen log lines and reports whether a line is a duplicate.
type Deduplicator struct {
	mu      sync.Mutex
	seen    map[string]struct{}
	window  int
	order   []string
}

// New creates a Deduplicator that remembers at most windowSize unique lines.
// Once the window is full the oldest entry is evicted. A windowSize of 0
// means unlimited (bounded only by memory).
func New(windowSize int) (*Deduplicator, error) {
	if windowSize < 0 {
		return nil, ErrInvalidWindow
	}
	return &Deduplicator{
		seen:   make(map[string]struct{}),
		window: windowSize,
	}, nil
}

// ErrInvalidWindow is returned when a negative window size is provided.
var ErrInvalidWindow = errInvalidWindow("window size must be >= 0")

type errInvalidWindow string

func (e errInvalidWindow) Error() string { return string(e) }

// IsDuplicate returns true if line has been seen before and records it if not.
func (d *Deduplicator) IsDuplicate(line string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.seen[line]; ok {
		return true
	}

	// Evict oldest entry when window is full.
	if d.window > 0 && len(d.order) >= d.window {
		oldest := d.order[0]
		d.order = d.order[1:]
		delete(d.seen, oldest)
	}

	d.seen[line] = struct{}{}
	d.order = append(d.order, line)
	return false
}

// Reset clears all tracked lines.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]struct{})
	d.order = d.order[:0]
}
