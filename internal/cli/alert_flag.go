package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/aggregator"
	"github.com/user/logslice/internal/source"
)

// runAlert handles the --alert flag: "field:operator:threshold".
// Example: --alert latency_ms:gt:500
func runAlert(args []string) error {
	fs := flag.NewFlagSet("alert", flag.ContinueOnError)
	alertExpr := fs.String("alert", "", "field:operator:threshold (e.g. latency_ms:gt:500)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *alertExpr == "" {
		return fmt.Errorf("--alert requires field:operator:threshold")
	}

	parts := strings.SplitN(*alertExpr, ":", 3)
	if len(parts) != 3 {
		return fmt.Errorf("invalid --alert format, expected field:operator:threshold")
	}

	field := parts[0]
	op := parts[1]
	threshold, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return fmt.Errorf("invalid threshold %q: %w", parts[2], err)
	}

	validOps := map[string]bool{"gt": true, "lt": true, "eq": true}
	if !validOps[op] {
		return fmt.Errorf("unknown operator %q, use gt, lt, or eq", op)
	}

	cond := aggregator.AlertCondition{
		Field:     field,
		Operator:  op,
		Threshold: threshold,
	}
	a := aggregator.NewAlert(cond)

	src := source.NewStdinSource()
	for line := range src {
		a.Record(line)
	}

	var out map[string]interface{}
	json.Unmarshal(a.Snapshot(), &out)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		return fmt.Errorf("encode output: %w", err)
	}
	return nil
}
