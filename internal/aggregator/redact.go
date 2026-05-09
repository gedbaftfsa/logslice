package aggregator

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Redact replaces the value of specified fields with a placeholder string.
type Redact struct {
	fields      []string
	placeholder string
	count       int
}

// NewRedact creates a Redact processor that masks the given fields.
// placeholder defaults to "***" if empty.
func NewRedact(fields []string, placeholder string) (*Redact, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("redact: at least one field is required")
	}
	if placeholder == "" {
		placeholder = "***"
	}
	return &Redact{fields: fields, placeholder: placeholder}, nil
}

// Record processes a JSON log line and redacts the configured fields.
func (r *Redact) Record(line []byte) ([]byte, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return line, nil
	}

	modified := false
	for _, f := range r.fields {
		if _, ok := obj[f]; ok {
			obj[f] = r.placeholder
			modified = true
		}
	}

	if !modified {
		return line, nil
	}

	r.count++
	out, err := json.Marshal(obj)
	if err != nil {
		return line, nil
	}
	return out, nil
}

// Snapshot returns a JSON summary of how many lines were redacted.
func (r *Redact) Snapshot() []byte {
	fields := strings.Join(r.fields, ",")
	b, _ := json.Marshal(map[string]interface{}{
		"redacted_fields": fields,
		"lines_redacted":  r.count,
		"placeholder":     r.placeholder,
	})
	return b
}

// Reset clears the redaction counter.
func (r *Redact) Reset() {
	r.count = 0
}
