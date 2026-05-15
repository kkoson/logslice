package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew_EmptyPath(t *testing.T) {
	_, err := New("")
	if err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestNew_MissingFile_ReturnsZeroState(t *testing.T) {
	cp, err := New(filepath.Join(t.TempDir(), "cp.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s := cp.Get(); s.Offset != 0 || s.Path != "" {
		t.Fatalf("expected zero state, got %+v", s)
	}
}

func TestNew_CorruptFile_ReturnsError(t *testing.T) {
	f := filepath.Join(t.TempDir(), "cp.json")
	if err := os.WriteFile(f, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := New(f)
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestSave_PersistsState(t *testing.T) {
	dir := t.TempDir()
	cp, _ := New(filepath.Join(dir, "cp.json"))
	if err := cp.Save("/var/log/app.log", 1024); err != nil {
		t.Fatalf("Save: %v", err)
	}
	s := cp.Get()
	if s.Path != "/var/log/app.log" || s.Offset != 1024 {
		t.Fatalf("unexpected state: %+v", s)
	}
}

func TestSave_CanBeReloaded(t *testing.T) {
	f := filepath.Join(t.TempDir(), "cp.json")
	cp, _ := New(f)
	cp.Save("/logs/out.log", 4096) //nolint:errcheck

	cp2, err := New(f)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	s := cp2.Get()
	if s.Path != "/logs/out.log" || s.Offset != 4096 {
		t.Fatalf("reloaded state mismatch: %+v", s)
	}
}

func TestReset_ClearsStateAndFile(t *testing.T) {
	f := filepath.Join(t.TempDir(), "cp.json")
	cp, _ := New(f)
	cp.Save("/logs/out.log", 512) //nolint:errcheck

	if err := cp.Reset(); err != nil {
		t.Fatalf("Reset: %v", err)
	}
	if s := cp.Get(); s.Offset != 0 || s.Path != "" {
		t.Fatalf("expected zero state after reset, got %+v", s)
	}
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Fatal("expected checkpoint file to be removed")
	}
}

func TestReset_NoFile_IsNoop(t *testing.T) {
	cp, _ := New(filepath.Join(t.TempDir(), "cp.json"))
	if err := cp.Reset(); err != nil {
		t.Fatalf("Reset on missing file: %v", err)
	}
}
