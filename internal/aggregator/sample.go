package aggregator

import (
	"encoding/json"
	"math/rand"
	"sync"
)

// Sample retains a random reservoir sample of up to N log entries.
type Sample struct {
	mu      sync.Mutex
	size    int
	bucket  []map[string]any
	count   int
	rng     *rand.Rand
}

// NewSample creates a new reservoir sampler keeping at most size entries.
func NewSample(size int) *Sample {
	if size <= 0 {
		size = 10
	}
	return &Sample{
		size:   size,
		bucket: make([]map[string]any, 0, size),
		rng:    rand.New(rand.NewSource(42)),
	}
}

// Record processes a raw JSON log line using reservoir sampling.
func (s *Sample) Record(line []byte) {
	var entry map[string]any
	if err := json.Unmarshal(line, &entry); err != nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count++
	if len(s.bucket) < s.size {
		s.bucket = append(s.bucket, entry)
		return
	}
	// Reservoir replacement
	idx := s.rng.Intn(s.count)
	if idx < s.size {
		s.bucket[idx] = entry
	}
}

// Snapshot returns the current sample entries.
func (s *Sample) Snapshot() []map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]map[string]any, len(s.bucket))
	copy(out, s.bucket)
	return out
}

// Reset clears all sampled entries.
func (s *Sample) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bucket = s.bucket[:0]
	s.count = 0
}

// MarshalJSON serialises the current sample as a JSON array.
func (s *Sample) MarshalJSON() ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return json.Marshal(map[string]any{
		"sample_size": s.size,
		"count":       s.count,
		"entries":     s.bucket,
	})
}
