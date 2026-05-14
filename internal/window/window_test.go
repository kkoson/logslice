package window

import (
	"testing"
	"time"
)

func TestNew_InvalidSize(t *testing.T) {
	_, err := New(-time.Second, 10)
	if err == nil {
		t.Fatal("expected error for negative size")
	}
}

func TestNew_ZeroBuckets(t *testing.T) {
	_, err := New(time.Minute, 0)
	if err == nil {
		t.Fatal("expected error for zero buckets")
	}
}

func TestNew_Valid(t *testing.T) {
	w, err := New(time.Minute, 6)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil window")
	}
}

func TestTotal_EmptyWindow(t *testing.T) {
	w, _ := New(time.Minute, 6)
	if got := w.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_AccumulatesInSameBucket(t *testing.T) {
	w, _ := New(time.Minute, 6)
	w.Add(3)
	w.Add(7)
	if got := w.Total(); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestRotate_ExpiredBucketsAreZeroed(t *testing.T) {
	// Use a very small window so we can force expiry.
	w, _ := New(100*time.Millisecond, 2)
	w.Add(5)
	if w.Total() != 5 {
		t.Fatal("expected 5 immediately after add")
	}
	// Sleep past the full window.
	time.Sleep(150 * time.Millisecond)
	if got := w.Total(); got != 0 {
		t.Fatalf("expected 0 after window expires, got %d", got)
	}
}

func TestRotate_PartialExpiry(t *testing.T) {
	// 4 buckets of 25 ms each = 100 ms window.
	w, _ := New(100*time.Millisecond, 4)
	w.Add(10)
	// Sleep past 2 bucket ticks (≈50 ms) but stay within the window.
	time.Sleep(55 * time.Millisecond)
	w.Add(4)
	// Both adds should still be within the 100 ms window.
	if got := w.Total(); got < 4 {
		t.Fatalf("expected at least 4 within partial window, got %d", got)
	}
}

func TestTotal_ConcurrentSafety(t *testing.T) {
	w, _ := New(time.Second, 10)
	done := make(chan struct{})
	for i := 0; i < 20; i++ {
		go func() {
			w.Add(1)
			_ = w.Total()
			done <- struct{}{}
		}()
	}
	for i := 0; i < 20; i++ {
		<-done
	}
}
