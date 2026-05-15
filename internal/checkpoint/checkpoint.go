package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// State holds the persisted position within a log file.
type State struct {
	Path   string `json:"path"`
	Offset int64  `json:"offset"`
}

// Checkpoint persists and restores read positions so that logslice can
// resume from where it left off after a restart.
type Checkpoint struct {
	mu       sync.Mutex
	filePath string
	state    State
}

// New loads an existing checkpoint file or creates a fresh one for the
// given state file path. Returns an error if the file exists but cannot
// be parsed.
func New(filePath string) (*Checkpoint, error) {
	if filePath == "" {
		return nil, fmt.Errorf("checkpoint: file path must not be empty")
	}
	cp := &Checkpoint{filePath: filePath}
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return cp, nil
		}
		return nil, fmt.Errorf("checkpoint: read %s: %w", filePath, err)
	}
	if err := json.Unmarshal(data, &cp.state); err != nil {
		return nil, fmt.Errorf("checkpoint: parse %s: %w", filePath, err)
	}
	return cp, nil
}

// Save atomically persists the current offset for the given log path.
func (c *Checkpoint) Save(logPath string, offset int64) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = State{Path: logPath, Offset: offset}
	data, err := json.Marshal(&c.state)
	if err != nil {
		return fmt.Errorf("checkpoint: marshal: %w", err)
	}
	tmp := c.filePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("checkpoint: write tmp: %w", err)
	}
	if err := os.Rename(tmp, c.filePath); err != nil {
		return fmt.Errorf("checkpoint: rename: %w", err)
	}
	return nil
}

// Get returns the most recently saved state.
func (c *Checkpoint) Get() State {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.state
}

// Reset clears the in-memory state and removes the checkpoint file if it
// exists.
func (c *Checkpoint) Reset() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.state = State{}
	if err := os.Remove(c.filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("checkpoint: remove: %w", err)
	}
	return nil
}
