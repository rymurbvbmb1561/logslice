package replace_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/replace"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Timestamp: time.Time{}, Fields: fields, Raw: ""}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	r := replace.New(nil)
	e := makeEntry(map[string]any{"msg": "hello world"})
	out := r.Apply(e)
	if out.Fields["msg"] != "hello world" {
		t.Fatalf("expected unchanged, got %v", out.Fields["msg"])
	}
}

func TestApply_ReplacesMatchingSubstring(t *testing.T) {
	rules, _ := replace.ParseRules([]string{"msg/world/Go"})
	r := replace.New(rules)
	e := makeEntry(map[string]any{"msg": "hello world"})
	out := r.Apply(e)
	if out.Fields["msg"] != "hello Go" {
		t.Fatalf("expected 'hello Go', got %v", out.Fields["msg"])
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	rules, _ := replace.ParseRules([]string{"missing/foo/bar"})
	r := replace.New(rules)
	e := makeEntry(map[string]any{"msg": "hello"})
	out := r.Apply(e)
	if _, ok := out.Fields["missing"]; ok {
		t.Fatal("expected missing field to remain absent")
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	rules, _ := replace.ParseRules([]string{"count/1/2"})
	r := replace.New(rules)
	e := makeEntry(map[string]any{"count": 42})
	out := r.Apply(e)
	if out.Fields["count"] != 42 {
		t.Fatalf("expected 42, got %v", out.Fields["count"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	rules, _ := replace.ParseRules([]string{"msg/hello/hi"})
	r := replace.New(rules)
	e := makeEntry(map[string]any{"msg": "hello world"})
	r.Apply(e)
	if e.Fields["msg"] != "hello world" {
		t.Fatal("original entry was mutated")
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := replace.ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil rules and no error, got %v %v", rules, err)
	}
}

func TestParseRules_InvalidSpec_ReturnsError(t *testing.T) {
	_, err := replace.ParseRules([]string{"badspec"})
	if err == nil {
		t.Fatal("expected error for bad spec")
	}
}

func TestParseRules_InvalidRegex_ReturnsError(t *testing.T) {
	_, err := replace.ParseRules([]string{"field/[invalid/x"})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}
