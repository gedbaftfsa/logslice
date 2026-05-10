package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/logslice/internal/aggregator"
)

// runGroupBy reads JSON log lines from r, groups them by the given field,
// and writes the resulting counts as JSON to w.
// Each line must be a valid JSON object; malformed lines are skipped with a
// warning written to stderr.
func runGroupBy(args []string, r io.Reader, w io.Writer) error {
	fs := flag.NewFlagSet("groupby", flag.ContinueOnError)
	field := fs.String("field", "", "field name to group by (required)")
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *field == "" {
		return fmt.Errorf("groupby: -field is required")
	}

	g := aggregator.NewGroupBy(*field)
	dec := json.NewDecoder(r)
	var skipped int
	for dec.More() {
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			skipped++
			fmt.Fprintf(os.Stderr, "groupby: skipping malformed line: %v\n", err)
			continue
		}
		g.Record([]byte(raw))
	}
	if skipped > 0 {
		fmt.Fprintf(os.Stderr, "groupby: skipped %d malformed line(s)\n", skipped)
	}

	b, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return fmt.Errorf("groupby: marshal: %w", err)
	}
	_, err = fmt.Fprintf(w, "%s\n", b)
	return err
}
