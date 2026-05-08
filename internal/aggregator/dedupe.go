package aggregator

import (
	"encoding/json"
	"sync"
	"time"
)

// Dedupe tracks recently seen log entries by a key field and suppresses duplicates
// within a configurable time window.
type Dedupe struct {
	mu      sync.Mutex
	field   string
	ttl     time.Duration
	seen    map[string]time.Time
	Dropped int
}

// NewDedupe creates a Dedupe that suppresses duplicate values of field within ttl.
func NewDedupe(field string, ttl time.Duration) *Dedupe {
	return &Dedupe{
		field: field,
		ttl:   ttl,
		seen:  make(map[string]time.Time),
	}
}

// IsDuplicate returns true if the entry's key field value was seen within the TTL.
// It also evicts expired entries and records the new value if not a duplicate.
func (d *Dedupe) IsDuplicate(entry []byte) (bool, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(entry, &m); err != nil {
		return false, err
	}

	val, ok := m[d.field]
	if !ok {
		return false, nil
	}

	key, ok := val.(string)
	if !ok {
		return false, nil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	d.evict(now)

	if last, exists := d.seen[key]; exists && now.Sub(last) < d.ttl {
		d.Dropped++
		return true, nil
	}

	d.seen[key] = now
	return false, nil
}

// evict removes entries older than the TTL. Must be called with mu held.
func (d *Dedupe) evict(now time.Time) {
	for k, t := range d.seen {
		if now.Sub(t) >= d.ttl {
			delete(d.seen, k)
		}
	}
}

// Reset clears all seen entries and resets the dropped counter.
func (d *Dedupe) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
	d.Dropped = 0
}
