package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/user/logslice/internal/aggregator"
	"github.com/user/logslice/internal/source"
)

// runSample reads from stdin and prints a reservoir sample of N log entries.
func runSample(args []string) {
	fs := flag.NewFlagSet("sample", flag.ExitOnError)
	n := fs.Int("n", 10, "reservoir sample size")
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "sample: %v\n", err)
		os.Exit(1)
	}

	src := source.NewStdinSource()
	sampler := aggregator.NewSample(*n)

	for line := range src {
		sampler.Record(line)
	}

	for _, entry := range sampler.Snapshot() {
		b, err := json.Marshal(entry)
		if err != nil {
			continue
		}
		fmt.Println(string(b))
	}
}
