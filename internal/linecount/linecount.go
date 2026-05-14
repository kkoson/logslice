package linecount

import "sync/atomic"

// Counter tracks the number of lines seen and emitted during a pipeline run.
// It is safe for concurrent use.
type Counter struct {
	seen    atomic.Int64
	emitted atomic.Int64
}

// New returns a new zero-value Counter.
func New() *Counter {
	return &Counter{}
}

// IncSeen records that one more line was read from the source.
func (c *Counter) IncSeen() {
	c.seen.Add(1)
}

// IncEmitted records that one more line passed all filters and was written.
func (c *Counter) IncEmitted() {
	c.emitted.Add(1)
}

// Seen returns the total number of lines read from the source.
func (c *Counter) Seen() int64 {
	return c.seen.Load()
}

// Emitted returns the total number of lines that passed all filters.
func (c *Counter) Emitted() int64 {
	return c.emitted.Load()
}

// Dropped returns the number of lines that were filtered out.
func (c *Counter) Dropped() int64 {
	return c.Seen() - c.Emitted()
}

// Summary holds a snapshot of counter values for reporting.
type Summary struct {
	Seen    int64
	Emitted int64
	Dropped int64
}

// Snapshot returns an immutable Summary of the current counts.
func (c *Counter) Snapshot() Summary {
	s := c.Seen()
	e := c.Emitted()
	return Summary{
		Seen:    s,
		Emitted: e,
		Dropped: s - e,
	}
}
