package filter

import (
	"testing"
)

func TestParse_Empty(t *testing.T) {
	f, err := Parse("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Exprs) != 0 {
		t.Errorf("expected 0 exprs, got %d", len(f.Exprs))
	}
}

func TestParse_Valid(t *testing.T) {
	f, err := Parse("level=error status>=400")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Exprs) != 2 {
		t.Fatalf("expected 2 exprs, got %d", len(f.Exprs))
	}
	if f.Exprs[0].Field != "level" || f.Exprs[0].Operator != "=" || f.Exprs[0].Value != "error" {
		t.Errorf("unexpected first expr: %+v", f.Exprs[0])
	}
	if f.Exprs[1].Field != "status" || f.Exprs[1].Operator != ">=" || f.Exprs[1].Value != "400" {
		t.Errorf("unexpected second expr: %+v", f.Exprs[1])
	}
}

func TestParse_Invalid(t *testing.T) {
	_, err := Parse("badexpr")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestMatch_EmptyFilter(t *testing.T) {
	f := &Filter{}
	if !f.Match([]byte(`{"level":"info"}`)) {
		t.Error("empty filter should match everything")
	}
}

func TestMatch_EqualString(t *testing.T) {
	f, _ := Parse("level=error")
	if !f.Match([]byte(`{"level":"error","msg":"oops"}`)) {
		t.Error("should match level=error")
	}
	if f.Match([]byte(`{"level":"info","msg":"ok"}`)) {
		t.Error("should not match level=info")
	}
}

func TestMatch_NumericGTE(t *testing.T) {
	f, _ := Parse("status>=400")
	if !f.Match([]byte(`{"status":500}`)) {
		t.Error("500 >= 400 should match")
	}
	if !f.Match([]byte(`{"status":400}`)) {
		t.Error("400 >= 400 should match")
	}
	if f.Match([]byte(`{"status":200}`)) {
		t.Error("200 >= 400 should not match")
	}
}

func TestMatch_NotEqual(t *testing.T) {
	f, _ := Parse("level!=debug")
	if !f.Match([]byte(`{"level":"info"}`)) {
		t.Error("info != debug should match")
	}
	if f.Match([]byte(`{"level":"debug"}`)) {
		t.Error("debug != debug should not match")
	}
}

func TestMatch_MissingField(t *testing.T) {
	f, _ := Parse("level=error")
	if f.Match([]byte(`{"msg":"hello"}`)) {
		t.Error("missing field should not match")
	}
}

func TestMatch_InvalidJSON(t *testing.T) {
	f, _ := Parse("level=error")
	if f.Match([]byte(`not json`)) {
		t.Error("invalid JSON should not match")
	}
}
