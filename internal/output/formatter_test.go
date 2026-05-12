package output

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNewFormatter_InvalidFormat(t *testing.T) {
	_, err := NewFormatter("xml", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestPlainFormatter_Write(t *testing.T) {
	var buf bytes.Buffer
	f, err := NewFormatter(FormatPlain, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := LogEntry{Line: "hello world", Index: 1}
	if err := f.Write(entry); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	got := strings.TrimSpace(buf.String())
	if got != "hello world" {
		t.Errorf("expected %q, got %q", "hello world", got)
	}
}

func TestJSONFormatter_Write(t *testing.T) {
	var buf bytes.Buffer
	f, err := NewFormatter(FormatJSON, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := LogEntry{Line: "error occurred", Index: 42}
	if err := f.Write(entry); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"line":"error occurred"`) {
		t.Errorf("expected JSON to contain line field, got: %s", out)
	}
	if !strings.Contains(out, `"index":42`) {
		t.Errorf("expected JSON to contain index field, got: %s", out)
	}
}

func TestCSVFormatter_Write(t *testing.T) {
	var buf bytes.Buffer
	f, err := NewFormatter(FormatCSV, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	entry := LogEntry{Line: "csv line", Index: 1, Timestamp: ts}
	if err := f.Write(entry); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines (header + row), got %d", len(lines))
	}
	if lines[0] != "index,timestamp,line" {
		t.Errorf("unexpected header: %q", lines[0])
	}
	if !strings.Contains(lines[1], "csv line") {
		t.Errorf("expected row to contain line content, got: %q", lines[1])
	}
}

func TestCSVFormatter_EscapesQuotes(t *testing.T) {
	var buf bytes.Buffer
	f, err := NewFormatter(FormatCSV, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entry := LogEntry{Line: `say "hello"`, Index: 1}
	if err := f.Write(entry); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"say ""hello"""
`) {
		t.Errorf("expected escaped quotes in CSV output, got: %s", out)
	}
}
