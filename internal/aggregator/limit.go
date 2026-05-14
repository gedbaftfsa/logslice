package aggregator

import (
	"encoding/json"
	"fmt"
)

// Limit passes through at most N log entries, then signals done.
type Limit struct {
	max     int
	count   int
	entries []map[string]any
}

// NewLimit creates a Limit aggregator that retains at most n entries.
func NewLimit(n int) (*Limit, error) {
	if n <= 0 {
		return nil, fmt.Errorf("limit: n must be greater than zero, got %d", n)
	}
	return &Limit{max: n}, nil
}

// Record accepts a raw JSON log line and stores it if under the limit.
// Returns false once the limit is reached (caller may stop feeding).
func (l *Limit) Record(line []byte) bool {
	if l.count >= l.max {
		return false
	}
	var entry map[string]any
	if err := json.Unmarshal(line, &entry); err != nil {
		return l.count < l.max
	}
	l.entries = append(l.entries, entry)
	l.count++
	return l.count < l.max
}

// Done reports whether the limit has been reached.
func (l *Limit) Done() bool {
	return l.count >= l.max
}

// Entries returns the collected entries.
func (l *Limit) Entries() []map[string]any {
	return l.entries
}

// Reset clears all collected entries and resets the counter.
func (l *Limit) Reset() {
	l.count = 0
	l.entries = nil
}

// Snapshot returns a JSON-encoded summary of the limit state.
func (l *Limit) Snapshot() ([]byte, error) {
	type snapshot struct {
		Max     int              `json:"max"`
		Count   int              `json:"count"`
		Entries []map[string]any `json:"entries"`
	}
	return json.Marshal(snapshot{
		Max:     l.max,
		Count:   l.count,
		Entries: l.entries,
	})
}
