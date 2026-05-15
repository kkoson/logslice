package rotate

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Writer is a log output writer that rotates files when they exceed MaxBytes
// or when a new day begins (if RotateDaily is set).
type Writer struct {
	mu          sync.Mutex
	dir         string
	prefix      string
	maxBytes    int64
	rotatDaily  bool
	file        *os.File
	written     int64
	currentDate string
}

// Config holds options for New.
type Config struct {
	Dir         string
	Prefix      string
	MaxBytes    int64
	RotateDaily bool
}

// New creates a Writer that writes to files under dir.
// MaxBytes must be > 0.
func New(cfg Config) (*Writer, error) {
	if cfg.MaxBytes <= 0 {
		return nil, fmt.Errorf("rotate: MaxBytes must be > 0, got %d", cfg.MaxBytes)
	}
	if cfg.Dir == "" {
		return nil, fmt.Errorf("rotate: Dir must not be empty")
	}
	if err := os.MkdirAll(cfg.Dir, 0o755); err != nil {
		return nil, fmt.Errorf("rotate: mkdir %s: %w", cfg.Dir, err)
	}
	w := &Writer{
		dir:        cfg.Dir,
		prefix:     cfg.Prefix,
		maxBytes:   cfg.MaxBytes,
		rotatDaily: cfg.RotateDaily,
	}
	if err := w.openNew(); err != nil {
		return nil, err
	}
	return w, nil
}

// Write implements io.Writer. It rotates the underlying file as needed.
func (w *Writer) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotateIfNeeded(int64(len(p))); err != nil {
		return 0, err
	}
	n, err := w.file.Write(p)
	w.written += int64(n)
	return n, err
}

// Close closes the underlying file.
func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func (w *Writer) rotateIfNeeded(incoming int64) error {
	newDate := time.Now().Format("2006-01-02")
	dailyRotate := w.rotatDaily && newDate != w.currentDate
	sizeRotate := w.written+incoming > w.maxBytes
	if dailyRotate || sizeRotate {
		if err := w.file.Close(); err != nil {
			return err
		}
		return w.openNew()
	}
	return nil
}

func (w *Writer) openNew() error {
	ts := time.Now().Format("20060102-150405")
	name := filepath.Join(w.dir, fmt.Sprintf("%s%s.log", w.prefix, ts))
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return fmt.Errorf("rotate: open %s: %w", name, err)
	}
	w.file = f
	w.written = 0
	w.currentDate = time.Now().Format("2006-01-02")
	return nil
}
