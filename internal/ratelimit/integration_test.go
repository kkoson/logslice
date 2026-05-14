package ratelimit_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/ratelimit"
)

// TestLimiter_RateOne_AllowsExactlyOne verifies that a rate-1 limiter allows
// exactly one line from a burst of many.
func TestLimiter_RateOne_AllowsExactlyOne(t *testing.T) {
	l, err := ratelimit.New(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	allowed := 0
	for i := 0; i < 100; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 1 {
		t.Fatalf("rate-1 limiter: expected 1 allowed, got %d", allowed)
	}
}

// TestLimiter_HighRate_AllowsAll verifies that when the burst size equals the
// number of lines, all lines are forwarded.
func TestLimiter_HighRate_AllowsAll(t *testing.T) {
	const lines = 50
	l, err := ratelimit.New(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	allowed := 0
	for i := 0; i < lines; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != lines {
		t.Fatalf("expected all %d lines allowed, got %d", lines, allowed)
	}
}

// TestLimiter_RateProperty checks that allowed count never exceeds the rate.
func TestLimiter_RateProperty(t *testing.T) {
	rates := []int{1, 5, 10, 100}
	for _, r := range rates {
		l, _ := ratelimit.New(r)
		allowed := 0
		for i := 0; i < r*3; i++ {
			if l.Allow() {
				allowed++
			}
		}
		if allowed > r {
			t.Errorf("rate %d: allowed %d exceeds rate", r, allowed)
		}
	}
}
