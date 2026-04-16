package typecast_test

import (
	"testing"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/typecast"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Fields: fields, Raw: "{}"}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	c := typecast.New(nil)
	e := makeEntry(map[string]any{"status": "200"})
	out := c.Apply(e)
	if out.Fields["status"] != "200" {
		t.Fatalf("expected 200, got %v", out.Fields["status"])
	}
}

func TestApply_CastsStringToInt(t *testing.T) {
	c := typecast.New([]typecast.Rule{{Field: "code", Type: "int"}})
	e := makeEntry(map[string]any{"code": "42"})
	out := c.Apply(e)
	if out.Fields["code"] != int64(42) {
		t.Fatalf("expected int64(42), got %T %v", out.Fields["code"], out.Fields["code"])
	}
}

func TestApply_CastsStringToFloat(t *testing.T) {
	c := typecast.New([]typecast.Rule{{Field: "ratio", Type: "float"}})
	e := makeEntry(map[string]any{"ratio": "3.14"})
	out := c.Apply(e)
	if out.Fields["ratio"] != 3.14 {
		t.Fatalf("expected 3.14, got %v", out.Fields["ratio"])
	}
}

func TestApply_CastsStringToBool(t *testing.T) {
	c := typecast.New([]typecast.Rule{{Field: "ok", Type: "bool"}})
	e := makeEntry(map[string]any{"ok": "true"})
	out := c.Apply(e)
	if out.Fields["ok"] != true {
		t.Fatalf("expected true, got %v", out.Fields["ok"])
	}
}

func TestApply_InvalidCast_LeavesFieldUnchanged(t *testing.T) {
	c := typecast.New([]typecast.Rule{{Field: "num", Type: "int"}})
	e := makeEntry(map[string]any{"num": "not-a-number"})
	out := c.Apply(e)
	if out.Fields["num"] != "not-a-number" {
		t.Fatalf("expected original value, got %v", out.Fields["num"])
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	c := typecast.New([]typecast.Rule{{Field: "missing", Type: "int"}})
	e := makeEntry(map[string]any{"other": "val"})
	out := c.Apply(e)
	if _, ok := out.Fields["missing"]; ok {
		t.Fatal("expected missing field to remain absent")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	c := typecast.New([]typecast.Rule{{Field: "n", Type: "int"}})
	orig := map[string]any{"n": "7"}
	e := makeEntry(orig)
	c.Apply(e)
	if orig["n"] != "7" {
		t.Fatal("input was mutated")
	}
}

func TestParseRules_Valid(t *testing.T) {
	rules, err := typecast.ParseRules([]string{"status=int", "score=float"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 || rules[0].Field != "status" || rules[1].Type != "float" {
		t.Fatalf("unexpected rules: %+v", rules)
	}
}

func TestParseRules_Invalid_ReturnsError(t *testing.T) {
	_, err := typecast.ParseRules([]string{"badspec"})
	if err == nil {
		t.Fatal("expected error for bad spec")
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := typecast.ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil rules and no error")
	}
}
