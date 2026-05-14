package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/logslice/internal/aggregator"
)

// runLimit reads newline-delimited JSON from lines, keeps at most n entries,
// and writes each retained entry as a JSON line to stdout.
func runLimit(lines <-chan []byte, n int) error {
	lim, err := aggregator.NewLimit(n)
	if err != nil {
		return fmt.Errorf("limit: %w", err)
	}

	for line := range lines {
		more := lim.Record(line)
		if !more {
			break
		}
	}

	enc := json.NewEncoder(os.Stdout)
	for _, entry := range lim.Entries() {
		if err := enc.Encode(entry); err != nil {
			return fmt.Errorf("limit: encode: %w", err)
		}
	}
	return nil
}
