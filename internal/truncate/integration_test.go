package truncate_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/truncate"
)

// TestTruncator_Pipeline simulates a stream of log lines being truncated.
func TestTruncator_Pipeline(t *testing.T) {
	const maxLen = 40
	const suffix = " [truncated]"

	tr, err := truncate.New(maxLen, suffix)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	lines := []string{
		"short",
		strings.Repeat("a", 40),
		strings.Repeat("b", 80),
		"2024-01-15T12:00:00Z ERROR something went catastrophically wrong in module X",
	}

	for _, line := range lines {
		result := tr.Apply(line)
		if len(result) > maxLen {
			t.Errorf("line %q: result length %d exceeds maxLen %d", line[:min(len(line), 20)], len(result), maxLen)
		}
		if len(line) > maxLen && !strings.HasSuffix(result, suffix) {
			t.Errorf("long line was not suffixed: got %q", result)
		}
		if len(line) <= maxLen && result != line {
			t.Errorf("short line was modified: got %q want %q", result, line)
		}
	}
}

func TestTruncator_AllASCII_ExactBoundary(t *testing.T) {
	tr, _ := truncate.New(20, ">>")
	line := strings.Repeat("x", 20)
	if got := tr.Apply(line); got != line {
		t.Errorf("exact-length line should be unchanged, got %q", got)
	}
	line2 := strings.Repeat("x", 21)
	got2 := tr.Apply(line2)
	if len(got2) != 20 {
		t.Errorf("expected length 20, got %d", len(got2))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
