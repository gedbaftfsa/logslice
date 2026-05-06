package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/user/logslice/internal/aggregator"
)

// runCounter reads lines from stdin, counts values for the given field,
// and writes the result as JSON to stdout. It is invoked when --count is set.
func runCounter(field string) error {
	if field == "" {
		return fmt.Errorf("--count requires a field name")
	}

	counter := aggregator.NewCounter(field)

	decoder := json.NewDecoder(os.Stdin)
	for decoder.More() {
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			// skip malformed lines
			continue
		}
		counter.Record([]byte(raw))
	}

	data, err := json.MarshalIndent(counter, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal counter: %w", err)
	}
	fmt.Println(string(data))
	return nil
}
