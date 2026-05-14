package ratelimit

import (
	"errors"
	"sync"
	"time"
)

// Limiter enforces a maximum number of lines emitted per second.
type Limiter struct {
	mu       sync.Mutex
	rate     int
	bucket   int
	lastTick time.Time
	clock    func() time.Time
}

// New creates a Limiter that allows at most ratePerSec lines per second.
// ratePerSec must be greater than zero.
func New(ratePerSec int) (*Limiter, error) {
	if ratePerSec <= 0 {
		return nil, errors.New("ratelimit: rate must be greater than zero")
	}
	return &Limiter{
		rate:     ratePerSec,
		bucket:   ratePerSec,
		lastTick: time.Now(),
		clock:    time.Now,
	}, nil
}

// Allow returns true if the current line should be forwarded, false if it
// should be dropped to stay within the configured rate.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.clock()
	elapsed := now.Sub(l.lastTick)

	// Refill tokens proportional to elapsed time.
	if elapsed >= time.Second {
		seconds := int(elapsed.Seconds())
		l.bucket += seconds * l.rate
		if l.bucket > l.rate {
			l.bucket = l.rate
		}
		l.lastTick = now
	}

	if l.bucket <= 0 {
		return false
	}
	l.bucket--
	return true
}

// Rate returns the configured rate limit.
func (l *Limiter) Rate() int {
	return l.rate
}
