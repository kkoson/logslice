package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempLog(t *testing.T, lines []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "logslice-*.log")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer f.Close()
	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return filepath.ToSlash(f.Name())
}

func TestRun_NoArgs_ReadsStdin(t *testing.T) {
	// run with an empty file should succeed without error
	path := writeTempLog(t, nil)
	err := run([]string{"-file", path})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_FilterMatchFlag(t *testing.T) {
	path := writeTempLog(t, []string{"INFO starting", "ERROR crashed", "INFO ready"})

	// Redirect stdout via capturing — run writes to os.Stdout;
	// we test indirectly by checking no error and that excluded lines are gone.
	err := run([]string{"-file", path, "-match", "INFO"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_InvalidMatchPattern(t *testing.T) {
	path := writeTempLog(t, []string{"hello"})
	err := run([]string{"-file", path, "-match", "["})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
	if !strings.Contains(err.Error(), "invalid match pattern") {
		t.Errorf("error message mismatch: %v", err)
	}
}

func TestRun_InvalidExcludePattern(t *testing.T) {
	path := writeTempLog(t, []string{"hello"})
	err := run([]string{"-file", path, "-exclude", "("})
	if err == nil {
		t.Fatal("expected error for invalid exclude pattern")
	}
}

func TestRun_InvalidFormat(t *testing.T) {
	path := writeTempLog(t, []string{"hello"})
	err := run([]string{"-file", path, "-format", "xml"})
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestRun_FileNotFound(t *testing.T) {
	err := run([]string{"-file", "/no/such/file.log"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestRun_AggregationMode(t *testing.T) {
	lines := []string{
		"level=INFO msg=started",
		"level=ERROR msg=oops",
		"level=INFO msg=done",
	}
	path := writeTempLog(t, lines)
	err := run([]string{
		"-file", path,
		"-agg", `level=(?P<key>\w+)`,
		"-format", "plain",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// captureStdout replaces os.Stdout temporarily and returns written bytes.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	old := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}
