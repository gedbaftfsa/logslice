package aggregator

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// Flatten reads JSON log entries and expands nested objects into dot-notation keys.
// For example, {"a":{"b":1}} becomes {"a.b":1}.
type Flatten struct {
	mu      sync.Mutex
	prefix  string
	results []map[string]any
}

// NewFlatten creates a Flatten transformer. prefix is prepended to all
// top-level keys (use "" for none).
func NewFlatten(prefix string) *Flatten {
	return &Flatten{prefix: prefix}
}

// Record ingests a raw JSON log line and stores the flattened result.
func (f *Flatten) Record(line []byte) error {
	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		return fmt.Errorf("flatten: invalid JSON: %w", err)
	}

	flat := make(map[string]any)
	flattenMap(f.prefix, obj, flat)

	f.mu.Lock()
	f.results = append(f.results, flat)
	f.mu.Unlock()
	return nil
}

// Snapshot returns all flattened entries recorded so far.
func (f *Flatten) Snapshot() []map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]map[string]any, len(f.results))
	copy(out, f.results)
	return out
}

// Reset clears all stored entries.
func (f *Flatten) Reset() {
	f.mu.Lock()
	f.results = f.results[:0]
	f.mu.Unlock()
}

// flattenMap recursively walks obj, building dot-notation keys into dst.
func flattenMap(prefix string, obj map[string]any, dst map[string]any) {
	for k, v := range obj {
		key := k
		if prefix != "" {
			key = strings.Join([]string{prefix, k}, ".")
		}
		switch child := v.(type) {
		case map[string]any:
			flattenMap(key, child, dst)
		default:
			dst[key] = v
		}
	}
}
