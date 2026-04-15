package coalesce_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/coalesce"
	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       "{}",
		Fields:    fields,
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	c := coalesce.New(nil)
	e := makeEntry(map[string]interface{}{"a": "hello"})
	out := c.Apply(e)
	if out.Fields["a"] != "hello" {
		t.Fatalf("expected 'hello', got %v", out.Fields["a"])
	}
}

func TestApply_FirstSourceUsed(t *testing.T) {
	c := coalesce.New([]coalesce.Rule{
		{Target: "msg", Sources: []string{"message", "log", "text"}},
	})
	e := makeEntry(map[string]interface{}{"message": "hello", "log": "world"})
	out := c.Apply(e)
	if out.Fields["msg"] != "hello" {
		t.Fatalf("expected 'hello', got %v", out.Fields["msg"])
	}
}

func TestApply_FallsBackToSecondSource(t *testing.T) {
	c := coalesce.New([]coalesce.Rule{
		{Target: "msg", Sources: []string{"message", "log"}},
	})
	e := makeEntry(map[string]interface{}{"log": "fallback"})
	out := c.Apply(e)
	if out.Fields["msg"] != "fallback" {
		t.Fatalf("expected 'fallback', got %v", out.Fields["msg"])
	}
}

func TestApply_EmptyStringSkipped(t *testing.T) {
	c := coalesce.New([]coalesce.Rule{
		{Target: "msg", Sources: []string{"a", "b"}},
	})
	e := makeEntry(map[string]interface{}{"a": "   ", "b": "used"})
	out := c.Apply(e)
	if out.Fields["msg"] != "used" {
		t.Fatalf("expected 'used', got %v", out.Fields["msg"])
	}
}

func TestApply_NoSourceFound_TargetNotSet(t *testing.T) {
	c := coalesce.New([]coalesce.Rule{
		{Target: "msg", Sources: []string{"missing"}},
	})
	e := makeEntry(map[string]interface{}{"other": "value"})
	out := c.Apply(e)
	if _, ok := out.Fields["msg"]; ok {
		t.Fatal("expected target not to be set when no source found")
	}
}

func TestApply_OriginalEntryUnmodified(t *testing.T) {
	c := coalesce.New([]coalesce.Rule{
		{Target: "msg", Sources: []string{"log"}},
	})
	origFields := map[string]interface{}{"log": "hello"}
	e := makeEntry(origFields)
	c.Apply(e)
	if _, ok := origFields["msg"]; ok {
		t.Fatal("original fields map should not be modified")
	}
}
