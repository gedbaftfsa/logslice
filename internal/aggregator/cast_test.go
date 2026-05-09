package aggregator

import (
	"encoding/json"
	"testing"
)

func makeCastEntry(t *testing.T, m map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewCast_InvalidTarget(t *testing.T) {
	_, err := NewCast("field", "bytes")
	if err == nil {
		t.Fatal("expected error for unsupported target")
	}
}

func TestCast_ToInt(t *testing.T) {
	c, _ := NewCast("count", "int")
	entry := makeCastEntry(t, map[string]interface{}{"count": "42"})
	out, err := c.Record(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	json.Unmarshal(out, &m)
	if v, ok := m["count"].(float64); !ok || int(v) != 42 {
		t.Errorf("expected count=42, got %v", m["count"])
	}
}

func TestCast_ToFloat(t *testing.T) {
	c, _ := NewCast("value", "float")
	entry := makeCastEntry(t, map[string]interface{}{"value": "3.14"})
	out, err := c.Record(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	json.Unmarshal(out, &m)
	if v, ok := m["value"].(float64); !ok || v < 3.13 || v > 3.15 {
		t.Errorf("expected value~3.14, got %v", m["value"])
	}
}

func TestCast_ToBool(t *testing.T) {
	c, _ := NewCast("active", "bool")
	entry := makeCastEntry(t, map[string]interface{}{"active": "true"})
	out, err := c.Record(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	json.Unmarshal(out, &m)
	if v, ok := m["active"].(bool); !ok || !v {
		t.Errorf("expected active=true, got %v", m["active"])
	}
}

func TestCast_ToString(t *testing.T) {
	c, _ := NewCast("code", "string")
	entry := makeCastEntry(t, map[string]interface{}{"code": 404})
	out, err := c.Record(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	json.Unmarshal(out, &m)
	if v, ok := m["code"].(string); !ok || v != "404" {
		t.Errorf("expected code=\"404\", got %v", m["code"])
	}
}

func TestCast_MissingField(t *testing.T) {
	c, _ := NewCast("missing", "int")
	entry := makeCastEntry(t, map[string]interface{}{"other": "val"})
	out, err := c.Record(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) == "" {
		t.Error("expected non-empty output")
	}
}

func TestCast_InvalidJSON(t *testing.T) {
	c, _ := NewCast("x", "int")
	_, err := c.Record([]byte("not-json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
