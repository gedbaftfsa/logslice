package aggregator

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

// Compute applies an arithmetic expression to a numeric field and emits
// the result as a new field on each log entry.
type Compute struct {
	srcField string
	dstField string
	op       string
	operand  float64
	results  []json.RawMessage
}

// NewCompute creates a Compute aggregator.
// op must be one of: add, sub, mul, div, abs, round.
func NewCompute(srcField, dstField, op string, operand float64) (*Compute, error) {
	switch op {
	case "add", "sub", "mul", "div", "abs", "round":
	default:
		return nil, fmt.Errorf("unsupported op %q: must be add, sub, mul, div, abs, or round", op)
	}
	if op == "div" && operand == 0 {
		return nil, fmt.Errorf("division by zero")
	}
	return &Compute{
		srcField: srcField,
		dstField: dstField,
		op:       op,
		operand:  operand,
	}, nil
}

// Record processes a single JSON log entry.
func (c *Compute) Record(entry json.RawMessage) {
	var m map[string]interface{}
	if err := json.Unmarshal(entry, &m); err != nil {
		return
	}
	raw, ok := m[c.srcField]
	if !ok {
		return
	}
	val, err := toFloat(raw)
	if err != nil {
		return
	}
	result := c.apply(val)
	m[c.dstField] = result
	out, err := json.Marshal(m)
	if err != nil {
		return
	}
	c.results = append(c.results, json.RawMessage(out))
}

func (c *Compute) apply(v float64) float64 {
	switch c.op {
	case "add":
		return v + c.operand
	case "sub":
		return v - c.operand
	case "mul":
		return v * c.operand
	case "div":
		return v / c.operand
	case "abs":
		return math.Abs(v)
	case "round":
		return math.Round(v)
	}
	return v
}

// Results returns all transformed entries.
func (c *Compute) Results() []json.RawMessage {
	return c.results
}

// Reset clears accumulated results.
func (c *Compute) Reset() {
	c.results = nil
}

func toFloat(v interface{}) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case string:
		return strconv.ParseFloat(n, 64)
	}
	return 0, fmt.Errorf("not a number")
}
