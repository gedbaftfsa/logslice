package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/aggregator"
	"github.com/yourorg/logslice/internal/source"
)

// runSort reads JSON lines from stdin, sorts by field, and writes to stdout.
func runSort(args []string) error {
	fs := flag.NewFlagSet("sort", flag.ContinueOnError)
	field := fs.String("field", "", "JSON field to sort by (required)")
	order := fs.String("order", "asc", "Sort order: asc or desc")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *field == "" {
		fs.Usage()
		return fmt.Errorf("sort: -field is required")
	}

	s, err := aggregator.NewSort(*field, *order)
	if err != nil {
		return err
	}

	src := source.NewStdinSource()
	for line := range src.Lines() {
		if err := s.Record(line); err != nil {
			fmt.Fprintf(os.Stderr, "sort: skipping invalid line: %v\n", err)
		}
	}

	sorted, err := s.Flush()
	if err != nil {
		return fmt.Errorf("sort: flush failed: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	for _, b := range sorted {
		var m map[string]interface{}
		if err := json.Unmarshal(b, &m); err != nil {
			continue
		}
		if err := enc.Encode(m); err != nil {
			return err
		}
	}
	return nil
}
