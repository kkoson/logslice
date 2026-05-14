package tail_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/tail"
)

func writeLine(t *testing.T, f *os.File, line string) {
	t.Helper()
	if _, err := f.WriteString(line + "\n"); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestNew_FileNotFound(t *testing.T) {
	_, err := tail.New("/nonexistent/path/log.txt")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestTailer_ReceivesNewLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer f.Close()

	tr, err := tail.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tr.Start(ctx)

	// Give the tailer time to seek to end.
	time.Sleep(50 * time.Millisecond)

	writeLine(t, f, "hello world")
	writeLine(t, f, "second line")

	var got []string
	timeout := time.After(2 * time.Second)
	for len(got) < 2 {
		select {
		case line, ok := <-tr.Lines():
			if !ok {
				t.Fatal("lines channel closed early")
			}
			got = append(got, line)
		case <-timeout:
			t.Fatalf("timed out waiting for lines; got %v", got)
		}
	}

	if got[0] != "hello world" || got[1] != "second line" {
		t.Errorf("unexpected lines: %v", got)
	}
}

func TestTailer_StopsOnContextCancel(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "app.log")
	if _, err := os.Create(path); err != nil {
		t.Fatalf("create: %v", err)
	}

	tr, err := tail.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	tr.Start(ctx)

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case <-tr.Lines():
		// channel closed — OK
	case <-time.After(time.Second):
		t.Fatal("lines channel not closed after context cancel")
	}
}
