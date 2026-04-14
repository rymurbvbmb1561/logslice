package highlight

import (
	"testing"
)

func TestParseRules_ValidSpecs(t *testing.T) {
	rules, err := ParseRules([]string{"level=red", "msg=cyan"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Field != "level" || rules[0].Color != Red {
		t.Errorf("rule[0]: got {%q, %q}", rules[0].Field, rules[0].Color)
	}
	if rules[1].Field != "msg" || rules[1].Color != Cyan {
		t.Errorf("rule[1]: got {%q, %q}", rules[1].Field, rules[1].Color)
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := ParseRules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules != nil {
		t.Errorf("expected nil rules, got %v", rules)
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := ParseRules([]string{"levelred"})
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestParseRules_EmptyField_ReturnsError(t *testing.T) {
	_, err := ParseRules([]string{"=red"})
	if err == nil {
		t.Error("expected error for empty field name")
	}
}

func TestParseRules_UnknownColor_ReturnsError(t *testing.T) {
	_, err := ParseRules([]string{"level=magenta"})
	if err == nil {
		t.Error("expected error for unknown color")
	}
}

func TestDefaultRules_NonEmpty(t *testing.T) {
	rules := DefaultRules()
	if len(rules) == 0 {
		t.Error("expected non-empty default rules")
	}
	for _, r := range rules {
		if r.Field == "" {
			t.Error("default rule has empty field")
		}
		if r.Color == Reset {
			t.Errorf("default rule for %q has no color set", r.Field)
		}
	}
}
