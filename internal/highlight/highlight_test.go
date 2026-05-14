package highlight_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/highlight"
)

func TestNew_InvalidPattern(t *testing.T) {
	_, err := highlight.New("[", "red")
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNew_UnknownColour(t *testing.T) {
	_, err := highlight.New("ERROR", "purple")
	if err == nil {
		t.Fatal("expected error for unknown colour, got nil")
	}
	if !strings.Contains(err.Error(), "unknown colour") {
		t.Errorf("error message should mention unknown colour, got: %v", err)
	}
}

func TestNew_Valid(t *testing.T) {
	h, err := highlight.New("ERROR", "red")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h == nil {
		t.Fatal("expected non-nil Highlighter")
	}
}

func TestApply_NoMatch_ReturnsOriginal(t *testing.T) {
	h, _ := highlight.New("ERROR", "red")
	line := "INFO everything is fine"
	if got := h.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_MatchWrapsWithANSI(t *testing.T) {
	h, _ := highlight.New("ERROR", "yellow")
	line := "2024-01-01 ERROR something broke"
	got := h.Apply(line)
	if !strings.Contains(got, highlight.Yellow+"ERROR"+highlight.Reset) {
		t.Errorf("expected ANSI-wrapped ERROR in output, got %q", got)
	}
	if !strings.Contains(got, "something broke") {
		t.Errorf("non-matched portion should be preserved, got %q", got)
	}
}

func TestApply_MultipleMatches(t *testing.T) {
	h, _ := highlight.New("WARN", "cyan")
	line := "WARN first WARN second"
	got := h.Apply(line)
	count := strings.Count(got, highlight.Cyan+"WARN"+highlight.Reset)
	if count != 2 {
		t.Errorf("expected 2 highlighted matches, got %d in %q", count, got)
	}
}

func TestApply_AllColours(t *testing.T) {
	colours := []string{"red", "green", "yellow", "blue", "cyan"}
	for _, c := range colours {
		h, err := highlight.New("test", c)
		if err != nil {
			t.Errorf("colour %q should be valid, got error: %v", c, err)
			continue
		}
		got := h.Apply("test line")
		if !strings.Contains(got, highlight.Reset) {
			t.Errorf("colour %q: expected ANSI reset in output %q", c, got)
		}
	}
}
