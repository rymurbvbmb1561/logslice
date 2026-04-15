package transform_test

import (
	"testing"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/transform"
)

func TestParseAndApply_RoundTrip(t *testing.T) {
	t.Parallel()
	specs := []string{
		"delete:password",
		"set:service=logslice",
		"rename:ts=timestamp",
	}
	rules, err := transform.ParseRules(specs)
	if err != nil {
		t.Fatalf("ParseRules: %v", err)
	}
	tr := transform.New(rules)

	entry := parser.Entry{
		"password": "hunter2",
		"ts":       "2024-01-01T00:00:00Z",
		"msg":      "login attempt",
	}
	out := tr.Apply(entry)

	if _, ok := out["password"]; ok {
		t.Error("password should be deleted")
	}
	if out["service"] != "logslice" {
		t.Errorf("expected service=logslice, got %v", out["service"])
	}
	if out["timestamp"] != "2024-01-01T00:00:00Z" {
		t.Errorf("expected timestamp field, got %v", out["timestamp"])
	}
	if _, ok := out["ts"]; ok {
		t.Error("old 'ts' field should be gone after rename")
	}
	if out["msg"] != "login attempt" {
		t.Errorf("msg should be unchanged, got %v", out["msg"])
	}
}

func TestParseAndApply_EmptyRules_EntryUnchanged(t *testing.T) {
	t.Parallel()
	rules, err := transform.ParseRules(nil)
	if err != nil {
		t.Fatalf("ParseRules: %v", err)
	}
	tr := transform.New(rules)
	entry := parser.Entry{"msg": "hello", "level": "info"}
	out := tr.Apply(entry)
	if len(out) != len(entry) {
		t.Errorf("expected same number of fields, got %d", len(out))
	}
}

func TestParseAndApply_Apply_DoesNotMutateInput(t *testing.T) {
	t.Parallel()
	specs := []string{
		"delete:secret",
		"set:env=prod",
	}
	rules, err := transform.ParseRules(specs)
	if err != nil {
		t.Fatalf("ParseRules: %v", err)
	}
	tr := transform.New(rules)

	entry := parser.Entry{"secret": "abc123", "msg": "test"}
	_ = tr.Apply(entry)

	// Verify the original entry was not modified by Apply.
	if _, ok := entry["secret"]; !ok {
		t.Error("Apply mutated the input entry: 'secret' was removed")
	}
	if _, ok := entry["env"]; ok {
		t.Error("Apply mutated the input entry: 'env' was added")
	}
}
