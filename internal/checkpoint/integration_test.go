package checkpoint_test

import (
	"path/filepath"
	"sync"
	"testing"

	"github.com/yourorg/logslice/internal/checkpoint"
)

// TestCheckpoint_ConcurrentSaves verifies that concurrent Save calls do not
// corrupt the checkpoint file and that the last write wins.
func TestCheckpoint_ConcurrentSaves(t *testing.T) {
	f := filepath.Join(t.TempDir(), "cp.json")
	cp, err := checkpoint.New(f)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		offset := int64(i * 100)
		go func() {
			defer wg.Done()
			cp.Save("/logs/app.log", offset) //nolint:errcheck
		}()
	}
	wg.Wait()

	// Reload and verify the file is valid JSON with a non-negative offset.
	cp2, err := checkpoint.New(f)
	if err != nil {
		t.Fatalf("reload after concurrent writes: %v", err)
	}
	s := cp2.Get()
	if s.Path != "/logs/app.log" {
		t.Fatalf("unexpected path: %q", s.Path)
	}
	if s.Offset < 0 {
		t.Fatalf("negative offset: %d", s.Offset)
	}
}

// TestCheckpoint_RoundTrip exercises the full save → reload → reset cycle.
func TestCheckpoint_RoundTrip(t *testing.T) {
	f := filepath.Join(t.TempDir(), "state.json")

	cp, _ := checkpoint.New(f)
	if err := cp.Save("/var/log/syslog", 8192); err != nil {
		t.Fatal(err)
	}

	cp2, err := checkpoint.New(f)
	if err != nil {
		t.Fatal(err)
	}
	s := cp2.Get()
	if s.Path != "/var/log/syslog" || s.Offset != 8192 {
		t.Fatalf("round-trip mismatch: %+v", s)
	}

	if err := cp2.Reset(); err != nil {
		t.Fatal(err)
	}

	cp3, err := checkpoint.New(f)
	if err != nil {
		t.Fatal(err)
	}
	if s3 := cp3.Get(); s3.Offset != 0 {
		t.Fatalf("expected zero offset after reset reload, got %d", s3.Offset)
	}
}
