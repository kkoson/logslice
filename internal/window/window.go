package window

import (
	"errors"
	"sync"
	"time"
)

// Window is a sliding time-window counter. It tracks how many lines were
// seen within the last Duration and exposes a snapshot for reporting.
type Window struct {
	mu       sync.Mutex
	size     time.Duration
	buckets  []int64
	nBuckets int
	tick     time.Duration
	last     time.Time
}

// New creates a Window that spans size, divided into nBuckets equal slots.
// size must be positive and nBuckets must be >= 1.
func New(size time.Duration, nBuckets int) (*Window, error) {
	if size <= 0 {
		return nil, errors.New("window: size must be positive")
	}
	if nBuckets < 1 {
		return nil, errors.New("window: nBuckets must be >= 1")
	}
	return &Window{
		size:     size,
		nBuckets: nBuckets,
		buckets:  make([]int64, nBuckets),
		tick:     size / time.Duration(nBuckets),
		last:     time.Now(),
	}, nil
}

// Add records n events at the current time.
func (w *Window) Add(n int64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.rotate(time.Now())
	w.buckets[0] += n
}

// Total returns the sum of all events within the window.
func (w *Window) Total() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.rotate(time.Now())
	var sum int64
	for _, v := range w.buckets {
		sum += v
	}
	return sum
}

// rotate advances expired buckets, zeroing them out.
func (w *Window) rotate(now time.Time) {
	elapsed := now.Sub(w.last)
	if elapsed < w.tick {
		return
	}
	slots := int(elapsed / w.tick)
	if slots >= w.nBuckets {
		for i := range w.buckets {
			w.buckets[i] = 0
		}
	} else {
		for i := 0; i < slots; i++ {
			// shift right, dropping the oldest
			copy(w.buckets[1:], w.buckets[:w.nBuckets-1])
			w.buckets[0] = 0
		}
	}
	w.last = now
}
