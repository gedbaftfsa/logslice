package aggregator

import (
	"encoding/json"
	"time"
)

// Window holds log entries within a sliding time window.
type Window struct {
	duration time.Duration
	entries  []windowEntry
}

type windowEntry struct {
	ts  time.Time
	raw json.RawMessage
}

// NewWindow creates a Window that retains entries within the given duration.
func NewWindow(d time.Duration) *Window {
	return &Window{duration: d}
}

// Add inserts a raw JSON entry into the window, parsing the timestamp from
// the given field name. If the field is missing or unparseable, time.Now()
// is used as a fallback.
func (w *Window) Add(raw json.RawMessage, tsField string) {
	ts := extractTime(raw, tsField)
	w.entries = append(w.entries, windowEntry{ts: ts, raw: raw})
	w.evict()
}

// Entries returns all entries currently inside the window.
func (w *Window) Entries() []json.RawMessage {
	w.evict()
	out := make([]json.RawMessage, len(w.entries))
	for i, e := range w.entries {
		out[i] = e.raw
	}
	return out
}

// Len returns the number of entries currently in the window.
func (w *Window) Len() int {
	w.evict()
	return len(w.entries)
}

// evict removes entries older than the window duration relative to the newest entry.
func (w *Window) evict() {
	if len(w.entries) == 0 {
		return
	}
	newest := w.entries[len(w.entries)-1].ts
	cutoff := newest.Add(-w.duration)
	i := 0
	for i < len(w.entries) && w.entries[i].ts.Before(cutoff) {
		i++
	}
	w.entries = w.entries[i:]
}

// extractTime attempts to parse a time from the named field in a JSON object.
func extractTime(raw json.RawMessage, field string) time.Time {
	if field == "" {
		return time.Now()
	}
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		return time.Now()
	}
	v, ok := m[field]
	if !ok {
		return time.Now()
	}
	s, ok := v.(string)
	if !ok {
		return time.Now()
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now()
	}
	return t
}
