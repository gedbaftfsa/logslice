package aggregator

import (
	"encoding/json"
	"testing"
)

func makeTopKEntry(field, value string) []byte {
	b, _ := json.Marshal(map[string]string{field: value})
	return b
}

func TestNewTopK_InitialState(t *testing.T) {
	tk := NewTopK("level", 5)
	if got := tk.Top(); len(got) != 0 {
		t.Fatalf("expected empty top, got %d entries", len(got))
	}
}

func TestTopK_Record_StringField(t *testing.T) {
	tk := NewTopK("level", 3)
	for i := 0; i < 5; i++ {
		_ = tk.Record(makeTopKEntry("level", "error"))
	}
	for i := 0; i < 3; i++ {
		_ = tk.Record(makeTopKEntry("level", "warn"))
	}
	_ = tk.Record(makeTopKEntry("level", "info"))

	top := tk.Top()
	if len(top) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(top))
	}
	if top[0].Value != "error" || top[0].Count != 5 {
		t.Errorf("expected error/5, got %s/%d", top[0].Value, top[0].Count)
	}
	if top[1].Value != "warn" || top[1].Count != 3 {
		t.Errorf("expected warn/3, got %s/%d", top[1].Value, top[1].Count)
	}
}

func TestTopK_Record_MissingField(t *testing.T) {
	tk := NewTopK("level", 5)
	_ = tk.Record(makeTopKEntry("other", "value"))
	if len(tk.Top()) != 0 {
		t.Error("expected no entries for missing field")
	}
}

func TestTopK_Record_InvalidJSON(t *testing.T) {
	tk := NewTopK("level", 5)
	err := tk.Record([]byte("not json"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestTopK_Reset(t *testing.T) {
	tk := NewTopK("level", 5)
	_ = tk.Record(makeTopKEntry("level", "error"))
	tk.Reset()
	if len(tk.Top()) != 0 {
		t.Error("expected empty top after reset")
	}
}

func TestTopK_DefaultK(t *testing.T) {
	tk := NewTopK("level", 0)
	if tk.k != 10 {
		t.Errorf("expected default k=10, got %d", tk.k)
	}
}

func TestTopK_LimitsToK(t *testing.T) {
	tk := NewTopK("svc", 2)
	for _, v := range []string{"a", "b", "c", "d"} {
		_ = tk.Record(makeTopKEntry("svc", v))
	}
	if len(tk.Top()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(tk.Top()))
	}
}
