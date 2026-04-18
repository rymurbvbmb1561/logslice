package trim

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Timestamp: time.Time{}, Fields: fields}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	tr := New(nil)
	e := makeEntry(map[string]any{"msg": "  hello  "})
	out := tr.Apply(e)
	if out.Fields["msg"] != "  hello  " {
		t.Fatalf("expected unchanged, got %q", out.Fields["msg"])
	}
}

func TestApply_TrimsWhitespace(t *testing.T) {
	tr := New([]Rule{{Field: "msg", Dir: Both}})
	e := makeEntry(map[string]any{"msg": "  hello  "})
	out := tr.Apply(e)
	if out.Fields["msg"] != "hello" {
		t.Fatalf("expected 'hello', got %q", out.Fields["msg"])
	}
}

func TestApply_TrimLeft(t *testing.T) {
	tr := New([]Rule{{Field: "msg", Dir: Left}})
	e := makeEntry(map[string]any{"msg": "  hi  "})
	out := tr.Apply(e)
	if out.Fields["msg"] != "hi  " {
		t.Fatalf("expected 'hi  ', got %q", out.Fields["msg"])
	}
}

func TestApply_TrimRight(t *testing.T) {
	tr := New([]Rule{{Field: "msg", Dir: Right}})
	e := makeEntry(map[string]any{"msg": "  hi  "})
	out := tr.Apply(e)
	if out.Fields["msg"] != "  hi" {
		t.Fatalf("expected '  hi', got %q", out.Fields["msg"])
	}
}

func TestApply_CustomCutset(t *testing.T) {
	tr := New([]Rule{{Field: "val", Cutset: "*", Dir: Both}})
	e := makeEntry(map[string]any{"val": "***data***"})
	out := tr.Apply(e)
	if out.Fields["val"] != "data" {
		t.Fatalf("expected 'data', got %q", out.Fields["val"])
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	tr := New([]Rule{{Field: "missing", Dir: Both}})
	e := makeEntry(map[string]any{"msg": "hello"})
	out := tr.Apply(e)
	if _, ok := out.Fields["missing"]; ok {
		t.Fatal("expected missing field to remain absent")
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	tr := New([]Rule{{Field: "count", Dir: Both}})
	e := makeEntry(map[string]any{"count": 42})
	out := tr.Apply(e)
	if out.Fields["count"] != 42 {
		t.Fatalf("expected 42, got %v", out.Fields["count"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	tr := New([]Rule{{Field: "msg", Dir: Both}})
	orig := "  hello  "
	e := makeEntry(map[string]any{"msg": orig})
	tr.Apply(e)
	if e.Fields["msg"] != orig {
		t.Fatal("original entry was mutated")
	}
}
