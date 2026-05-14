package truncate

import (
	"strings"
	"testing"
)

func TestNew_ZeroMaxLen(t *testing.T) {
	_, err := New(0, "...")
	if err == nil {
		t.Fatal("expected error for maxLen=0")
	}
}

func TestNew_NegativeMaxLen(t *testing.T) {
	_, err := New(-5, "...")
	if err == nil {
		t.Fatal("expected error for negative maxLen")
	}
}

func TestNew_SuffixTooLong(t *testing.T) {
	_, err := New(3, "...")
	if err == nil {
		t.Fatal("expected error when suffix length >= maxLen")
	}
}

func TestNew_Valid(t *testing.T) {
	tr, err := New(80, "...")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Truncator")
	}
}

func TestApply_ShortLine_Unchanged(t *testing.T) {
	tr, _ := New(80, "...")
	line := "short line"
	if got := tr.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_ExactLength_Unchanged(t *testing.T) {
	tr, _ := New(10, "...")
	line := "0123456789" // exactly 10 bytes
	if got := tr.Apply(line); got != line {
		t.Errorf("expected %q, got %q", line, got)
	}
}

func TestApply_LongLine_Truncated(t *testing.T) {
	tr, _ := New(10, "...")
	line := "0123456789EXTRA"
	got := tr.Apply(line)
	if len(got) > 10 {
		t.Errorf("result length %d exceeds maxLen 10", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Errorf("expected suffix '...', got %q", got)
	}
}

func TestApply_NoSuffix(t *testing.T) {
	tr, _ := New(5, "")
	line := "hello world"
	got := tr.Apply(line)
	if got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}

func TestApply_MultibyteRune_SafeCut(t *testing.T) {
	// Each Japanese character is 3 bytes; total = 9 bytes.
	line := "日本語"
	// maxLen=5 with suffix="…" (3 bytes) => cutAt=2, but 2 is mid-rune => back to 0.
	// Result should be suffix only.
	tr, _ := New(5, "…")
	got := tr.Apply(line)
	// Should not panic and should end with suffix.
	if !strings.HasSuffix(got, "…") {
		t.Errorf("expected suffix '…', got %q", got)
	}
}
