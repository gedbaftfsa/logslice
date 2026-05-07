package aggregator

import (
	"encoding/json"
	"sync"
)

// GroupBy aggregates counts of log entries grouped by the value of a field.
type GroupBy struct {
	mu     sync.Mutex
	field  string
	counts map[string]int
}

// NewGroupBy creates a new GroupBy aggregator for the given field name.
func NewGroupBy(field string) *GroupBy {
	return &GroupBy{
		field:  field,
		counts: make(map[string]int),
	}
}

// Record processes a raw JSON log line and increments the count for the
// value of the configured field. Invalid JSON or missing fields are skipped.
func (g *GroupBy) Record(line []byte) {
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return
	}
	v, ok := entry[g.field]
	if !ok {
		return
	}
	key, ok := v.(string)
	if !ok {
		return
	}
	g.mu.Lock()
	g.counts[key]++
	g.mu.Unlock()
}

// Snapshot returns a copy of the current group counts.
func (g *GroupBy) Snapshot() map[string]int {
	g.mu.Lock()
	defer g.mu.Unlock()
	out := make(map[string]int, len(g.counts))
	for k, v := range g.counts {
		out[k] = v
	}
	return out
}

// Reset clears all accumulated counts.
func (g *GroupBy) Reset() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counts = make(map[string]int)
}

// MarshalJSON serialises the snapshot as a JSON object.
func (g *GroupBy) MarshalJSON() ([]byte, error) {
	return json.Marshal(g.Snapshot())
}
