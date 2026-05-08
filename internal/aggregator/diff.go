package aggregator

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Diff tracks changes in a numeric field between consecutive log entries.
// It reports the delta (difference) between the current and previous value.
type Diff struct {
	mu       sync.Mutex
	field    string
	hasPrev  bool
	prevVal  float64
	deltas   []float64
	invalid  int
}

// NewDiff creates a new Diff tracker for the given numeric field.
func NewDiff(field string) *Diff {
	return &Diff{field: field}
}

// Record processes a JSON log entry and computes the delta from the previous value.
func (d *Diff) Record(entry []byte) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var m map[string]interface{}
	if err := json.Unmarshal(entry, &m); err != nil {
		d.invalid++
		return
	}

	raw, ok := m[d.field]
	if !ok {
		return
	}

	var val float64
	switch v := raw.(type) {
	case float64:
		val = v
	case int:
		val = float64(v)
	default:
		d.invalid++
		return
	}

	if d.hasPrev {
		d.deltas = append(d.deltas, val-d.prevVal)
	}
	d.prevVal = val
	d.hasPrev = true
}

// DiffSnapshot holds a point-in-time view of diff results.
type DiffSnapshot struct {
	Field   string    `json:"field"`
	Deltas  []float64 `json:"deltas"`
	Count   int       `json:"count"`
	Invalid int       `json:"invalid"`
}

// Snapshot returns the current diff state without resetting.
func (d *Diff) Snapshot() DiffSnapshot {
	d.mu.Lock()
	defer d.mu.Unlock()

	deltas := make([]float64, len(d.deltas))
	copy(deltas, d.deltas)
	return DiffSnapshot{
		Field:   d.field,
		Deltas:  deltas,
		Count:   len(deltas),
		Invalid: d.invalid,
	}
}

// Reset clears all recorded deltas and resets state.
func (d *Diff) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.deltas = nil
	d.hasPrev = false
	d.prevVal = 0
	d.invalid = 0
}

// MarshalJSON implements json.Marshaler.
func (d *Diff) MarshalJSON() ([]byte, error) {
	snap := d.Snapshot()
	return json.Marshal(snap)
}

// String returns a human-readable summary.
func (d *Diff) String() string {
	snap := d.Snapshot()
	return fmt.Sprintf("diff(%s): %d deltas, %d invalid", snap.Field, snap.Count, snap.Invalid)
}
