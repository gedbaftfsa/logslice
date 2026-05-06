// Package tail provides functionality to follow a file as it grows,
// similar to `tail -f`, emitting new lines on a channel.
package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Options configures the tail behaviour.
type Options struct {
	// PollInterval is how often to check for new data when the reader is
	// exhausted. Defaults to 200ms if zero.
	PollInterval time.Duration
}

// Follow reads lines from path, sending each non-empty line on the returned
// channel. It blocks until ctx is cancelled or an unrecoverable read error
// occurs, after which the channel is closed.
func Follow(ctx context.Context, path string, opts Options) (<-chan string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	if opts.PollInterval == 0 {
		opts.PollInterval = 200 * time.Millisecond
	}

	lines := make(chan string, 64)
	go func() {
		defer close(lines)
		defer f.Close()
		r := bufio.NewReader(f)
		for {
			line, err := r.ReadString('\n')
			if len(line) > 0 {
				// Strip trailing newline.
				if line[len(line)-1] == '\n' {
					line = line[:len(line)-1]
				}
				if line != "" {
					select {
					case lines <- line:
					case <-ctx.Done():
						return
					}
				}
			}
			if err != nil {
				if err != io.EOF {
					return
				}
				// EOF — wait for more data or cancellation.
				select {
				case <-ctx.Done():
					return
				case <-time.After(opts.PollInterval):
				}
			}
		}
	}()

	return lines, nil
}
