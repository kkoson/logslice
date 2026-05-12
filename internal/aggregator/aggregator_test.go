package aggregator_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/aggregator"
)

func TestNew_InvalidPattern(t *testing.T) {
	_, err := aggregator.New("[", "key")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_MissingGroup(t *testing.T) {
	_, err := aggregator.New(`(?P<level>\w+)`, "missing")
	if err == nil {
		t.Fatal("expected error for missing group")
	}
}

func TestAggregator_Add_NoMatch(t *testing.T) {
	a, err := aggregator.New(`(?P<level>\w+)`, "level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.Add("   ") // should not panic or count
	if got := a.Results(); len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestAggregator_CountsCorrectly(t *testing.T) {
	a, _ := aggregator.New(`level=(?P<level>\w+)`, "level")
	lines := []string{
		"level=INFO msg=started",
		"level=ERROR msg=failed",
		"level=INFO msg=running",
		"level=INFO msg=done",
	}
	for _, l := range lines {
		a.Add(l)
	}
	results := a.Results()
	if len(results) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(results))
	}
	if results[0].Key != "INFO" || results[0].Count != 3 {
		t.Errorf("expected INFO=3, got %s=%d", results[0].Key, results[0].Count)
	}
	if results[1].Key != "ERROR" || results[1].Count != 1 {
		t.Errorf("expected ERROR=1, got %s=%d", results[1].Key, results[1].Count)
	}
}

func TestAggregator_SortTieBreakByKey(t *testing.T) {
	a, _ := aggregator.New(`svc=(?P<svc>\w+)`, "svc")
	a.Add("svc=alpha")
	a.Add("svc=beta")
	results := a.Results()
	if results[0].Key != "alpha" {
		t.Errorf("expected alpha first (tie-break), got %s", results[0].Key)
	}
}

func TestAggregator_Reset(t *testing.T) {
	a, _ := aggregator.New(`(?P<level>\w+)`, "level")
	a.Add("INFO")
	a.Reset()
	if got := a.Results(); len(got) != 0 {
		t.Fatalf("expected 0 entries after reset, got %d", len(got))
	}
}
