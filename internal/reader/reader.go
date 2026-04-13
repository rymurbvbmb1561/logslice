// Package reader provides utilities for reading log files line by line,
// supporting both regular files and stdin input.
package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// LineReader reads lines from an input source and sends them to a channel.
type LineReader struct {
	source io.ReadCloser
	name   string
}

// NewFileReader opens a file at the given path and returns a LineReader.
func NewFileReader(path string) (*LineReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("reader: opening file %q: %w", path, err)
	}
	return &LineReader{source: f, name: path}, nil
}

// NewStdinReader returns a LineReader that reads from standard input.
func NewStdinReader() *LineReader {
	return &LineReader{source: io.NopCloser(os.Stdin), name: "<stdin>"}
}

// Name returns the name of the underlying source (file path or "<stdin>").
func (r *LineReader) Name() string {
	return r.name
}

// Close closes the underlying source.
func (r *LineReader) Close() error {
	return r.source.Close()
}

// Lines reads all lines from the source and returns them as a slice.
// Empty lines are included so callers can decide how to handle them.
func (r *LineReader) Lines() ([]string, error) {
	scanner := bufio.NewScanner(r.source)
	// Allow lines up to 1 MiB to handle large JSON log entries.
	const maxCapacity = 1 << 20
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reader: scanning %q: %w", r.name, err)
	}
	return lines, nil
}

// Stream reads lines from the source and sends each to the provided channel.
// The channel is closed when reading is complete or an error occurs.
// Any scan error is returned after the channel is closed.
func (r *LineReader) Stream(ch chan<- string) error {
	defer close(ch)
	scanner := bufio.NewScanner(r.source)
	const maxCapacity = 1 << 20
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	for scanner.Scan() {
		ch <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reader: streaming %q: %w", r.name, err)
	}
	return nil
}
