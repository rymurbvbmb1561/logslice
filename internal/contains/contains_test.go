package contains_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/contains"
)

func makeEntry(fields map[string]interface{}) map[string]interface{} {
	return fields
}

func TestApply_NoRules_AllowsAll(t *testing.T) {
	p := contains.New(nil)
	entry := makeEntry(map[string]interface{}{"msg": "hello world"})
	if !p.Apply(entry) {
		t.Fatal("expected entry to pass with no rules")
	}
}

func TestApply_FieldContainsSubstring_Passes(t *testing.T) {
	rules := []contains.Rule{{Field: "msg", Substring: "error"}}
	p := contains.New(rules)
	entry := makeEntry(map[string]interface{}{"msg": "an error occurred"})
	if !p.Apply(entry) {
		t.Fatal("expected entry to pass")
	}
}

func TestApply_FieldMissingSubstring_Drops(t *testing.T) {
	rules := []contains.Rule{{Field: "msg", Substring: "error"}}
	p := contains.New(rules)
	entry := makeEntry(map[string]interface{}{"msg": "all good"})
	if p.Apply(entry) {
		t.Fatal("expected entry to be dropped")
	}
}

func TestApply_NegateRule_DropsWhenContains(t *testing.T) {
	rules := []contains.Rule{{Field: "msg", Substring: "debug", Negate: true}}
	p := contains.New(rules)
	entry := makeEntry(map[string]interface{}{"msg": "debug: starting up"})
	if p.Apply(entry) {
		t.Fatal("expected entry to be dropped by negate rule")
	}
}

func TestApply_NegateRule_PassesWhenAbsent(t *testing.T) {
	rules := []contains.Rule{{Field: "msg", Substring: "debug", Negate: true}}
	p := contains.New(rules)
	entry := makeEntry(map[string]interface{}{"msg": "info: all good"})
	if !p.Apply(entry) {
		t.Fatal("expected entry to pass negate rule")
	}
}

func TestApply_MissingField_TreatedAsEmpty(t *testing.T) {
	rules := []contains.Rule{{Field: "level", Substring: "error"}}
	p := contains.New(rules)
	entry := makeEntry(map[string]interface{}{"msg": "something"})
	if p.Apply(entry) {
		t.Fatal("expected entry to be dropped when field missing and substring non-empty")
	}
}

func TestParseRules_ValidSpec(t *testing.T) {
	rules, err := contains.ParseRules([]string{"msg=error"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Field != "msg" || rules[0].Substring != "error" {
		t.Fatalf("unexpected rules: %+v", rules)
	}
}

func TestParseRules_NegatePrefix(t *testing.T) {
	rules, err := contains.ParseRules([]string{"!level=debug"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rules[0].Negate || rules[0].Field != "level" || rules[0].Substring != "debug" {
		t.Fatalf("unexpected rules: %+v", rules)
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := contains.ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil rules and no error, got %v / %v", rules, err)
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := contains.ParseRules([]string{"msgnoequals"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseRules_EmptyField_ReturnsError(t *testing.T) {
	_, err := contains.ParseRules([]string{"=value"})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}
