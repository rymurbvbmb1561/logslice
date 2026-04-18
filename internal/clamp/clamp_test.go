package clamp_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/clamp"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Timestamp: time.Time{}, Fields: fields, Raw: ""}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	c := clamp.New(nil)
	e := makeEntry(map[string]any{"val": 5.0})
	out := c.Apply(e)
	if out.Fields["val"] != 5.0 {
		t.Fatalf("expected 5.0, got %v", out.Fields["val"])
	}
}

func TestApply_ClampsAboveMax(t *testing.T) {
	c := clamp.New([]clamp.Rule{{Field: "score", Min: 0, Max: 100}})
	out := c.Apply(makeEntry(map[string]any{"score": 150.0}))
	if out.Fields["score"] != 100.0 {
		t.Fatalf("expected 100, got %v", out.Fields["score"])
	}
}

func TestApply_ClampsBelowMin(t *testing.T) {
	c := clamp.New([]clamp.Rule{{Field: "score", Min: 0, Max: 100}})
	out := c.Apply(makeEntry(map[string]any{"score": -10.0}))
	if out.Fields["score"] != 0.0 {
		t.Fatalf("expected 0, got %v", out.Fields["score"])
	}
}

func TestApply_WithinRange_Unchanged(t *testing.T) {
	c := clamp.New([]clamp.Rule{{Field: "score", Min: 0, Max: 100}})
	out := c.Apply(makeEntry(map[string]any{"score": 42.0}))
	if out.Fields["score"] != 42.0 {
		t.Fatalf("expected 42, got %v", out.Fields["score"])
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	c := clamp.New([]clamp.Rule{{Field: "missing", Min: 0, Max: 10}})
	out := c.Apply(makeEntry(map[string]any{"other": 5.0}))
	if _, ok := out.Fields["missing"]; ok {
		t.Fatal("expected missing field to remain absent")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	c := clamp.New([]clamp.Rule{{Field: "x", Min: 0, Max: 1}})
	e := makeEntry(map[string]any{"x": 99.0})
	c.Apply(e)
	if e.Fields["x"] != 99.0 {
		t.Fatal("input was mutated")
	}
}

func TestParseRules_ValidSpec(t *testing.T) {
	rules, err := clamp.ParseRules([]string{"temp=0:100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Field != "temp" || rules[0].Min != 0 || rules[0].Max != 100 {
		t.Fatalf("unexpected rules: %+v", rules)
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := clamp.ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil, got %v %v", rules, err)
	}
}

func TestParseRules_MinGreaterThanMax_ReturnsError(t *testing.T) {
	_, err := clamp.ParseRules([]string{"x=10:1"})
	if err == nil {
		t.Fatal("expected error for min > max")
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := clamp.ParseRules([]string{"badspec"})
	if err == nil {
		t.Fatal("expected error")
	}
}
