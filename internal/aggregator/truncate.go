package aggregator

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Truncate trims string field values to a maximum length.
type Truncate struct {
	mu     sync.Mutex
	field  string
	maxLen int
	results []map[string]interface{}
}

// NewTruncate creates a Truncate aggregator that trims the given field to maxLen characters.
func NewTruncate(field string, maxLen int) (*Truncate, error) {
	if field == "" {
		return nil, fmt.Errorf("truncate: field name must not be empty")
	}
	if maxLen <= 0 {
		return nil, fmt.Errorf("truncate: maxLen must be greater than zero")
	}
	return &Truncate{field: field, maxLen: maxLen}, nil
}

// Record processes a JSON log entry, truncating the target field if it is a string.
func (t *Truncate) Record(line []byte) error {
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return fmt.Errorf("truncate: invalid JSON: %w", err)
	}

	if val, ok := entry[t.field]; ok {
		if s, ok := val.(string); ok {
			if len(s) > t.maxLen {
				entry[t.field] = s[:t.maxLen]
			}
		}
	}

	t.mu.Lock()
	t.results = append(t.results, entry)
	t.mu.Unlock()
	return nil
}

// Snapshot returns all processed entries as JSON lines.
func (t *Truncate) Snapshot() [][]byte {
	t.mu.Lock()
	defer t.mu.Unlock()

	out := make([][]byte, 0, len(t.results))
	for _, entry := range t.results {
		b, err := json.Marshal(entry)
		if err != nil {
			continue
		}
		out = append(out, b)
	}
	return out
}

// Reset clears all stored results.
func (t *Truncate) Reset() {
	t.mu.Lock()
	t.results = nil
	t.mu.Unlock()
}
