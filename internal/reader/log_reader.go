// Package reader provides utilities for reading log input from various sources.
package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Source represents a log input source.
type Source interface {
	Lines() (<-chan string, <-chan error)
	Close() error
}

// FileSource reads log lines from a file.
type FileSource struct {
	path string
	file *os.File
}

// NewFileSource opens the given file path for reading.
func NewFileSource(path string) (*FileSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("reader: open file %q: %w", path, err)
	}
	return &FileSource{path: path, file: f}, nil
}

// Lines returns a channel of lines and a channel of errors.
// The lines channel is closed when EOF is reached or an error occurs.
func (fs *FileSource) Lines() (<-chan string, <-chan error) {
	return linesFromReader(fs.file, fs.path)
}

// Close closes the underlying file.
func (fs *FileSource) Close() error {
	return fs.file.Close()
}

// StdinSource reads log lines from standard input.
type StdinSource struct{}

// NewStdinSource creates a source that reads from os.Stdin.
func NewStdinSource() *StdinSource {
	return &StdinSource{}
}

// Lines returns a channel of lines read from stdin.
func (s *StdinSource) Lines() (<-chan string, <-chan error) {
	return linesFromReader(os.Stdin, "stdin")
}

// Close is a no-op for stdin.
func (s *StdinSource) Close() error { return nil }

// linesFromReader reads lines from r, labeling any scan errors with name.
func linesFromReader(r io.Reader, name string) (<-chan string, <-chan error) {
	lines := make(chan string)
	errs := make(chan error, 1)
	go func() {
		defer close(lines)
		defer close(errs)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		if err := scanner.Err(); err != nil {
			errs <- fmt.Errorf("reader: scanning %s: %w", name, err)
		}
	}()
	return lines, errs
}
