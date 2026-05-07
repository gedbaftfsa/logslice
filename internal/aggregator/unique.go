package aggregator

import (
	"encoding/json"
	"sync"
)

// Unique tracks the count of distinct values for a given field across log entries.
type Unique struct {
	mu    sync.Mutex
	field string
	seen  map[string]struct{}
}

// NewUnique creates a new Unique tracker for the specified field.
func NewUnique(field string) *Unique {
	return &Unique{
		field: field,
		seen:  make(map[string]struct{}),
	}
}

// Record processes a raw JSON log entry and records the value of the tracked field.
func (u *Unique) Record(raw []byte) {
	var entry map[string]interface{}
	if err := json.Unmarshal(raw, &entry); err != nil {
		return
	}
	v, ok := entry[u.field]
	if !ok {
		return
	}
	var key string
	switch val := v.(type) {
	case string:
		key = val
	default:
		b, err := json.Marshal(val)
		if err != nil {
			return
		}
		key = string(b)
	}
	u.mu.Lock()
	u.seen[key] = struct{}{}
	u.mu.Unlock()
}

// Count returns the number of distinct values seen so far.
func (u *Unique) Count() int {
	u.mu.Lock()
	defer u.mu.Unlock()
	return len(u.seen)
}

// Values returns a snapshot of all distinct values seen.
func (u *Unique) Values() []string {
	u.mu.Lock()
	defer u.mu.Unlock()
	out := make([]string, 0, len(u.seen))
	for k := range u.seen {
		out = append(out, k)
	}
	return out
}

// Reset clears all tracked values.
func (u *Unique) Reset() {
	u.mu.Lock()
	u.seen = make(map[string]struct{})
	u.mu.Unlock()
}

// MarshalJSON serialises the unique tracker state.
func (u *Unique) MarshalJSON() ([]byte, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	values := make([]string, 0, len(u.seen))
	for k := range u.seen {
		values = append(values, k)
	}
	return json.Marshal(map[string]interface{}{
		"field":  u.field,
		"count":  len(u.seen),
		"values": values,
	})
}
