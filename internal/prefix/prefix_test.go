package prefix_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/prefix"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Timestamp: time.Time{}, Fields: fields, Raw: ""}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	p := prefix.New(nil)
	e := makeEntry(map[string]any{"msg": "hello"})
	out := p.Apply(e)
	if out.Fields["msg"] != "hello" {
		t.Fatalf("expected hello, got %v", out.Fields["msg"])
	}
}

func TestApply_PrefixesStringField(t *testing.T) {
	rules := []prefix.Rule{{Field: "msg", Prefix: "[INFO] "}}
	p := prefix.New(rules)
	e := makeEntry(map[string]any{"msg": "started"})
	out := p.Apply(e)
	if out.Fields["msg"] != "[INFO] started" {
		t.Fatalf("unexpected value: %v", out.Fields["msg"])
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	rules := []prefix.Rule{{Field: "missing", Prefix: "x-"}}
	p := prefix.New(rules)
	e := makeEntry(map[string]any{"msg": "hello"})
	out := p.Apply(e)
	if _, ok := out.Fields["missing"]; ok {
		t.Fatal("expected missing field to remain absent")
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	rules := []prefix.Rule{{Field: "count", Prefix: "n="}}
	p := prefix.New(rules)
	e := makeEntry(map[string]any{"count": 42})
	out := p.Apply(e)
	if out.Fields["count"] != 42 {
		t.Fatalf("expected 42, got %v", out.Fields["count"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	rules := []prefix.Rule{{Field: "msg", Prefix: "pre-"}}
	p := prefix.New(rules)
	e := makeEntry(map[string]any{"msg": "original"})
	p.Apply(e)
	if e.Fields["msg"] != "original" {
		t.Fatal("input entry was mutated")
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := prefix.ParseRules([]string{"level=LOG:", "env=prod-"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Field != "level" || rules[0].Prefix != "LOG:" {
		t.Fatalf("unexpected rule: %+v", rules[0])
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := prefix.ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil, nil; got %v, %v", rules, err)
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := prefix.ParseRules([]string{"noequals"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseRules_EmptyField_ReturnsError(t *testing.T) {
	_, err := prefix.ParseRules([]string{"=someprefix"})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}
