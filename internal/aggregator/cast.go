package aggregator

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Cast converts a field's value to a target type (string, int, float, bool).
type Cast struct {
	field  string
	target string
}

// NewCast creates a Cast aggregator that converts field to the target type.
// Supported targets: "string", "int", "float", "bool".
func NewCast(field, target string) (*Cast, error) {
	switch target {
	case "string", "int", "float", "bool":
		// valid
	default:
		return nil, fmt.Errorf("unsupported cast target %q: must be string, int, float, or bool", target)
	}
	return &Cast{field: field, target: target}, nil
}

// Record converts the field in the JSON entry and returns the modified entry.
// Returns an error if the entry is invalid JSON or the conversion fails.
func (c *Cast) Record(entry []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(entry, &m); err != nil {
		return nil, fmt.Errorf("cast: invalid JSON: %w", err)
	}

	raw, ok := m[c.field]
	if !ok {
		return entry, nil
	}

	converted, err := castValue(raw, c.target)
	if err != nil {
		return nil, fmt.Errorf("cast: field %q: %w", c.field, err)
	}

	m[c.field] = converted
	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("cast: marshal: %w", err)
	}
	return out, nil
}

func castValue(raw interface{}, target string) (interface{}, error) {
	s := fmt.Sprintf("%v", raw)
	switch target {
	case "string":
		return s, nil
	case "int":
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %q to int", s)
		}
		return int64(v), nil
	case "float":
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %q to float", s)
		}
		return v, nil
	case "bool":
		v, err := strconv.ParseBool(s)
		if err != nil {
			return nil, fmt.Errorf("cannot convert %q to bool", s)
		}
		return v, nil
	}
	return nil, fmt.Errorf("unknown target type %q", target)
}
