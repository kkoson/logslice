package fieldextract

import (
	"testing"
)

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("[invalid")
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNew_NoNamedGroups(t *testing.T) {
	_, err := New(`(\d+) (\w+)`)
	if err == nil {
		t.Fatal("expected error when no named groups present")
	}
}

func TestNew_Valid(t *testing.T) {
	e, err := New(`(?P<level>\w+) (?P<msg>.*)`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := e.Fields()
	if len(fields) != 2 || fields[0] != "level" || fields[1] != "msg" {
		t.Fatalf("unexpected fields: %v", fields)
	}
}

func TestExtract_Match(t *testing.T) {
	e, _ := New(`(?P<level>\w+)\s+(?P<msg>.*)`)
	fields, ok := e.Extract("ERROR something went wrong")
	if !ok {
		t.Fatal("expected match")
	}
	if fields["level"] != "ERROR" {
		t.Errorf("level: got %q, want %q", fields["level"], "ERROR")
	}
	if fields["msg"] != "something went wrong" {
		t.Errorf("msg: got %q, want %q", fields["msg"], "something went wrong")
	}
}

func TestExtract_NoMatch(t *testing.T) {
	e, _ := New(`^(?P<ts>\d{4}-\d{2}-\d{2})`)
	_, ok := e.Extract("not a timestamp line")
	if ok {
		t.Fatal("expected no match")
	}
}

func TestExtract_PartialGroups(t *testing.T) {
	// Optional group may not participate in match.
	e, _ := New(`(?P<level>ERROR|WARN)(?:\s+(?P<code>\d+))?`)
	fields, ok := e.Extract("ERROR")
	if !ok {
		t.Fatal("expected match")
	}
	if fields["level"] != "ERROR" {
		t.Errorf("level: got %q", fields["level"])
	}
	// code group did not participate — value should be empty string
	if v, exists := fields["code"]; exists && v != "" {
		t.Errorf("code: expected empty, got %q", v)
	}
}

func TestFields_ReturnsCopy(t *testing.T) {
	e, _ := New(`(?P<a>x)(?P<b>y)`)
	f1 := e.Fields()
	f1[0] = "mutated"
	f2 := e.Fields()
	if f2[0] == "mutated" {
		t.Error("Fields() should return a copy, not the internal slice")
	}
}
