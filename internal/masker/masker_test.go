package masker

import (
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("", "***")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("[invalid", "***")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestNew_NoNamedGroups(t *testing.T) {
	_, err := New(`\d+`, "***")
	if err == nil {
		t.Fatal("expected error when no named groups present")
	}
}

func TestNew_Valid(t *testing.T) {
	m, err := New(`(?P<token>[A-Za-z0-9]+)`, "***")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("expected non-nil masker")
	}
}

func TestNew_DefaultMask(t *testing.T) {
	m, err := New(`(?P<token>[A-Za-z0-9]+)`, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.mask != "***" {
		t.Fatalf("expected default mask '***', got %q", m.mask)
	}
}

func TestApply_NoMatch_ReturnsOriginal(t *testing.T) {
	m, _ := New(`token=(?P<value>[A-Z]+)`, "[REDACTED]")
	input := "no token here"
	if got := m.Apply(input); got != input {
		t.Fatalf("expected %q, got %q", input, got)
	}
}

func TestApply_MasksNamedGroup(t *testing.T) {
	m, _ := New(`token=(?P<value>\S+)`, "[REDACTED]")
	got := m.Apply("auth token=abc123 ok")
	want := "auth token=[REDACTED] ok"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApply_MasksEmailLocalPart(t *testing.T) {
	m, _ := New(`(?P<local>[^@\s]+)@(?P<domain>[^\s]+)`, "***")
	got := m.Apply("user=alice@example.com login")
	// both local and domain groups are masked
	if got == "user=alice@example.com login" {
		t.Fatal("expected masking to occur")
	}
}

func TestApply_MultipleOccurrences(t *testing.T) {
	m, _ := New(`password=(?P<pw>\S+)`, "XXX")
	got := m.Apply("password=secret1 password=secret2")
	want := "password=XXX password=XXX"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
