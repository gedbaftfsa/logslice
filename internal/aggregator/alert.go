package aggregator

import (
	"encoding/json"
	"fmt"
	"sync"
)

// AlertCondition defines when an alert should fire.
type AlertCondition struct {
	Field    string
	Operator string // "gt", "lt", "eq"
	Threshold float64
}

// Alert tracks how many log entries match a threshold condition on a numeric field.
type Alert struct {
	mu        sync.Mutex
	condition AlertCondition
	triggered int
	total     int
}

// NewAlert creates an Alert for the given condition.
func NewAlert(cond AlertCondition) *Alert {
	return &Alert{condition: cond}
}

// Record evaluates a JSON log entry against the alert condition.
func (a *Alert) Record(entry []byte) {
	a.mu.Lock()
	defer a.mu.Unlock()

	var m map[string]interface{}
	if err := json.Unmarshal(entry, &m); err != nil {
		return
	}
	a.total++

	raw, ok := m[a.condition.Field]
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
		return
	}

	switch a.condition.Operator {
	case "gt":
		if val > a.condition.Threshold {
			a.triggered++
		}
	case "lt":
		if val < a.condition.Threshold {
			a.triggered++
		}
	case "eq":
		if val == a.condition.Threshold {
			a.triggered++
		}
	}
}

// Snapshot returns a JSON summary of alert state.
func (a *Alert) Snapshot() []byte {
	a.mu.Lock()
	defer a.mu.Unlock()

	out := map[string]interface{}{
		"field":     a.condition.Field,
		"operator":  a.condition.Operator,
		"threshold": a.condition.Threshold,
		"triggered": a.triggered,
		"total":     a.total,
		"firing":    a.triggered > 0,
	}
	b, _ := json.Marshal(out)
	return b
}

// Reset clears the alert counters.
func (a *Alert) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.triggered = 0
	a.total = 0
}

// String returns a human-readable summary.
func (a *Alert) String() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return fmt.Sprintf("alert[%s %s %.2f]: %d/%d triggered",
		a.condition.Field, a.condition.Operator, a.condition.Threshold,
		a.triggered, a.total)
}
