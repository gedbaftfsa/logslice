package aggregator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Extract pulls a value from a source field using a regex capture group
// and writes it to a new destination field in each log entry.
type Extract struct {
	srcField  string
	dstField  string
	pattern   *regexp.Regexp
	entries   []map[string]any
}

// NewExtract creates an Extract aggregator.
// srcField is the field to read from, dstField is the field to write to,
// and pattern must contain exactly one capture group.
func NewExtract(srcField, dstField, pattern string) (*Extract, error) {
	if srcField == "" || dstField == "" {
		return nil, fmt.Errorf("extract: srcField and dstField must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("extract: invalid pattern: %w", err)
	}
	if re.NumSubexp() < 1 {
		return nil, fmt.Errorf("extract: pattern must contain at least one capture group")
	}
	return &Extract{
		srcField: srcField,
		dstField: dstField,
		pattern:  re,
	}, nil
}

// Record processes a single JSON log line.
func (e *Extract) Record(line string) error {
	var m map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &m); err != nil {
		return nil // skip invalid JSON
	}

	if val, ok := m[e.srcField]; ok {
		str := fmt.Sprintf("%v", val)
		matches := e.pattern.FindStringSubmatch(str)
		if len(matches) >= 2 {
			m[e.dstField] = matches[1]
		}
	}

	e.entries = append(e.entries, m)
	return nil
}

// Snapshot returns all processed entries as JSON lines.
func (e *Extract) Snapshot() []string {
	out := make([]string, 0, len(e.entries))
	for _, m := range e.entries {
		b, err := json.Marshal(m)
		if err != nil {
			continue
		}
		out = append(out, string(b))
	}
	return out
}

// Reset clears all recorded entries.
func (e *Extract) Reset() {
	e.entries = nil
}
