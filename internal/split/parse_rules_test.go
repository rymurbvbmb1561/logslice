package split_test

import (
	"testing"

	"github.com/logslice/logslice/internal/split"
)

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := split.ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil rules and no error, got %v %v", rules, err)
	}
}

func TestParseRules_ValidSpec(t *testing.T) {
	rules, err := split.ParseRules([]string{"addr:host,port|:"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	r := rules[0]
	if r.Source != "addr" {
		t.Errorf("source = %q", r.Source)
	}
	if r.Delim != ":" {
		t.Errorf("delim = %q", r.Delim)
	}
	if len(r.Targets) != 2 || r.Targets[0] != "host" || r.Targets[1] != "port" {
		t.Errorf("targets = %v", r.Targets)
	}
}

func TestParseRules_MultipleSpecs(t *testing.T) {
	specs := []string{"a:x,y|,", "b:p,q| "}
	rules, err := split.ParseRules(specs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}

func TestParseRules_MissingPipe_ReturnsError(t *testing.T) {
	_, err := split.ParseRules([]string{"addr:host,port"})
	if err == nil {
		t.Fatal("expected error for missing pipe")
	}
}

func TestParseRules_EmptyDelim_ReturnsError(t *testing.T) {
	_, err := split.ParseRules([]string{"addr:host,port|"})
	if err == nil {
		t.Fatal("expected error for empty delimiter")
	}
}

func TestParseRules_MissingColon_ReturnsError(t *testing.T) {
	_, err := split.ParseRules([]string{"addr|,"})
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseRules_EmptySource_ReturnsError(t *testing.T) {
	_, err := split.ParseRules([]string{":host,port|,"})
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}
