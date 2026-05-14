package multiline

import (
	"testing"
)

func TestNew_EmptyPattern(t *testing.T) {
	_, err := New("", "\n")
	if err != ErrEmptyPattern {
		t.Fatalf("expected ErrEmptyPattern, got %v", err)
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New("[invalid", "\n")
	if err == nil {
		t.Fatal("expected error for invalid regex, got nil")
	}
}

func TestNew_Valid(t *testing.T) {
	f, err := New(`^\d{4}-\d{2}-\d{2}`, "\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Folder")
	}
}

func TestAdd_SingleLineEvents(t *testing.T) {
	f, _ := New(`^START`, " ")

	_, ok := f.Add("START first")
	if ok {
		t.Fatal("first line should not emit an event")
	}

	event, ok := f.Add("START second")
	if !ok {
		t.Fatal("second start line should flush previous event")
	}
	if event != "START first" {
		t.Fatalf("unexpected event: %q", event)
	}
}

func TestAdd_ContinuationLines(t *testing.T) {
	f, _ := New(`^\d`, "\n")

	f.Add("1 main event")
	f.Add("  continuation a")
	f.Add("  continuation b")

	event, ok := f.Add("2 next event")
	if !ok {
		t.Fatal("expected event to be emitted")
	}
	want := "1 main event\n  continuation a\n  continuation b"
	if event != want {
		t.Fatalf("got %q, want %q", event, want)
	}
}

func TestFlush_EmptyBuffer(t *testing.T) {
	f, _ := New(`^START`, "\n")
	_, ok := f.Flush()
	if ok {
		t.Fatal("flush on empty buffer should return ok=false")
	}
}

func TestFlush_ReturnsPending(t *testing.T) {
	f, _ := New(`^START`, "\n")
	f.Add("START line")
	f.Add("  extra")

	event, ok := f.Flush()
	if !ok {
		t.Fatal("expected flush to return ok=true")
	}
	if event != "START line\n  extra" {
		t.Fatalf("unexpected event: %q", event)
	}

	// second flush should be empty
	_, ok = f.Flush()
	if ok {
		t.Fatal("second flush should return ok=false")
	}
}
