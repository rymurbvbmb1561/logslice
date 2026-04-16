package rename_test

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/rename"
)

func makeEntry(fields map[string]any) parser.Entry {
	e := make(parser.Entry)
	for k, v := range fields {
		e[k] = v
	}
	return e
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	r := rename.New(nil)
	entry := makeEntry(map[string]any{"msg": "hello"})
	out := r.Apply(entry)
	if out["msg"] != "hello" {
		t.Fatalf("expected msg=hello, got %v", out["msg"])
	}
}

func TestApply_RenamesField(t *testing.T) {
	r := rename.New([]rename.Rule{{From: "msg", To: "message"}})
	entry := makeEntry(map[string]any{"msg": "hello", "level": "info"})
	out := r.Apply(entry)
	if _, ok := out["msg"]; ok {
		t.Fatal("old field 'msg' should be removed")
	}
	if out["message"] != "hello" {
		t.Fatalf("expected message=hello, got %v", out["message"])
	}
	if out["level"] != "info" {
		t.Fatal("unrelated field should be unchanged")
	}
}

func TestApply_MissingSourceSkipped(t *testing.T) {
	r := rename.New([]rename.Rule{{From: "ghost", To: "spirit"}})
	entry := makeEntry(map[string]any{"msg": "hi"})
	out := r.Apply(entry)
	if _, ok := out["spirit"]; ok {
		t.Fatal("target field should not be created when source is missing")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	r := rename.New([]rename.Rule{{From: "a", To: "b"}})
	entry := makeEntry(map[string]any{"a": "val"})
	r.Apply(entry)
	if _, ok := entry["a"]; !ok {
		t.Fatal("original entry should not be mutated")
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := rename.ParseRules([]string{"old=new", "foo=bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].From != "old" || rules[0].To != "new" {
		t.Fatalf("unexpected rule: %+v", rules[0])
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := rename.ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil, nil; got %v, %v", rules, err)
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := rename.ParseRules([]string{"badspec"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseRules_EmptyTarget_ReturnsError(t *testing.T) {
	_, err := rename.ParseRules([]string{"field="})
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}
