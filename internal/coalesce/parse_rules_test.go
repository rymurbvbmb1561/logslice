package coalesce_test

import (
	"testing"

	"github.com/user/logslice/internal/coalesce"
)

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := coalesce.ParseRules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules != nil {
		t.Fatal("expected nil rules")
	}
}

func TestParseRules_ValidSpec(t *testing.T) {
	rules, err := coalesce.ParseRules([]string{"msg=message,log,text"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Target != "msg" {
		t.Errorf("expected target 'msg', got %q", rules[0].Target)
	}
	if len(rules[0].Sources) != 3 {
		t.Errorf("expected 3 sources, got %d", len(rules[0].Sources))
	}
}

func TestParseRules_MultipleSpecs(t *testing.T) {
	rules, err := coalesce.ParseRules([]string{"msg=message,log", "host=hostname,host_name"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := coalesce.ParseRules([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseRules_EmptyTarget_ReturnsError(t *testing.T) {
	_, err := coalesce.ParseRules([]string{"=src1,src2"})
	if err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestParseRules_EmptySources_ReturnsError(t *testing.T) {
	_, err := coalesce.ParseRules([]string{"target="})
	if err == nil {
		t.Fatal("expected error for empty sources")
	}
}
