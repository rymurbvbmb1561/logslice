package ceiling_test

import (
	"testing"

	"github.com/logslice/logslice/internal/ceiling"
)

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := ceiling.ParseRules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules != nil {
		t.Fatal("expected nil rules")
	}
}

func TestParseRules_ValidSpec(t *testing.T) {
	rules, err := ceiling.ParseRules([]string{"latency=100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Field != "latency" || rules[0].Multiple != 100 {
		t.Fatalf("unexpected rule: %+v", rules[0])
	}
}

func TestParseRules_MultipleSpecs(t *testing.T) {
	rules, err := ceiling.ParseRules([]string{"latency=100", "size=1024"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := ceiling.ParseRules([]string{"latency100"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseRules_EmptyField_ReturnsError(t *testing.T) {
	_, err := ceiling.ParseRules([]string{"=100"})
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestParseRules_NonPositiveMultiple_ReturnsError(t *testing.T) {
	_, err := ceiling.ParseRules([]string{"latency=0"})
	if err == nil {
		t.Fatal("expected error for zero multiple")
	}
}

func TestParseRules_InvalidMultiple_ReturnsError(t *testing.T) {
	_, err := ceiling.ParseRules([]string{"latency=abc"})
	if err == nil {
		t.Fatal("expected error for non-numeric multiple")
	}
}
