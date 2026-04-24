package conditional_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/conditional"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
		Raw:       "{}",
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	a := conditional.New(nil)
	e := makeEntry(map[string]any{"level": "error"})
	out := a.Apply(e)
	if out.Fields["level"] != "error" {
		t.Fatalf("expected 'error', got %v", out.Fields["level"])
	}
}

func TestApply_MatchingSetsTarget(t *testing.T) {
	rules := []conditional.Rule{
		{Field: "level", Match: "error", Target: "alert", Value: "true"},
	}
	a := conditional.New(rules)
	e := makeEntry(map[string]any{"level": "error"})
	out := a.Apply(e)
	if out.Fields["alert"] != "true" {
		t.Fatalf("expected alert=true, got %v", out.Fields["alert"])
	}
}

func TestApply_NonMatchingDoesNotSetTarget(t *testing.T) {
	rules := []conditional.Rule{
		{Field: "level", Match: "error", Target: "alert", Value: "true"},
	}
	a := conditional.New(rules)
	e := makeEntry(map[string]any{"level": "info"})
	out := a.Apply(e)
	if _, ok := out.Fields["alert"]; ok {
		t.Fatal("expected alert field to be absent")
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	rules := []conditional.Rule{
		{Field: "status", Match: "500", Target: "critical", Value: "yes"},
	}
	a := conditional.New(rules)
	e := makeEntry(map[string]any{"level": "error"})
	out := a.Apply(e)
	if _, ok := out.Fields["critical"]; ok {
		t.Fatal("expected critical field to be absent")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	rules := []conditional.Rule{
		{Field: "level", Match: "warn", Target: "flag", Value: "1"},
	}
	a := conditional.New(rules)
	origFields := map[string]any{"level": "warn"}
	e := makeEntry(origFields)
	_ = a.Apply(e)
	if _, ok := origFields["flag"]; ok {
		t.Fatal("original fields were mutated")
	}
}

func TestParseRules_ValidSpec(t *testing.T) {
	rules, err := conditional.ParseRules([]string{"level=error:alert=true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	r := rules[0]
	if r.Field != "level" || r.Match != "error" || r.Target != "alert" || r.Value != "true" {
		t.Fatalf("unexpected rule: %+v", r)
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := conditional.ParseRules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules != nil {
		t.Fatal("expected nil rules")
	}
}

func TestParseRules_MissingColon_ReturnsError(t *testing.T) {
	_, err := conditional.ParseRules([]string{"level=error"})
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseRules_EmptyField_ReturnsError(t *testing.T) {
	_, err := conditional.ParseRules([]string{"=error:alert=true"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}
