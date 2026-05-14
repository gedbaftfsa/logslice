package aggregator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Mask partially obscures the value of specified fields using a regex pattern,
// replacing matched groups with a mask string (default "***").
type Mask struct {
	field   string
	re      *regexp.Regexp
	maskStr string
	buf     []map[string]any
}

// NewMask creates a Mask aggregator. pattern is a Go regex; any capturing group
// in the match is replaced by maskStr. If maskStr is empty, "***" is used.
func NewMask(field, pattern, maskStr string) (*Mask, error) {
	if field == "" {
		return nil, fmt.Errorf("mask: field must not be empty")
	}
	if pattern == "" {
		return nil, fmt.Errorf("mask: pattern must not be empty")
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("mask: invalid pattern: %w", err)
	}
	if maskStr == "" {
		maskStr = "***"
	}
	return &Mask{field: field, re: re, maskStr: maskStr}, nil
}

// Record processes a single JSON log entry, masking the target field if present.
func (m *Mask) Record(entry []byte) error {
	var row map[string]any
	if err := json.Unmarshal(entry, &row); err != nil {
		return nil // skip invalid JSON
	}
	val, ok := row[m.field]
	if !ok {
		m.buf = append(m.buf, row)
		return nil
	}
	str, ok := val.(string)
	if !ok {
		m.buf = append(m.buf, row)
		return nil
	}
	masked := m.re.ReplaceAllStringFunc(str, func(match string) string {
		// Replace each captured sub-group; if no groups, replace whole match.
		result := match
		for i, sub := range m.re.FindStringSubmatch(match) {
			if i == 0 {
				continue
			}
			result = strings.ReplaceAll(result, sub, m.maskStr)
		}
		if result == match {
			// no capturing groups — replace whole match
			return m.maskStr
		}
		return result
	})
	row[m.field] = masked
	m.buf = append(m.buf, row)
	return nil
}

// Snapshot returns all processed entries as JSON lines.
func (m *Mask) Snapshot() [][]byte {
	out := make([][]byte, 0, len(m.buf))
	for _, row := range m.buf {
		b, err := json.Marshal(row)
		if err != nil {
			continue
		}
		out = append(out, b)
	}
	return out
}

// Reset clears buffered entries.
func (m *Mask) Reset() {
	m.buf = nil
}
