package transform

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	e := make(parser.Entry)
	for k, v := range fields {
		e[k] = v
	}
	return e
}

func TestApply_DeleteField(t *testing.T) {
	t.Parallel()
	tr := New([]Rule{{Field: "secret", Action: ActionDelete}})
	out := tr.Apply(makeEntry(map[string]any{"msg": "hello", "secret": "topsecret"}))
	if _, ok := out["secret"]; ok {
		t.Error("expected 'secret' to be deleted")
	}
	if out["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", out["msg"])
	}
}

func TestApply_SetField(t *testing.T) {
	t.Parallel()
	tr := New([]Rule{{Field: "env", Action: ActionSet, Value: "production"}})
	out := tr.Apply(makeEntry(map[string]any{"msg": "hi"}))
	if out["env"] != "production" {
		t.Errorf("expected env=production, got %v", out["env"])
	}
}

func TestApply_RenameField(t *testing.T) {
	t.Parallel()
	tr := New([]Rule{{Field: "message", Action: ActionRename, Value: "msg"}})
	out := tr.Apply(makeEntry(map[string]any{"message": "hello"}))
	if out["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", out["msg"])
	}
	if _, ok := out["message"]; ok {
		t.Error("expected old field 'message' to be removed")
	}
}

func TestApply_RenameNonExistentField(t *testing.T) {
	t.Parallel()
	tr := New([]Rule{{Field: "missing", Action: ActionRename, Value: "new"}})
	out := tr.Apply(makeEntry(map[string]any{"msg": "hi"}))
	if _, ok := out["new"]; ok {
		t.Error("expected no 'new' field when source does not exist")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	t.Parallel()
	original := makeEntry(map[string]any{"a": "1", "b": "2"})
	tr := New([]Rule{{Field: "a", Action: ActionDelete}})
	_ = tr.Apply(original)
	if _, ok := original["a"]; !ok {
		t.Error("original entry should not be mutated")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	t.Parallel()
	rules := []Rule{
		{Field: "level", Action: ActionSet, Value: "info"},
		{Field: "debug", Action: ActionDelete},
	}
	tr := New(rules)
	out := tr.Apply(makeEntry(map[string]any{"debug": "verbose", "msg": "ok"}))
	if out["level"] != "info" {
		t.Errorf("expected level=info, got %v", out["level"])
	}
	if _, ok := out["debug"]; ok {
		t.Error("expected 'debug' field to be deleted")
	}
}
