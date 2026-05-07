package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/user/logslice/internal/aggregator"
)

// runTimeBucket reads stdin line-by-line and prints time-bucketed counts.
func runTimeBucket(args []string) {
	fs := flag.NewFlagSet("timebucket", flag.ExitOnError)
	field := fs.String("field", "ts", "JSON field containing the timestamp (RFC3339 or Unix seconds)")
	intervalStr := fs.String("interval", "1m", "bucket interval (e.g. 1m, 5m, 1h)")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "timebucket: %v\n", err)
		os.Exit(1)
	}

	interval, err := time.ParseDuration(*intervalStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "timebucket: invalid interval %q: %v\n", *intervalStr, err)
		os.Exit(1)
	}

	tb := aggregator.NewTimeBucket(*field, interval)

	dec := json.NewDecoder(os.Stdin)
	for dec.More() {
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			continue
		}
		tb.Record(raw)
	}

	out, err := json.MarshalIndent(tb, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "timebucket: marshal error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
