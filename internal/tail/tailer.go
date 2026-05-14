package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// PollInterval is how often the tailer checks for new data.
const PollInterval = 200 * time.Millisecond

// Tailer follows a file as new lines are appended, similar to `tail -f`.
type Tailer struct {
	path   string
	lines  chan string
	errors chan error
}

// New creates a Tailer for the given file path.
func New(path string) (*Tailer, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	return &Tailer{
		path:   path,
		lines:  make(chan string, 64),
		errors: make(chan error, 1),
	}, nil
}

// Lines returns the channel on which new log lines are delivered.
func (t *Tailer) Lines() <-chan string { return t.lines }

// Errors returns the channel on which read errors are delivered.
func (t *Tailer) Errors() <-chan error { return t.errors }

// Start begins tailing the file. It seeks to the end before watching for new
// content. Cancel ctx to stop.
func (t *Tailer) Start(ctx context.Context) {
	go t.run(ctx)
}

func (t *Tailer) run(ctx context.Context) {
	defer close(t.lines)

	f, err := os.Open(t.path)
	if err != nil {
		t.errors <- err
		return
	}
	defer f.Close()

	// Seek to end so we only see new lines.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		t.errors <- err
		return
	}

	reader := bufio.NewReader(f)
	ticker := time.NewTicker(PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for {
				line, err := reader.ReadString('\n')
				if len(line) > 0 {
					// Strip trailing newline before sending.
					if len(line) > 0 && line[len(line)-1] == '\n' {
						line = line[:len(line)-1]
					}
					select {
					case t.lines <- line:
					case <-ctx.Done():
						return
					}
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					t.errors <- err
					return
				}
			}
		}
	}
}
