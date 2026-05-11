package aggregator

import (
	"encoding/json"
	"fmt"
)

// Coalesce picks the first non-empty value from a list of fields and writes
// it to a destination field. Useful for normalising logs that use different
// field names for the same concept across sources.
type Coalesce struct {
	fields []string
	dest   string
	results []map[string]any
}

// NewCoalesce creates a Coalesce processor.
// fields is the ordered list of source fields to try; dest is the output field.
func NewCoalesce(fields []string, dest string) (*Coalesce, error) {
	if len(fields) < 2 {
		return nil, fmt.Errorf("coalesce requires at least 2 source fields")
	}
	if dest == "" {
		return nil, fmt.Errorf("coalesce requires a non-empty destination field")
	}
	return &Coalesce{fields: fields, dest: dest}, nil
}

// Record processes a single JSON log line.
func (c *Coalesce) Record(line []byte) error {
	var entry map[string]any
	if err := json.Unmarshal(line, &entry); err != nil {
		return err
	}

	for _, f := range c.fields {
		val, ok := entry[f]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok || s == "" {
			continue
		}
		entry[c.dest] = s
		c.results = append(c.results, entry)
		return nil
	}

	// No non-empty value found; still pass the entry through unchanged.
	c.results = append(c.results, entry)
	return nil
}

// Snapshot returns all processed entries as JSON lines.
func (c *Coalesce) Snapshot() [][]byte {
	out := make([][]byte, 0, len(c.results))
	for _, e := range c.results {
		b, err := json.Marshal(e)
		if err != nil {
			continue
		}
		out = append(out, b)
	}
	return out
}

// Reset clears accumulated results.
func (c *Coalesce) Reset() {
	c.results = nil
}
