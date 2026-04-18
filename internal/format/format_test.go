package format_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/format"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Timestamp: time.Now(), Fields: fields, Raw: "{}"}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	f := format.New(nil)
	e := makeEntry(map[string]any{"latency": 1.23456})
	out := f.Apply(e)
	if out.Fields["latency"] != 1.23456 {
		t.Fatalf("expected unchanged value, got %v", out.Fields["latency"])
	}
}

func TestApply_FormatsFloatField(t *testing.T) {
	f := format.New([]format.Rule{{Field: "latency", Format: "%.2f"}})
	e := makeEntry(map[string]any{"latency": 3.14159})
	out := f.Apply(e)
	if out.Fields["latency"] != "3.14" {
		t.Fatalf("expected \"3.14\", got %v", out.Fields["latency"])
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	f := format.New([]format.Rule{{Field: "missing", Format: "%d"}})
	e := makeEntry(map[string]any{"other": 42})
	out := f.Apply(e)
	if _, ok := out.Fields["missing"]; ok {
		t.Fatal("expected missing field to remain absent")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	f := format.New([]format.Rule{{Field: "val", Format: "%05d"}})
	orig := map[string]any{"val": 7}
	e := makeEntry(orig)
	f.Apply(e)
	if orig["val"] != 7 {
		t.Fatal("input entry was mutated")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	f := format.New([]format.Rule{
		{Field: "a", Format: "%.1f"},
		{Field: "b", Format: "%05d"},
	})
	e := makeEntry(map[string]any{"a": 2.567, "b": 42})
	out := f.Apply(e)
	if out.Fields["a"] != "2.6" {
		t.Fatalf("expected \"2.6\", got %v", out.Fields["a"])
	}
	if out.Fields["b"] != "00042" {
		t.Fatalf("expected \"00042\", got %v", out.Fields["b"])
	}
}

func TestParseRules_ValidSpec(t *testing.T) {
	rules, err := format.ParseRules([]string{"latency=%.3f", "code=%d"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := format.ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil, nil; got %v, %v", rules, err)
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := format.ParseRules([]string{"badspec"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseRules_EmptyField_ReturnsError(t *testing.T) {
	_, err := format.ParseRules([]string{"=%.2f"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}
