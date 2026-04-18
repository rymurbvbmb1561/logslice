package join_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/join"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Timestamp: time.Now(), Fields: fields, Raw: "{}"}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	j := join.New(nil)
	e := makeEntry(map[string]any{"a": "x"})
	out := j.Apply(e)
	if out.Fields["a"] != "x" {
		t.Fatalf("expected x, got %v", out.Fields["a"])
	}
}

func TestApply_JoinsTwoFields(t *testing.T) {
	rules := []join.Rule{{Target: "full", Sources: []string{"first", "last"}, Separator: " "}}
	j := join.New(rules)
	e := makeEntry(map[string]any{"first": "John", "last": "Doe"})
	out := j.Apply(e)
	if out.Fields["full"] != "John Doe" {
		t.Fatalf("expected 'John Doe', got %v", out.Fields["full"])
	}
}

func TestApply_MissingSourceSkipped(t *testing.T) {
	rules := []join.Rule{{Target: "result", Sources: []string{"a", "b", "c"}, Separator: "-"}}
	j := join.New(rules)
	e := makeEntry(map[string]any{"a": "x", "c": "z"})
	out := j.Apply(e)
	if out.Fields["result"] != "x-z" {
		t.Fatalf("expected 'x-z', got %v", out.Fields["result"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	rules := []join.Rule{{Target: "merged", Sources: []string{"x", "y"}, Separator: ","}}
	j := join.New(rules)
	e := makeEntry(map[string]any{"x": "1", "y": "2"})
	j.Apply(e)
	if _, ok := e.Fields["merged"]; ok {
		t.Fatal("input entry was mutated")
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := join.ParseRules([]string{"full=first,last| "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Target != "full" || rules[0].Separator != " " {
		t.Fatalf("unexpected rule: %+v", rules[0])
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := join.ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil, got %v %v", rules, err)
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := join.ParseRules([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseRules_SingleSource_ReturnsError(t *testing.T) {
	_, err := join.ParseRules([]string{"out=onlyone"})
	if err == nil {
		t.Fatal("expected error for single source")
	}
}
