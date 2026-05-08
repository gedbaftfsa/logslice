package aggregator

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Join correlates log entries from two named streams by a shared key field.
// When both sides have seen a value for the key, it emits a merged record.
type Join struct {
	mu      sync.Mutex
	keyField string
	leftTag  string
	rightTag string
	left    map[string]map[string]interface{}
	right   map[string]map[string]interface{}
	matched []map[string]interface{}
}

// NewJoin creates a Join that merges records sharing keyField across leftTag and rightTag streams.
func NewJoin(keyField, leftTag, rightTag string) *Join {
	return &Join{
		keyField: keyField,
		leftTag:  leftTag,
		rightTag: rightTag,
		left:    make(map[string]map[string]interface{}),
		right:   make(map[string]map[string]interface{}),
	}
}

// Record ingests a raw JSON line tagged with a source label.
func (j *Join) Record(tag string, line []byte) error {
	var entry map[string]interface{}
	if err := json.Unmarshal(line, &entry); err != nil {
		return fmt.Errorf("join: invalid JSON: %w", err)
	}
	keyVal, ok := entry[j.keyField]
	if !ok {
		return nil
	}
	key := fmt.Sprintf("%v", keyVal)

	j.mu.Lock()
	defer j.mu.Unlock()

	switch tag {
	case j.leftTag:
		j.left[key] = entry
		if right, found := j.right[key]; found {
			j.matched = append(j.matched, j.merge(entry, right))
		}
	case j.rightTag:
		j.right[key] = entry
		if left, found := j.left[key]; found {
			j.matched = append(j.matched, j.merge(left, entry))
		}
	}
	return nil
}

func (j *Join) merge(left, right map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(left)+len(right))
	for k, v := range left {
		out[j.leftTag+"."+k] = v
	}
	for k, v := range right {
		out[j.rightTag+"."+k] = v
	}
	return out
}

// Flush returns all matched (joined) records and resets the match buffer.
func (j *Join) Flush() []map[string]interface{} {
	j.mu.Lock()
	defer j.mu.Unlock()
	out := j.matched
	j.matched = nil
	return out
}

// Reset clears all buffered state.
func (j *Join) Reset() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.left = make(map[string]map[string]interface{})
	j.right = make(map[string]map[string]interface{})
	j.matched = nil
}
