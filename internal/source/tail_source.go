package source

import (
	"context"
	"time"

	"github.com/yourusername/logslice/internal/tail"
)

// NewTailSource returns a Source that follows path as it grows, emitting each
// new JSON line. The source stops when ctx is cancelled.
func NewTailSource(ctx context.Context, path string) (Source, error) {
	lines, err := tail.Follow(ctx, path, tail.Options{
		PollInterval: 100 * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}
	return &chanSource{lines: lines}, nil
}

// chanSource adapts a string channel to the Source interface.
type chanSource struct {
	lines <-chan string
	current string
	done bool
}

func (s *chanSource) Next() bool {
	if s.done {
		return false
	}
	line, ok := <-s.lines
	if !ok {
		s.done = true
		return false
	}
	s.current = line
	return true
}

func (s *chanSource) Line() string {
	return s.current
}

func (s *chanSource) Err() error {
	return nil
}
