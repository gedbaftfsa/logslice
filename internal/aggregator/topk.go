package aggregator

import (
	"encoding/json"
	"sort"
	"sync"
)

// TopK tracks the top K most frequent values for a given field.
type TopK struct {
	mu    sync.Mutex
	field string
	k     int
	counts map[string]int
}

type TopKEntry struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

// NewTopK creates a new TopK tracker for the given field and k.
func NewTopK(field string, k int) *TopK {
	if k <= 0 {
		k = 10
	}
	return &TopK{
		field:  field,
		k:      k,
		counts: make(map[string]int),
	}
}

// Record parses a JSON log line and increments the count for the field value.
func (t *TopK) Record(line []byte) error {
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return err
	}
	v, ok := entry[t.field]
	if !ok {
		return nil
	}
	val, ok := v.(string)
	if !ok {
		return nil
	}
	t.mu.Lock()
	t.counts[val]++
	t.mu.Unlock()
	return nil
}

// Top returns the top K entries sorted by count descending.
func (t *TopK) Top() []TopKEntry {
	t.mu.Lock()
	defer t.mu.Unlock()

	entries := make([]TopKEntry, 0, len(t.counts))
	for val, cnt := range t.counts {
		entries = append(entries, TopKEntry{Value: val, Count: cnt})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Count != entries[j].Count {
			return entries[i].Count > entries[j].Count
		}
		return entries[i].Value < entries[j].Value
	})
	if len(entries) > t.k {
		entries = entries[:t.k]
	}
	return entries
}

// Reset clears all counts.
func (t *TopK) Reset() {
	t.mu.Lock()
	t.counts = make(map[string]int)
	t.mu.Unlock()
}
