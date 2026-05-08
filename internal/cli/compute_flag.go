package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/user/logslice/internal/aggregator"
)

// runCompute parses --compute flags and runs the Compute aggregator over
// lines read from stdin, writing results to stdout.
//
// Flag format: --compute src:dst:op:operand
// Example:     --compute latency_ms:latency_s:div:1000
func runCompute(args []string) error {
	fs := flag.NewFlagSet("compute", flag.ContinueOnError)
	expr := fs.String("compute", "", "src:dst:op:operand")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *expr == "" {
		return fmt.Errorf("--compute requires src:dst:op:operand")
	}
	parts := strings.SplitN(*expr, ":", 4)
	if len(parts) != 4 {
		return fmt.Errorf("--compute: expected src:dst:op:operand, got %q", *expr)
	}
	operand, err := strconv.ParseFloat(parts[3], 64)
	if err != nil {
		return fmt.Errorf("--compute: invalid operand %q: %w", parts[3], err)
	}
	c, err := aggregator.NewCompute(parts[0], parts[1], parts[2], operand)
	if err != nil {
		return fmt.Errorf("--compute: %w", err)
	}
	dec := json.NewDecoder(os.Stdin)
	enc := json.NewEncoder(os.Stdout)
	for dec.More() {
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			continue
		}
		c.Record(raw)
	}
	for _, r := range c.Results() {
		if err := enc.Encode(r); err != nil {
			return err
		}
	}
	return nil
}
