package aggregator

import (
	"encoding/json"
	"testing"
)

func makeSortEntry(t *testing.T, m map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestNewSort_InvalidArgs(t *testing.T) {
	if _, err := NewSort("", "asc"); err == nil {
		t.Fatal("expected error for empty field")
	}
	if _, err := NewSort("level", "random"); err == nil {
		t.Fatal("expected error for invalid order")
	}
}

func TestSort_NumericAsc(t *testing.T) {
	s, _ := NewSort("code", "asc")
	for _, v := range []float64{3, 1, 2} {
		s.Record(makeSortEntry(t, map[string]interface{}{"code": v}))
	}
	out, err := s.Flush()
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	expected := []float64{1, 2, 3}
	for i, b := range out {
		var m map[string]interface{}
		json.Unmarshal(b, &m)
		if m["code"].(float64) != expected[i] {
			t.Errorf("pos %d: got %v, want %v", i, m["code"], expected[i])
		}
	}
}

func TestSort_NumericDesc(t *testing.T) {
	s, _ := NewSort("code", "desc")
	for _, v := range []float64{1, 3, 2} {
		s.Record(makeSortEntry(t, map[string]interface{}{"code": v}))
	}
	out, _ := s.Flush()
	expected := []float64{3, 2, 1}
	for i, b := range out {
		var m map[string]interface{}
		json.Unmarshal(b, &m)
		if m["code"].(float64) != expected[i] {
			t.Errorf("pos %d: got %v, want %v", i, m["code"], expected[i])
		}
	}
}

func TestSort_StringAsc(t *testing.T) {
	s, _ := NewSort("name", "asc")
	for _, v := range []string{"charlie", "alice", "bob"} {
		s.Record(makeSortEntry(t, map[string]interface{}{"name": v}))
	}
	out, _ := s.Flush()
	expected := []string{"alice", "bob", "charlie"}
	for i, b := range out {
		var m map[string]interface{}
		json.Unmarshal(b, &m)
		if m["name"].(string) != expected[i] {
			t.Errorf("pos %d: got %v, want %v", i, m["name"], expected[i])
		}
	}
}

func TestSort_InvalidJSON(t *testing.T) {
	s, _ := NewSort("field", "asc")
	if err := s.Record([]byte(`not-json`)); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestSort_Reset(t *testing.T) {
	s, _ := NewSort("x", "asc")
	s.Record(makeSortEntry(t, map[string]interface{}{"x": 1.0}))
	s.Reset()
	out, _ := s.Flush()
	if len(out) != 0 {
		t.Fatalf("expected 0 after reset, got %d", len(out))
	}
}
