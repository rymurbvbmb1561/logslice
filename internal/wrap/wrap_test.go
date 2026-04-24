package wrap

import (
	"testing"
)

func makeEntry(fields map[string]any) Entry {
	return Entry{Fields: fields, Raw: "{}"}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	w := New(nil)
	e := makeEntry(map[string]any{"a": "1", "b": "2"})
	out := w.Apply(e)
	if len(out.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out.Fields))
	}
}

func TestApply_WrapsFields(t *testing.T) {
	rules := []Rule{{Target: "meta", Fields: []string{"host", "pid"}, Drop: false}}
	w := New(rules)
	e := makeEntry(map[string]any{"host": "srv1", "pid": 42, "msg": "hello"})
	out := w.Apply(e)

	nested, ok := out.Fields["meta"].(map[string]any)
	if !ok {
		t.Fatal("expected 'meta' to be a map")
	}
	if nested["host"] != "srv1" {
		t.Errorf("expected host=srv1, got %v", nested["host"])
	}
	if nested["pid"] != 42 {
		t.Errorf("expected pid=42, got %v", nested["pid"])
	}
	// original fields still present when Drop=false
	if _, exists := out.Fields["host"]; !exists {
		t.Error("expected original 'host' field to remain")
	}
}

func TestApply_DropRemovesSourceFields(t *testing.T) {
	rules := []Rule{{Target: "meta", Fields: []string{"host", "pid"}, Drop: true}}
	w := New(rules)
	e := makeEntry(map[string]any{"host": "srv1", "pid": 42, "msg": "hello"})
	out := w.Apply(e)

	if _, exists := out.Fields["host"]; exists {
		t.Error("expected 'host' to be removed after drop")
	}
	if _, exists := out.Fields["pid"]; exists {
		t.Error("expected 'pid' to be removed after drop")
	}
	if _, exists := out.Fields["msg"]; !exists {
		t.Error("expected 'msg' to remain")
	}
}

func TestApply_MissingFieldsSkipped(t *testing.T) {
	rules := []Rule{{Target: "meta", Fields: []string{"host", "missing"}, Drop: false}}
	w := New(rules)
	e := makeEntry(map[string]any{"host": "srv1"})
	out := w.Apply(e)

	nested, ok := out.Fields["meta"].(map[string]any)
	if !ok {
		t.Fatal("expected 'meta' to be a map")
	}
	if _, exists := nested["missing"]; exists {
		t.Error("expected missing field to be absent from nested map")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	rules := []Rule{{Target: "meta", Fields: []string{"host"}, Drop: true}}
	w := New(rules)
	orig := map[string]any{"host": "srv1", "msg": "hi"}
	e := makeEntry(orig)
	w.Apply(e)
	if _, ok := orig["host"]; !ok {
		t.Error("Apply mutated the original entry fields")
	}
}

func TestParseRules_ValidSpec(t *testing.T) {
	rules, err := ParseRules([]string{"meta=host,pid"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Target != "meta" {
		t.Errorf("expected target=meta, got %s", rules[0].Target)
	}
	if len(rules[0].Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(rules[0].Fields))
	}
}

func TestParseRules_DropSuffix(t *testing.T) {
	rules, err := ParseRules([]string{"meta=host,pid+drop"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rules[0].Drop {
		t.Error("expected Drop=true")
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := ParseRules(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rules != nil {
		t.Error("expected nil rules for empty input")
	}
}

func TestParseRules_MissingEquals_ReturnsError(t *testing.T) {
	_, err := ParseRules([]string{"noequalssign"})
	if err == nil {
		t.Error("expected error for missing '='")
	}
}

func TestParseRules_EmptyTarget_ReturnsError(t *testing.T) {
	_, err := ParseRules([]string{"=host,pid"})
	if err == nil {
		t.Error("expected error for empty target")
	}
}

func TestParseRules_NoFields_ReturnsError(t *testing.T) {
	_, err := ParseRules([]string{"meta="})
	if err == nil {
		t.Error("expected error for no fields")
	}
}
