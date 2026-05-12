package filter

import (
	"testing"
)

func TestNewFilter_InvalidPattern(t *testing.T) {
	_, err := NewFilter("[invalid", false)
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestFilter_Match(t *testing.T) {
	tests := []struct {
		pattern string
		invert  bool
		line    string
		want    bool
	}{
		{"ERROR", false, "2024-01-01 ERROR something failed", true},
		{"ERROR", false, "2024-01-01 INFO all good", false},
		{"ERROR", true, "2024-01-01 INFO all good", true},
		{"ERROR", true, "2024-01-01 ERROR something failed", false},
		{`\d{3}`, false, "response code 404 not found", true},
		{`\d{3}`, false, "no digits here", false},
	}

	for _, tc := range tests {
		f, err := NewFilter(tc.pattern, tc.invert)
		if err != nil {
			t.Fatalf("NewFilter(%q): unexpected error: %v", tc.pattern, err)
		}
		got := f.Match(tc.line)
		if got != tc.want {
			t.Errorf("Filter{pattern:%q, invert:%v}.Match(%q) = %v, want %v",
				tc.pattern, tc.invert, tc.line, got, tc.want)
		}
	}
}

func TestPipeline_EmptyMatchesAll(t *testing.T) {
	p := NewPipeline()
	if !p.Match("any line should pass") {
		t.Error("empty pipeline should match every line")
	}
}

func TestPipeline_MultipleFilters(t *testing.T) {
	p := NewPipeline()

	f1, _ := NewFilter("ERROR", false)
	f2, _ := NewFilter("timeout", false)
	p.Add(f1)
	p.Add(f2)

	if p.Len() != 2 {
		t.Errorf("expected pipeline length 2, got %d", p.Len())
	}

	if !p.Match("ERROR: connection timeout reached") {
		t.Error("expected match for line containing both ERROR and timeout")
	}
	if p.Match("ERROR: disk full") {
		t.Error("expected no match for line missing 'timeout'")
	}
	if p.Match("INFO: timeout warning") {
		t.Error("expected no match for line missing 'ERROR'")
	}
}

func TestPipeline_InvertFilter(t *testing.T) {
	p := NewPipeline()

	f1, _ := NewFilter("ERROR", false)
	f2, _ := NewFilter("DEBUG", true) // exclude DEBUG lines
	p.Add(f1)
	p.Add(f2)

	if !p.Match("ERROR: something went wrong") {
		t.Error("expected match for ERROR line without DEBUG")
	}
	if p.Match("ERROR DEBUG: verbose error") {
		t.Error("expected no match for line containing both ERROR and DEBUG")
	}
}
