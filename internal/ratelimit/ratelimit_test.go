package ratelimit

import (
	"testing"
	"time"
)

func TestNew_ZeroRate(t *testing.T) {
	_, err := New(0)
	if err == nil {
		t.Fatal("expected error for zero rate, got nil")
	}
}

func TestNew_NegativeRate(t *testing.T) {
	_, err := New(-5)
	if err == nil {
		t.Fatal("expected error for negative rate, got nil")
	}
}

func TestNew_ValidRate(t *testing.T) {
	l, err := New(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.Rate() != 10 {
		t.Fatalf("expected rate 10, got %d", l.Rate())
	}
}

func TestAllow_FullBucketPermitsUpToRate(t *testing.T) {
	l, _ := New(5)
	allowed := 0
	for i := 0; i < 10; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 5 {
		t.Fatalf("expected 5 allowed, got %d", allowed)
	}
}

func TestAllow_RefillAfterOneSec(t *testing.T) {
	base := time.Now()
	l, _ := New(3)
	l.clock = func() time.Time { return base }

	// Drain the bucket.
	for i := 0; i < 3; i++ {
		l.Allow()
	}
	if l.Allow() {
		t.Fatal("expected bucket to be empty")
	}

	// Advance clock by 1 second to trigger refill.
	l.clock = func() time.Time { return base.Add(time.Second) }
	if !l.Allow() {
		t.Fatal("expected token after refill")
	}
}

func TestAllow_NoRefillBeforeOneSec(t *testing.T) {
	base := time.Now()
	l, _ := New(2)
	l.clock = func() time.Time { return base }

	l.Allow()
	l.Allow()

	// Advance clock by less than a second.
	l.clock = func() time.Time { return base.Add(500 * time.Millisecond) }
	if l.Allow() {
		t.Fatal("expected no token before full second elapsed")
	}
}
