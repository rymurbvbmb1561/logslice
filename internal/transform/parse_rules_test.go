package transform

import (
	"testing"
)

func TestParseRules_Delete(t *testing.T) {
	t.Parallel()
	rules, err := ParseRules([]string{"delete:secret"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Action != ActionDelete || rules[0].Field != "secret" {
		t.Errorf("unexpected rule: %+v", rules)
	}
}

func TestParseRules_Set(t *testing.T) {
	t.Parallel()
	rules, err := ParseRules([]string{"set:env=production"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Field != "env" || rules[0].Value != "production" {
		t.Errorf("unexpected rule: %+v", rules)
	}
}

func TestParseRules_Rename(t *testing.T) {
	t.Parallel()
	rules, err := ParseRules([]string{"rename:message=msg"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Field != "message" || rules[0].Value != "msg" {
		t.Errorf("unexpected rule: %+v", rules)
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	t.Parallel()
	rules, err := ParseRules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules != nil {
		t.Errorf("expected nil rules, got %v", rules)
	}
}

func TestParseRules_UnknownAction_ReturnsError(t *testing.T) {
	t.Parallel()
	_, err := ParseRules([]string{"uppercase:field"})
	if err == nil {
		t.Error("expected error for unknown action")
	}
}

func TestParseRules_MissingColon_ReturnsError(t *testing.T) {
	t.Parallel()
	_, err := ParseRules([]string{"deletefield"})
	if err == nil {
		t.Error("expected error for missing colon")
	}
}

func TestParseRules_DeleteEmptyField_ReturnsError(t *testing.T) {
	t.Parallel()
	_, err := ParseRules([]string{"delete:"})
	if err == nil {
		t.Error("expected error for empty field name")
	}
}

func TestParseRules_SetMissingEquals_ReturnsError(t *testing.T) {
	t.Parallel()
	_, err := ParseRules([]string{"set:noequals"})
	if err == nil {
		t.Error("expected error for missing equals in set")
	}
}
