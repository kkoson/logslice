package transform

import (
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("", "x")
	if err == nil {
		t.Fatal("expected error for empty pattern")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("[invalid", "x")
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}

func TestNew_Valid(t *testing.T) {
	tr, err := New(`\d+`, "NUM")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr == nil {
		t.Fatal("expected non-nil Transformer")
	}
}

func TestApply_NoMatch_ReturnsOriginal(t *testing.T) {
	tr := MustNew(`\d+`, "NUM")
	line := "no digits here"
	if got := tr.Apply(line); got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestApply_ReplacesMatch(t *testing.T) {
	tr := MustNew(`\d+`, "NUM")
	got := tr.Apply("error code 404 on line 12")
	want := "error code NUM on line NUM"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestApply_NamedGroup(t *testing.T) {
	tr := MustNew(`(?P<level>INFO|WARN|ERROR)`, "[${level}]")
	got := tr.Apply("2024-01-01 INFO something happened")
	want := "2024-01-01 [INFO] something happened"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestChain_AppliesInOrder(t *testing.T) {
	t1 := MustNew(`\d+`, "NUM")
	t2 := MustNew(`NUM`, "<redacted>")
	got := Chain([]*Transformer{t1, t2}, "code 42")
	want := "code <redacted>"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestChain_Empty_ReturnsOriginal(t *testing.T) {
	line := "unchanged"
	got := Chain(nil, line)
	if got != line {
		t.Fatalf("expected %q, got %q", line, got)
	}
}

func TestMustNew_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for invalid pattern")
		}
	}()
	MustNew("[bad", "x")
}
