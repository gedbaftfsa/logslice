package aggregator

import (
	"encoding/json"
	"testing"
)

func makeComputeEntry(t *testing.T, fields map[string]interface{}) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return json.RawMessage(b)
}

func TestNewCompute_InvalidOp(t *testing.T) {
	_, err := NewCompute("x", "y", "mod", 2)
	if err == nil {
		t.Fatal("expected error for unsupported op")
	}
}

func TestNewCompute_DivByZero(t *testing.T) {
	_, err := NewCompute("x", "y", "div", 0)
	if err == nil {
		t.Fatal("expected error for division by zero")
	}
}

func TestCompute_Record_Add(t *testing.T) {
	c, _ := NewCompute("val", "val_plus_10", "add", 10)
	c.Record(makeComputeEntry(t, map[string]interface{}{"val": 5.0}))
	results := c.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	var m map[string]interface{}
	if err := json.Unmarshal(results[0], &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if m["val_plus_10"].(float64) != 15.0 {
		t.Errorf("expected 15, got %v", m["val_plus_10"])
	}
}

func TestCompute_Record_Mul(t *testing.T) {
	c, _ := NewCompute("n", "n2", "mul", 3)
	c.Record(makeComputeEntry(t, map[string]interface{}{"n": 4.0}))
	var m map[string]interface{}
	json.Unmarshal(c.Results()[0], &m)
	if m["n2"].(float64) != 12.0 {
		t.Errorf("expected 12, got %v", m["n2"])
	}
}

func TestCompute_Record_MissingField(t *testing.T) {
	c, _ := NewCompute("missing", "out", "add", 1)
	c.Record(makeComputeEntry(t, map[string]interface{}{"other": 1.0}))
	if len(c.Results()) != 0 {
		t.Error("expected no results for missing field")
	}
}

func TestCompute_Record_InvalidJSON(t *testing.T) {
	c, _ := NewCompute("x", "y", "add", 1)
	c.Record(json.RawMessage(`not-json`))
	if len(c.Results()) != 0 {
		t.Error("expected no results for invalid JSON")
	}
}

func TestCompute_Record_Abs(t *testing.T) {
	c, _ := NewCompute("v", "abs_v", "abs", 0)
	c.Record(makeComputeEntry(t, map[string]interface{}{"v": -7.5}))
	var m map[string]interface{}
	json.Unmarshal(c.Results()[0], &m)
	if m["abs_v"].(float64) != 7.5 {
		t.Errorf("expected 7.5, got %v", m["abs_v"])
	}
}

func TestCompute_Reset(t *testing.T) {
	c, _ := NewCompute("v", "out", "add", 1)
	c.Record(makeComputeEntry(t, map[string]interface{}{"v": 1.0}))
	c.Reset()
	if len(c.Results()) != 0 {
		t.Error("expected empty results after reset")
	}
}
