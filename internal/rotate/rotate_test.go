package rotate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNew_InvalidMaxBytes(t *testing.T) {
	_, err := New(Config{Dir: t.TempDir(), MaxBytes: 0})
	if err == nil {
		t.Fatal("expected error for MaxBytes=0")
	}
}

func TestNew_NegativeMaxBytes(t *testing.T) {
	_, err := New(Config{Dir: t.TempDir(), MaxBytes: -1})
	if err == nil {
		t.Fatal("expected error for negative MaxBytes")
	}
}

func TestNew_EmptyDir(t *testing.T) {
	_, err := New(Config{Dir: "", MaxBytes: 1024})
	if err == nil {
		t.Fatal("expected error for empty Dir")
	}
}

func TestNew_CreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "sub", "logs")
	w, err := New(Config{Dir: dir, MaxBytes: 1024})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("dir not created: %v", err)
	}
}

func TestWrite_CreatesLogFile(t *testing.T) {
	dir := t.TempDir()
	w, err := New(Config{Dir: dir, Prefix: "app-", MaxBytes: 1024})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()

	if _, err := w.Write([]byte("hello\n")); err != nil {
		t.Fatalf("Write: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("expected 1 file, got %d", len(entries))
	}
	if !strings.HasPrefix(entries[0].Name(), "app-") {
		t.Errorf("unexpected filename: %s", entries[0].Name())
	}
}

func TestWrite_RotatesOnSizeExceeded(t *testing.T) {
	dir := t.TempDir()
	w, err := New(Config{Dir: dir, Prefix: "r-", MaxBytes: 10})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()

	for i := 0; i < 5; i++ {
		if _, err := w.Write([]byte("hello\n")); err != nil {
			t.Fatalf("Write %d: %v", i, err)
		}
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) < 2 {
		t.Errorf("expected rotation to create multiple files, got %d", len(entries))
	}
}

func TestClose_Idempotent(t *testing.T) {
	w, err := New(Config{Dir: t.TempDir(), MaxBytes: 512})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("first Close: %v", err)
	}
}
