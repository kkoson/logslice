// Package sampler provides rate-based log line sampling for logslice.
// It allows processing only a fraction of matching lines, which is useful
// for high-volume log streams where full analysis is not required.
package sampler

import (
	"errors"
	"math/rand"
)

// Sampler decides whether a given log line should be kept based on a
// configured sampling rate in the range (0, 1].
type Sampler struct {
	rate float64
	rng  *rand.Rand
}

// New creates a Sampler with the given rate. rate must be in (0, 1].
// A rate of 1.0 keeps every line; 0.1 keeps roughly 10% of lines.
func New(rate float64, src rand.Source) (*Sampler, error) {
	if rate <= 0 || rate > 1 {
		return nil, errors.New("sampler: rate must be in (0, 1]")
	}
	if src == nil {
		src = rand.NewSource(42)
	}
	return &Sampler{rate: rate, rng: rand.New(src)}, nil
}

// Keep returns true if the line should be kept according to the sampling rate.
func (s *Sampler) Keep(_ string) bool {
	return s.rng.Float64() < s.rate
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() float64 {
	return s.rate
}
