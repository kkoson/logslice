package template

import (
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("", "{{.level}}")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("[invalid", "{{.level}}")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_NoNamedGroups(t *testing.T) {
	_, err := New(`(\d+)`, "{{.level}}")
	if err == nil {
		t.Fatal("expected error when no named groups present")
	}
}

func TestNew_InvalidTemplate(t *testing.T) {
	_, err := New(`(?P<level>\w+)`, "{{.level")
	if err == nil {
		t.Fatal("expected error for invalid template")
	}
}

func TestNew_Valid(t *testing.T) {
	r, err := New(`(?P<level>\w+)\s+(?P<msg>.+)`, "[{{.level}}] {{.msg}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil renderer")
	}
}

func TestApply_NoMatch_ReturnsOriginal(t *testing.T) {
	r, _ := New(`(?P<level>\w+)\s+(?P<msg>.+)`, "[{{.level}}] {{.msg}}")
	got, err := r.Apply("no match here!!!")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "no match here!!!" {
		t.Errorf("expected original line, got %q", got)
	}
}

func TestApply_ReformatsLine(t *testing.T) {
	r, _ := New(`(?P<level>\w+)\s+(?P<msg>.+)`, "[{{.level}}] {{.msg}}")
	got, err := r.Apply("ERROR something went wrong")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "[ERROR] something went wrong"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApply_SingleGroup(t *testing.T) {
	r, _ := New(`level=(?P<level>\w+)`, "severity:{{.level}}")
	got, err := r.Apply("level=WARN extra stuff")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "severity:WARN"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestApply_MissingGroupInTemplate_ReturnsEmpty(t *testing.T) {
	r, _ := New(`(?P<level>\w+)`, "{{.level}} {{.missing}}")
	// text/template renders missing map keys as empty string by default
	got, err := r.Apply("INFO")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "INFO " {
		t.Errorf("got %q, want %q", got, "INFO ")
	}
}
