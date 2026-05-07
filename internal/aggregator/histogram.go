package aggregator

import (
	"encoding/json"
	"math"
	"sort"
	"sync"
)

// Histogram tracks the distribution of a numeric field across log entries.
type Histogram struct {
	mu     sync.Mutex
	field  string
	values []float64
}

// NewHistogram creates a Histogram that tracks the given numeric field.
func NewHistogram(field string) *Histogram {
	return &Histogram{field: field}
}

// Record ingests a JSON log entry and records the numeric field value if present.
func (h *Histogram) Record(line []byte) {
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return
	}
	v, ok := entry[h.field]
	if !ok {
		return
	}
	var f float64
	switch val := v.(type) {
	case float64:
		f = val
	default:
		return
	}
	h.mu.Lock()
	h.values = append(h.values, f)
	h.mu.Unlock()
}

// Reset clears all recorded values.
func (h *Histogram) Reset() {
	h.mu.Lock()
	h.values = h.values[:0]
	h.mu.Unlock()
}

// HistogramSnapshot holds computed distribution statistics.
type HistogramSnapshot struct {
	Field  string             `json:"field"`
	Count  int                `json:"count"`
	Min    float64            `json:"min"`
	Max    float64            `json:"max"`
	Mean   float64            `json:"mean"`
	P50    float64            `json:"p50"`
	P90    float64            `json:"p90"`
	P99    float64            `json:"p99"`
}

// Snapshot returns a point-in-time view of the histogram statistics.
func (h *Histogram) Snapshot() HistogramSnapshot {
	h.mu.Lock()
	copied := make([]float64, len(h.values))
	copy(copied, h.values)
	h.mu.Unlock()

	snap := HistogramSnapshot{Field: h.field, Count: len(copied)}
	if len(copied) == 0 {
		return snap
	}
	sort.Float64s(copied)
	snap.Min = copied[0]
	snap.Max = copied[len(copied)-1]
	var sum float64
	for _, v := range copied {
		sum += v
	}
	snap.Mean = sum / float64(len(copied))
	snap.P50 = percentile(copied, 50)
	snap.P90 = percentile(copied, 90)
	snap.P99 = percentile(copied, 99)
	return snap
}

func percentile(sorted []float64, p float64) float64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(math.Ceil(p/100.0*float64(len(sorted)))) - 1
	if idx < 0 {
		idx = 0
	}
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
