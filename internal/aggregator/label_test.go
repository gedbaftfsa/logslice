package aggregator_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/aggregator"
)

func TestInjectSourceField_AddsKey(t *testing.T) {
	fields := map[string]interface{}{"msg": "hello"}
	out := aggregator.InjectSourceField(fields, "myapp")
	if out["_source"] != "myapp" {
		t.Errorf("expected _source=myapp, got %v", out["_source"])
	}
	if out["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", out["msg"])
	}
}

func TestInjectSourceField_DoesNotMutateOriginal(t *testing.T) {
	fields := map[string]interface{}{"k": "v"}
	aggregator.InjectSourceField(fields, "x")
	if _, ok := fields["_source"]; ok {
		t.Error("original map should not be mutated")
	}
}

func TestFormatLabel(t *testing.T) {
	got := aggregator.FormatLabel("svc")
	if got != "[svc]" {
		t.Errorf("expected [svc], got %q", got)
	}
}
