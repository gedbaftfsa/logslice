package aggregator

import (
	"encoding/json"
	"testing"
)

func makeJoinEntry(t *testing.T, fields map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(fields)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestNewJoin_InitialState(t *testing.T) {
	j := NewJoin("id", "left", "right")
	if got := j.Flush(); len(got) != 0 {
		t.Fatalf("expected empty flush, got %d entries", len(got))
	}
}

func TestJoin_NoMatchWithoutBothSides(t *testing.T) {
	j := NewJoin("id", "left", "right")
	_ = j.Record("left", makeJoinEntry(t, map[string]interface{}{"id": "1", "val": "a"}))
	if got := j.Flush(); len(got) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(got))
	}
}

func TestJoin_MatchWhenBothSidesPresent(t *testing.T) {
	j := NewJoin("id", "left", "right")
	_ = j.Record("left", makeJoinEntry(t, map[string]interface{}{"id": "42", "name": "alice"}))
	_ = j.Record("right", makeJoinEntry(t, map[string]interface{}{"id": "42", "status": "ok"}))

	matches := j.Flush()
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0]["left.name"] != "alice" {
		t.Errorf("expected left.name=alice, got %v", matches[0]["left.name"])
	}
	if matches[0]["right.status"] != "ok" {
		t.Errorf("expected right.status=ok, got %v", matches[0]["right.status"])
	}
}

func TestJoin_RightBeforeLeft(t *testing.T) {
	j := NewJoin("id", "left", "right")
	_ = j.Record("right", makeJoinEntry(t, map[string]interface{}{"id": "7", "code": 200}))
	_ = j.Record("left", makeJoinEntry(t, map[string]interface{}{"id": "7", "path": "/api"}))

	matches := j.Flush()
	if len(matches) != 1 {
		t.Fatalf("expected 1 match, got %d", len(matches))
	}
	if matches[0]["left.path"] != "/api" {
		t.Errorf("expected left.path=/api, got %v", matches[0]["left.path"])
	}
}

func TestJoin_MissingKeySkipped(t *testing.T) {
	j := NewJoin("id", "left", "right")
	_ = j.Record("left", makeJoinEntry(t, map[string]interface{}{"other": "x"}))
	_ = j.Record("right", makeJoinEntry(t, map[string]interface{}{"id": "1", "v": "y"}))
	if got := j.Flush(); len(got) != 0 {
		t.Fatalf("expected 0 matches, got %d", len(got))
	}
}

func TestJoin_InvalidJSONReturnsError(t *testing.T) {
	j := NewJoin("id", "left", "right")
	if err := j.Record("left", []byte(`not-json`)); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestJoin_Reset(t *testing.T) {
	j := NewJoin("id", "left", "right")
	_ = j.Record("left", makeJoinEntry(t, map[string]interface{}{"id": "1", "x": 1}))
	_ = j.Record("right", makeJoinEntry(t, map[string]interface{}{"id": "1", "y": 2}))
	j.Reset()
	if got := j.Flush(); len(got) != 0 {
		t.Fatalf("expected 0 after reset, got %d", len(got))
	}
}
