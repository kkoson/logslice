package labelmap

import (
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("", nil)
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New(`(?P<x>[invalid`, nil)
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_NoNamedGroups(t *testing.T) {
	_, err := New(`(\w+)`, nil)
	if err == nil {
		t.Fatal("expected error when no named groups present")
	}
}

func TestNew_Valid(t *testing.T) {
	_, err := New(`(?P<level>\w+)\s+(?P<msg>.+)`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMap_Match(t *testing.T) {
	m, err := New(`(?P<level>\w+)\s+(?P<msg>.+)`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels, ok := m.Map("ERROR something went wrong")
	if !ok {
		t.Fatal("expected match")
	}
	if labels["level"] != "ERROR" {
		t.Errorf("level: got %q, want %q", labels["level"], "ERROR")
	}
	if labels["msg"] != "something went wrong" {
		t.Errorf("msg: got %q, want %q", labels["msg"], "something went wrong")
	}
}

func TestMap_NoMatch(t *testing.T) {
	m, err := New(`(?P<level>ERROR)\s+(?P<msg>.+)`, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := m.Map("DEBUG nothing here")
	if ok {
		t.Fatal("expected no match")
	}
}

func TestMap_StaticOverrides(t *testing.T) {
	overrides := map[string]string{"service": "logslice", "env": "prod"}
	m, err := New(`(?P<level>\w+)`, overrides)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels, ok := m.Map("INFO")
	if !ok {
		t.Fatal("expected match")
	}
	if labels["service"] != "logslice" {
		t.Errorf("service: got %q, want %q", labels["service"], "logslice")
	}
	if labels["env"] != "prod" {
		t.Errorf("env: got %q, want %q", labels["env"], "prod")
	}
	if labels["level"] != "INFO" {
		t.Errorf("level: got %q, want %q", labels["level"], "INFO")
	}
}

func TestMap_PatternOverridesStatic(t *testing.T) {
	overrides := map[string]string{"level": "static"}
	m, err := New(`(?P<level>\w+)`, overrides)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	labels, ok := m.Map("WARN")
	if !ok {
		t.Fatal("expected match")
	}
	// captured value takes precedence over static override
	if labels["level"] != "WARN" {
		t.Errorf("level: got %q, want %q", labels["level"], "WARN")
	}
}
