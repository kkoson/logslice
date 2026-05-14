package tail

import "context"

// Source wraps a Tailer so it satisfies the reader.Source interface expected
// by the rest of logslice (Next() string, bool and Close()).
type Source struct {
	tailer *Tailer
	cancel context.CancelFunc
	lines  <-chan string
}

// NewSource creates a Source that tails path. Call Close to stop tailing.
func NewSource(path string) (*Source, error) {
	tr, err := New(path)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	tr.Start(ctx)

	return &Source{
		tailer: tr,
		cancel: cancel,
		lines:  tr.Lines(),
	}, nil
}

// Next blocks until the next line is available. Returns (line, true) on
// success, or ("", false) when the source is closed or the file is gone.
func (s *Source) Next() (string, bool) {
	line, ok := <-s.lines
	return line, ok
}

// Close stops the underlying tailer and releases resources.
func (s *Source) Close() error {
	s.cancel()
	// Drain remaining lines so the goroutine can exit cleanly.
	for range s.lines {
	}
	return nil
}
