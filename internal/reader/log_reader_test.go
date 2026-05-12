package reader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/reader"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "test.log")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestNewFileSource_NotFound(t *testing.T) {
	_, err := reader.NewFileSource("/nonexistent/path/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestFileSource_LinesReadsAll(t *testing.T) {
	content := "line one\nline two\nline three\n"
	p := writeTempFile(t, content)

	src, err := reader.NewFileSource(p)
	if err != nil {
		t.Fatalf("NewFileSource: %v", err)
	}
	defer src.Close()

	lines, errs := src.Lines()
	var got []string
	for l := range lines {
		got = append(got, l)
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"line one", "line two", "line three"}
	if len(got) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(got))
	}
	for i, e := range expected {
		if got[i] != e {
			t.Errorf("line %d: want %q, got %q", i, e, got[i])
		}
	}
}

func TestFileSource_EmptyFile(t *testing.T) {
	p := writeTempFile(t, "")
	src, err := reader.NewFileSource(p)
	if err != nil {
		t.Fatalf("NewFileSource: %v", err)
	}
	defer src.Close()

	lines, errs := src.Lines()
	var count int
	for range lines {
		count++
	}
	if err := <-errs; err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 lines from empty file, got %d", count)
	}
}

func TestStdinSource_Close(t *testing.T) {
	src := reader.NewStdinSource()
	if err := src.Close(); err != nil {
		t.Errorf("Close() should be no-op, got: %v", err)
	}
}
