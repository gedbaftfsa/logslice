package aggregator

import (
	"encoding/json"
	"sync"
)

// Counter tracks per-field value frequencies across log entries.
type Counter struct {
	mu     sync.Mutex
	field  string
	counts map[string]int64
}

// NewCounter creates a Counter that tallies occurrences of distinct values for the given field.
func NewCounter(field string) *Counter {
	return &Counter{
		field:  field,
		counts: make(map[string]int64),
	}
}

// Record extracts the target field from a raw JSON entry and increments its count.
// Entries missing the field or with non-string values are silently skipped.
func (c *Counter) Record(entry []byte) {
	var obj map[string]interface{}
	if err := json.Unmarshal(entry, &obj); err != nil {
		return
	}
	v, ok := obj[c.field]
	if !ok {
		return
	}
	var key string
	switch val := v.(type) {
	case string:
		key = val
	case float64:
		key = string(json.RawMessage(entry)) // fallback: use raw representation
		_ = val
		return // skip numeric values for now
	default:
		return
	}
	c.mu.Lock()
	c.counts[key]++
	c.mu.Unlock()
}

// Snapshot returns a copy of the current counts.
func (c *Counter) Snapshot() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]int64, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}

// Reset clears all counts.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts = make(map[string]int64)
}

// MarshalJSON serialises the counter as {"field":"...","counts":{...}}.
func (c *Counter) MarshalJSON() ([]byte, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return json.Marshal(struct {
		Field  string           `json:"field"`
		Counts map[string]int64 `json:"counts"`
	}{
		Field:  c.field,
		Counts: c.counts,
	})
}
