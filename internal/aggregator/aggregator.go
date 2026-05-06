// Package aggregator merges log entries from multiple sources into a single
// ordered channel, tagging each entry with its origin label.
package aggregator

import (
	"encoding/json"
	"sync"
)

// Entry is a single log record enriched with its source label.
type Entry struct {
	Source string
	Fields map[string]interface{}
	Raw    string
}

// Source is anything that emits raw JSON log lines.
type Source interface {
	Lines() <-chan string
}

// Merge fans-in lines from all labelled sources, parses each line as JSON,
// and sends the resulting Entry values on the returned channel.
// The channel is closed once every source is exhausted.
func Merge(sources map[string]Source) <-chan Entry {
	out := make(chan Entry, 64)
	var wg sync.WaitGroup

	for label, src := range sources {
		wg.Add(1)
		go func(label string, src Source) {
			defer wg.Done()
			for line := range src.Lines() {
				var fields map[string]interface{}
				if err := json.Unmarshal([]byte(line), &fields); err != nil {
					// Skip non-JSON lines silently.
					continue
				}
				out <- Entry{
					Source: label,
					Fields: fields,
					Raw:    line,
				}
			}
		}(label, src)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
