package extract_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/extract"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Time{},
		Fields:    fields,
		Raw:       "",
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	e := extract.New(nil)
	entry := makeEntry(map[string]any{"msg": "hello world"})
	out := e.Apply(entry)
	if len(out.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(out.Fields))
	}
}

func TestApply_ExtractsNamedGroups(t *testing.T) {
	pat := regexp.MustCompile(`(?P<level>\w+) (?P<msg>.+)`)
	e := extract.New([]extract.Rule{{Source: "text", Pattern: pat}})
	entry := makeEntry(map[string]any{"text": "ERROR something failed"})
	out := e.Apply(entry)
	if out.Fields["level"] != "ERROR" {
		t.Errorf("expected level=ERROR, got %v", out.Fields["level"])
	}
	if out.Fields["msg"] != "something failed" {
		t.Errorf("expected msg='something failed', got %v", out.Fields["msg"])
	}
}

func TestApply_MissingSourceSkipped(t *testing.T) {
	pat := regexp.MustCompile(`(?P<x>\d+)`)
	e := extract.New([]extract.Rule{{Source: "missing", Pattern: pat}})
	entry := makeEntry(map[string]any{"other": "123"})
	out := e.Apply(entry)
	if _, ok := out.Fields["x"]; ok {
		t.Error("expected x not to be set")
	}
}

func TestApply_NoMatch_FieldsUnchanged(t *testing.T) {
	pat := regexp.MustCompile(`(?P<num>\d+)`)
	e := extract.New([]extract.Rule{{Source: "msg", Pattern: pat}})
	entry := makeEntry(map[string]any{"msg": "no digits here"})
	out := e.Apply(entry)
	if _, ok := out.Fields["num"]; ok {
		t.Error("expected num not to be set")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	pat := regexp.MustCompile(`(?P<code>\d+)`)
	e := extract.New([]extract.Rule{{Source: "msg", Pattern: pat}})
	entry := makeEntry(map[string]any{"msg": "code 42"})
	e.Apply(entry)
	if _, ok := entry.Fields["code"]; ok {
		t.Error("original entry should not be mutated")
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := extract.ParseRules([]string{`msg=(?P<level>\w+)`})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Source != "msg" {
		t.Errorf("expected source=msg, got %s", rules[0].Source)
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := extract.ParseRules([]string{"noequals"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestParseRules_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := extract.ParseRules([]string{"msg=("})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := extract.ParseRules(nil)
	if err != nil || rules != nil {
		t.Errorf("expected nil, nil; got %v, %v", rules, err)
	}
}
