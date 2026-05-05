// Package source provides interfaces and implementations for reading
// structured JSON log entries from various input sources such as files,
// stdin, and remote streams.
package source

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Entry represents a single parsed JSON log entry as a map of key-value pairs.
type Entry map[string]interface{}

// Source is the interface that wraps the basic Read method.
// Implementations return log entries one at a time.
type Source interface {
	// Read returns the next log entry. Returns io.EOF when no more entries
	// are available.
	Read() (Entry, error)
	// Close releases any resources held by the source.
	Close() error
	// Name returns a human-readable identifier for this source.
	Name() string
}

// ReaderSource reads JSON log entries line-by-line from an io.Reader.
type ReaderSource struct {
	name    string
	reader  io.ReadCloser
	scanner *bufio.Scanner
}

// NewReaderSource creates a Source that reads from the given io.ReadCloser.
// The name parameter is used to identify the source in output.
func NewReaderSource(name string, r io.ReadCloser) *ReaderSource {
	return &ReaderSource{
		name:    name,
		reader:  r,
		scanner: bufio.NewScanner(r),
	}
}

// NewFileSource opens the file at the given path and returns a Source
// that reads JSON log entries from it.
func NewFileSource(path string) (*ReaderSource, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("source: open file %q: %w", path, err)
	}
	return NewReaderSource(path, f), nil
}

// NewStdinSource returns a Source that reads JSON log entries from os.Stdin.
func NewStdinSource() *ReaderSource {
	return NewReaderSource("<stdin>", io.NopCloser(os.Stdin))
}

// Read scans the next non-empty line and attempts to parse it as a JSON object.
// Returns io.EOF when the underlying reader is exhausted.
func (s *ReaderSource) Read() (Entry, error) {
	for s.scanner.Scan() {
		line := s.scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var entry Entry
		if err := json.Unmarshal(line, &entry); err != nil {
			return nil, fmt.Errorf("source %s: parse JSON: %w", s.name, err)
		}
		return entry, nil
	}
	if err := s.scanner.Err(); err != nil {
		return nil, fmt.Errorf("source %s: scan: %w", s.name, err)
	}
	return nil, io.EOF
}

// Close closes the underlying reader.
func (s *ReaderSource) Close() error {
	return s.reader.Close()
}

// Name returns the identifier for this source.
func (s *ReaderSource) Name() string {
	return s.name
}

// MultiSource fans in multiple Sources, yielding entries in arrival order.
// Each entry is tagged with the originating source name under the key
// "_source" if that key is not already present in the entry.
type MultiSource struct {
	sources []Source
	current int
}

// NewMultiSource creates a MultiSource from the provided slice of Sources.
func NewMultiSource(sources []Source) *MultiSource {
	return &MultiSource{sources: sources}
}

// Read returns the next entry from the current source, advancing to the next
// source when the current one is exhausted. Returns io.EOF when all sources
// are exhausted.
func (m *MultiSource) Read() (Entry, error) {
	for m.current < len(m.sources) {
		entry, err := m.sources[m.current].Read()
		if err == io.EOF {
			m.current++
			continue
		}
		if err != nil {
			return nil, err
		}
		if _, ok := entry["_source"]; !ok {
			entry["_source"] = m.sources[m.current].Name()
		}
		return entry, nil
	}
	return nil, io.EOF
}

// Close closes all underlying sources, returning the first error encountered.
func (m *MultiSource) Close() error {
	var firstErr error
	for _, s := range m.sources {
		if err := s.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Name returns a combined name for the multi-source.
func (m *MultiSource) Name() string {
	return "<multi>"
}
