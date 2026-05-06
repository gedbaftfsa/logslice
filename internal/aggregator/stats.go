package aggregator

import (
	"encoding/json"
	"sync"
)

// Stats tracks basic metrics about log lines processed during aggregation.
type Stats struct {
	mu      sync.Mutex
	Total   int64            `json:"total"`
	BySource map[string]int64 `json:"by_source"`
	Invalid int64            `json:"invalid"`
}

// NewStats creates a new Stats instance.
func NewStats() *Stats {
	return &Stats{
		BySource: make(map[string]int64),
	}
}

// Record increments the total and per-source counters.
func (s *Stats) Record(source string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Total++
	s.BySource[source]++
}

// RecordInvalid increments the invalid line counter.
func (s *Stats) RecordInvalid() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Invalid++
}

// Snapshot returns a copy of the current stats as a JSON-serialisable map.
func (s *Stats) Snapshot() map[string]interface{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	byCopy := make(map[string]int64, len(s.BySource))
	for k, v := range s.BySource {
		byCopy[k] = v
	}
	return map[string]interface{}{
		"total":     s.Total,
		"by_source": byCopy,
		"invalid":   s.Invalid,
	}
}

// MarshalJSON implements json.Marshaler so Stats can be serialised directly.
func (s *Stats) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Snapshot())
}
