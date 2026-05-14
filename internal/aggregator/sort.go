package aggregator

import (
	"encoding/json"
	"fmt"
	"sort"
)

// Sort buffers entries and emits them ordered by a numeric or string field.
type Sort struct {
	field string
	order string // "asc" or "desc"
	entries []map[string]interface{}
}

// NewSort creates a Sort aggregator. order must be "asc" or "desc".
func NewSort(field, order string) (*Sort, error) {
	if field == "" {
		return nil, fmt.Errorf("sort: field must not be empty")
	}
	if order != "asc" && order != "desc" {
		return nil, fmt.Errorf("sort: order must be 'asc' or 'desc', got %q", order)
	}
	return &Sort{field: field, order: order}, nil
}

// Record buffers a JSON entry for later sorting.
func (s *Sort) Record(line []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return fmt.Errorf("sort: invalid JSON: %w", err)
	}
	s.entries = append(s.entries, m)
	return nil
}

// Flush sorts the buffered entries and returns them as JSON lines.
func (s *Sort) Flush() ([][]byte, error) {
	sorted := make([]map[string]interface{}, len(s.entries))
	copy(sorted, s.entries)

	sort.SliceStable(sorted, func(i, j int) bool {
		vi := sorted[i][s.field]
		vj := sorted[j][s.field]
		less := compareValues(vi, vj)
		if s.order == "desc" {
			return !less
		}
		return less
	})

	out := make([][]byte, 0, len(sorted))
	for _, m := range sorted {
		b, err := json.Marshal(m)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

// Reset clears buffered entries.
func (s *Sort) Reset() {
	s.entries = nil
}

// compareValues returns true if a < b for numeric or string types.
func compareValues(a, b interface{}) bool {
	switch av := a.(type) {
	case float64:
		if bv, ok := b.(float64); ok {
			return av < bv
		}
	case string:
		if bv, ok := b.(string); ok {
			return av < bv
		}
	}
	return false
}
