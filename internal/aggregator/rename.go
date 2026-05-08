package aggregator

import (
	"encoding/json"
	"sync"
)

// Rename renames fields in each JSON log entry.
// It applies a map of oldKey -> newKey transformations.
type Rename struct {
	mu      sync.Mutex
	mapping map[string]string
	results []map[string]interface{}
}

// NewRename creates a new Rename transformer with the given field mapping.
func NewRename(mapping map[string]string) *Rename {
	m := make(map[string]string, len(mapping))
	for k, v := range mapping {
		m[k] = v
	}
	return &Rename{mapping: m}
}

// Record processes a raw JSON line, renaming fields as configured.
func (r *Rename) Record(line []byte) error {
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return err
	}

	for oldKey, newKey := range r.mapping {
		if val, ok := entry[oldKey]; ok {
			entry[newKey] = val
			delete(entry, oldKey)
		}
	}

	r.mu.Lock()
	r.results = append(r.results, entry)
	r.mu.Unlock()
	return nil
}

// Results returns all transformed entries as JSON lines.
func (r *Rename) Results() [][]byte {
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([][]byte, 0, len(r.results))
	for _, entry := range r.results {
		b, err := json.Marshal(entry)
		if err == nil {
			out = append(out, b)
		}
	}
	return out
}

// Reset clears all stored results.
func (r *Rename) Reset() {
	r.mu.Lock()
	r.results = r.results[:0]
	r.mu.Unlock()
}

// AddMapping adds or updates a field rename rule.
// If oldKey already exists in the mapping, its target is overwritten.
func (r *Rename) AddMapping(oldKey, newKey string) {
	r.mu.Lock()
	r.mapping[oldKey] = newKey
	r.mu.Unlock()
}
