package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/logslice/internal/aggregator"
	"github.com/yourorg/logslice/internal/source"
)

// runPivot is invoked when --pivot is provided on the CLI.
// Usage: logslice --pivot key=<field>,val=<field>,op=<op>
func runPivot(args []string) error {
	fs := flag.NewFlagSet("pivot", flag.ContinueOnError)
	key := fs.String("key", "", "field to group by")
	val := fs.String("val", "", "field to aggregate (not required for count)")
	op := fs.String("op", "count", "aggregation op: count|sum|avg|min|max")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *key == "" {
		return fmt.Errorf("--pivot requires -key")
	}

	p, err := aggregator.NewPivot(*key, *val, *op)
	if err != nil {
		return err
	}

	src := source.NewStdinSource()
	lines, errc := src.Lines()

	for line := range lines {
		p.Record(line)
	}
	if err := <-errc; err != nil {
		return fmt.Errorf("source error: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	for _, row := range p.Snapshot() {
		if err := enc.Encode(row); err != nil {
			return err
		}
	}
	return nil
}
