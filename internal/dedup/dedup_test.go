package dedup

import (
	"testing"
)

func TestNew_NegativeWindow(t *testing.T) {
	_, err := New(-1)
	if err == nil {
		t.Fatal("expected error for negative window, got nil")
	}
}

func TestNew_ZeroWindow(t *testing.T) {
	d, err := New(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Deduplicator")
	}
}

func TestIsDuplicate_FirstOccurrence(t *testing.T) {
	d, _ := New(10)
	if d.IsDuplicate("hello") {
		t.Error("first occurrence should not be a duplicate")
	}
}

func TestIsDuplicate_SecondOccurrence(t *testing.T) {
	d, _ := New(10)
	d.IsDuplicate("hello")
	if !d.IsDuplicate("hello") {
		t.Error("second occurrence should be a duplicate")
	}
}

func TestIsDuplicate_DifferentLines(t *testing.T) {
	d, _ := New(10)
	d.IsDuplicate("line1")
	if d.IsDuplicate("line2") {
		t.Error("different line should not be a duplicate")
	}
}

func TestWindow_EvictsOldest(t *testing.T) {
	// window of 2: after adding "a" and "b", adding "c" evicts "a".
	d, _ := New(2)
	d.IsDuplicate("a")
	d.IsDuplicate("b")
	d.IsDuplicate("c") // evicts "a"

	// "a" should no longer be tracked — not a duplicate.
	if d.IsDuplicate("a") {
		t.Error("evicted line 'a' should not be reported as duplicate")
	}
	// "b" is still in the window.
	if !d.IsDuplicate("b") {
		t.Error("line 'b' still in window should be reported as duplicate")
	}
}

func TestReset_ClearsState(t *testing.T) {
	d, _ := New(10)
	d.IsDuplicate("hello")
	d.Reset()
	if d.IsDuplicate("hello") {
		t.Error("after reset, line should not be a duplicate")
	}
}

func TestUnlimitedWindow(t *testing.T) {
	d, _ := New(0)
	for i := 0; i < 1000; i++ {
		line := string(rune('a' + i%26))
		_ = d.IsDuplicate(line)
	}
	// All single-letter lines should now be duplicates.
	if !d.IsDuplicate("a") {
		t.Error("expected 'a' to be a duplicate in unlimited window")
	}
}
