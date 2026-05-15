package rotate_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/logslice/internal/rotate"
)

func TestWriter_HighVolume_RotatesMultipleTimes(t *testing.T) {
	dir := t.TempDir()
	w, err := rotate.New(rotate.Config{
		Dir:      dir,
		Prefix:   "hv-",
		MaxBytes: 256,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()

	line := []byte("2024-01-01T00:00:00Z level=info msg=\"test log line\" id=12345\n")
	for i := 0; i < 100; i++ {
		if _, err := w.Write(line); err != nil {
			t.Fatalf("Write %d: %v", i, err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) < 5 {
		t.Errorf("expected at least 5 rotated files, got %d", len(entries))
	}
	for _, e := range entries {
		info, err := os.Stat(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Fatalf("Stat: %v", err)
		}
		if info.Size() == 0 {
			t.Errorf("file %s is empty", e.Name())
		}
	}
}

func TestWriter_ConcurrentWrites(t *testing.T) {
	dir := t.TempDir()
	w, err := rotate.New(rotate.Config{
		Dir:      dir,
		Prefix:   "conc-",
		MaxBytes: 512,
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer w.Close()

	done := make(chan struct{})
	for g := 0; g < 8; g++ {
		go func(id int) {
			for i := 0; i < 20; i++ {
				w.Write([]byte(fmt.Sprintf("goroutine=%d seq=%d\n", id, i))) //nolint:errcheck
			}
			done <- struct{}{}
		}(g)
	}
	for i := 0; i < 8; i++ {
		<-done
	}
}
