package aggregator

import (
	"encoding/json"
	"sort"
	"testing"
)

func makeUniqueEntry(t *testing.T, fields map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewUnique_InitialState(t *testing.T) {
	u := NewUnique("level")
	if u.Count() != 0 {
		t.Fatalf("expected 0, got %d", u.Count())
	}
	if len(u.Values()) != 0 {
		t.Fatal("expected empty values")
	}
}

func TestUnique_Record_StringField(t *testing.T) {
	u := NewUnique("level")
	u.Record(makeUniqueEntry(t, map[string]interface{}{"level": "info"}))
	u.Record(makeUniqueEntry(t, map[string]interface{}{"level": "warn"}))
	u.Record(makeUniqueEntry(t, map[string]interface{}{"level": "info"}))
	if u.Count() != 2 {
		t.Fatalf("expected 2 unique values, got %d", u.Count())
	}
}

func TestUnique_Record_MissingField(t *testing.T) {
	u := NewUnique("level")
	u.Record(makeUniqueEntry(t, map[string]interface{}{"msg": "hello"}))
	if u.Count() != 0 {
		t.Fatalf("expected 0, got %d", u.Count())
	}
}

func TestUnique_Record_InvalidJSON(t *testing.T) {
	u := NewUnique("level")
	u.Record([]byte("not-json"))
	if u.Count() != 0 {
		t.Fatalf("expected 0, got %d", u.Count())
	}
}

func TestUnique_Values(t *testing.T) {
	u := NewUnique("env")
	u.Record(makeUniqueEntry(t, map[string]interface{}{"env": "prod"}))
	u.Record(makeUniqueEntry(t, map[string]interface{}{"env": "staging"}))
	u.Record(makeUniqueEntry(t, map[string]interface{}{"env": "prod"}))
	vals := u.Values()
	sort.Strings(vals)
	if len(vals) != 2 || vals[0] != "prod" || vals[1] != "staging" {
		t.Fatalf("unexpected values: %v", vals)
	}
}

func TestUnique_Reset(t *testing.T) {
	u := NewUnique("level")
	u.Record(makeUniqueEntry(t, map[string]interface{}{"level": "info"}))
	u.Reset()
	if u.Count() != 0 {
		t.Fatalf("expected 0 after reset, got %d", u.Count())
	}
}

func TestUnique_MarshalJSON(t *testing.T) {
	u := NewUnique("level")
	u.Record(makeUniqueEntry(t, map[string]interface{}{"level": "info"}))
	b, err := u.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out["field"] != "level" {
		t.Errorf("expected field=level, got %v", out["field"])
	}
	if out["count"].(float64) != 1 {
		t.Errorf("expected count=1, got %v", out["count"])
	}
}
