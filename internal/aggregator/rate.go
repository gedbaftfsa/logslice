package aggregator

import (
	"encoding/json"
	"sync"
	"time"
)

// Rate tracks the number of log entries per second over a sliding window.
type Rate struct {
	mu       sync.Mutex
	window   time.Duration
	timestamps []time.Time
}

// NewRate creates a Rate tracker with the given sliding window duration.
func NewRate(window time.Duration) *Rate {
	return &Rate{
		window: window,
	}
}

// Record registers a new event at the current time.
func (r *Rate) Record(line []byte) {
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	r.timestamps = append(r.timestamps, now)
	r.evict(now)
}

// evict removes timestamps outside the current window. Must be called with lock held.
func (r *Rate) evict(now time.Time) {
	cutoff := now.Add(-r.window)
	i := 0
	for i < len(r.timestamps) && r.timestamps[i].Before(cutoff) {
		i++
	}
	r.timestamps = r.timestamps[i:]
}

// PerSecond returns the average events per second within the window.
func (r *Rate) PerSecond() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	r.evict(now)
	if len(r.timestamps) == 0 {
		return 0
	}
	secs := r.window.Seconds()
	if secs <= 0 {
		return 0
	}
	return float64(len(r.timestamps)) / secs
}

// Count returns the number of events currently within the window.
func (r *Rate) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.evict(time.Now())
	return len(r.timestamps)
}

// Reset clears all recorded timestamps.
func (r *Rate) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.timestamps = nil
}
