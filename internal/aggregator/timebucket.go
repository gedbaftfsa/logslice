package aggregator

import (
	"encoding/json"
	"sync"
	"time"
)

// TimeBucket aggregates counts of log entries bucketed by a fixed time interval.
type TimeBucket struct {
	mu       sync.Mutex
	field    string
	interval time.Duration
	buckets  map[int64]int64
}

// NewTimeBucket creates a TimeBucket that groups entries by the given time field
// and buckets them into intervals of the given duration.
func NewTimeBucket(field string, interval time.Duration) *TimeBucket {
	return &TimeBucket{
		field:    field,
		interval: interval,
		buckets:  make(map[int64]int64),
	}
}

// Record parses a JSON log entry and increments the bucket for the entry's timestamp.
func (tb *TimeBucket) Record(line []byte) {
	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return
	}
	val, ok := obj[tb.field]
	if !ok {
		return
	}
	var t time.Time
	switch v := val.(type) {
	case string:
		parsed, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return
		}
		t = parsed
	case float64:
		t = time.Unix(int64(v), 0)
	default:
		return
	}
	key := t.Truncate(tb.interval).Unix()
	tb.mu.Lock()
	tb.buckets[key]++
	tb.mu.Unlock()
}

// Snapshot returns a copy of the current bucket counts keyed by bucket start time.
func (tb *TimeBucket) Snapshot() map[time.Time]int64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	out := make(map[time.Time]int64, len(tb.buckets))
	for k, v := range tb.buckets {
		out[time.Unix(k, 0).UTC()] = v
	}
	return out
}

// Reset clears all bucket data.
func (tb *TimeBucket) Reset() {
	tb.mu.Lock()
	tb.buckets = make(map[int64]int64)
	tb.mu.Unlock()
}

// MarshalJSON serialises the snapshot as a JSON object with RFC3339 keys.
func (tb *TimeBucket) MarshalJSON() ([]byte, error) {
	snap := tb.Snapshot()
	out := make(map[string]int64, len(snap))
	for t, c := range snap {
		out[t.Format(time.RFC3339)] = c
	}
	return json.Marshal(out)
}
