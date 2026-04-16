package compute

import (
	"testing"
)

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil rules and nil error, got %v %v", rules, err)
	}
}

func TestParseRules_ValidMultiply(t *testing.T) {
	rules, err := ParseRules([]string{"ms=latency*1000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	r := rules[0]
	if r.Target != "ms" || r.Source != "latency" || r.Op != "*" || r.Operand != 1000 {
		t.Fatalf("unexpected rule: %+v", r)
	}
}

func TestParseRules_MultipleSpecs(t *testing.T) {
	specs := []string{"a=x*2", "b=y+5"}
	rules, err := ParseRules(specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := ParseRules([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseRules_InvalidExpr_ReturnsError(t *testing.T) {
	_, err := ParseRules([]string{"out=justfield"})
	if err == nil {
		t.Fatal("expected error for expr without operator")
	}
}

func TestParseRules_NonNumericOperand_ReturnsError(t *testing.T) {
	_, err := ParseRules([]string{"out=field*notanumber"})
	if err == nil {
		t.Fatal("expected error for non-numeric operand")
	}
}
