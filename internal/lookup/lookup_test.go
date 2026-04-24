package lookup_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/lookup"
	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
		Raw:       "{}",
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	a := lookup.New(nil)
	e := makeEntry(map[string]any{"level": "info"})
	out := a.Apply(e)
	if out.Fields["level"] != "info" {
		t.Fatalf("expected 'info', got %v", out.Fields["level"])
	}
}

func TestApply_LooksUpAndSetsTarget(t *testing.T) {
	rules := []lookup.Rule{{
		SourceField: "code",
		TargetField: "label",
		Table:       map[string]string{"200": "OK", "404": "Not Found"},
	}}
	a := lookup.New(rules)
	e := makeEntry(map[string]any{"code": "200"})
	out := a.Apply(e)
	if out.Fields["label"] != "OK" {
		t.Fatalf("expected 'OK', got %v", out.Fields["label"])
	}
}

func TestApply_MissingSourceSkipped(t *testing.T) {
	rules := []lookup.Rule{{
		SourceField: "code",
		TargetField: "label",
		Table:       map[string]string{"200": "OK"},
	}}
	a := lookup.New(rules)
	e := makeEntry(map[string]any{"other": "x"})
	out := a.Apply(e)
	if _, ok := out.Fields["label"]; ok {
		t.Fatal("expected 'label' to be absent")
	}
}

func TestApply_NoMatch_TargetUnchanged(t *testing.T) {
	rules := []lookup.Rule{{
		SourceField: "code",
		TargetField: "label",
		Table:       map[string]string{"200": "OK"},
	}}
	a := lookup.New(rules)
	e := makeEntry(map[string]any{"code": "500"})
	out := a.Apply(e)
	if _, ok := out.Fields["label"]; ok {
		t.Fatal("expected 'label' to remain absent")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	rules := []lookup.Rule{{
		SourceField: "env",
		TargetField: "region",
		Table:       map[string]string{"prod": "us-east-1"},
	}}
	a := lookup.New(rules)
	e := makeEntry(map[string]any{"env": "prod"})
	_ = a.Apply(e)
	if _, ok := e.Fields["region"]; ok {
		t.Fatal("original entry was mutated")
	}
}

func TestParseRules_ValidSpec(t *testing.T) {
	rules, err := lookup.ParseRules([]string{"level:severity=info->low,error->high"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	r := rules[0]
	if r.SourceField != "level" || r.TargetField != "severity" {
		t.Fatalf("unexpected fields: %+v", r)
	}
	if r.Table["info"] != "low" || r.Table["error"] != "high" {
		t.Fatalf("unexpected table: %v", r.Table)
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := lookup.ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil, nil; got %v, %v", rules, err)
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := lookup.ParseRules([]string{"level:severity"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestParseRules_InvalidHeader_ReturnsError(t *testing.T) {
	_, err := lookup.ParseRules([]string{"levelseverity=a->b"})
	if err == nil {
		t.Fatal("expected error for missing colon in header")
	}
}
