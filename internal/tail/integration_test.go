package tail_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/tail"
)

// TestTailer_HighVolume writes many lines and verifies all are received.
func TestTailer_HighVolume(t *testing.T) {
	const lineCount = 200

	dir := t.TempDir()
	path := filepath.Join(dir, "volume.log")

	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	defer f.Close()

	tr, err := tail.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tr.Start(ctx)
	time.Sleep(50 * time.Millisecond)

	// Write all lines.
	for i := 0; i < lineCount; i++ {
		if _, err := fmt.Fprintf(f, "line %d\n", i); err != nil {
			t.Fatalf("write line %d: %v", i, err)
		}
	}

	received := make([]string, 0, lineCount)
	timeout := time.After(8 * time.Second)
	for len(received) < lineCount {
		select {
		case line, ok := <-tr.Lines():
			if !ok {
				t.Fatalf("channel closed after %d lines", len(received))
			}
			received = append(received, line)
		case <-timeout:
			t.Fatalf("timeout: only received %d/%d lines", len(received), lineCount)
		}
	}

	for i, line := range received {
		want := fmt.Sprintf("line %d", i)
		if line != want {
			t.Errorf("line[%d] = %q, want %q", i, line, want)
		}
	}
}
