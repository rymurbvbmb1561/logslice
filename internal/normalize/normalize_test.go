package normalize_test

import (
	"testing"

	"logslice/internal/normalize"
	"logslice/internal/parser"
)

func makeEntry(kv ...interface{}) parser.Entry {
	e := make(parser.Entry)
	for i := 0; i+1 < len(kv); i += 2 {
		e[kv[i].(string)] = kv[i+1]
	}
	return e
}

func TestApply_NoFields_ReturnsOriginal(t *testing.T) {
	n := normalize.New()
	e := makeEntry("msg", "hello world")
	out := n.Apply(e)
	if out["msg"] != "hello world" {
		t.Fatalf("expected 'hello world', got %v", out["msg"])
	}
}

func TestApply_CollapsesWhitespace(t *testing.T) {
	n := normalize.New()
	e := makeEntry("msg", "  hello   world  ")
	out := n.Apply(e)
	if out["msg"] != "hello world" {
		t.Fatalf("expected 'hello world', got %v", out["msg"])
	}
}

func TestApply_WithLowercase(t *testing.T) {
	n := normalize.New(normalize.WithLowercase())
	e := makeEntry("level", "  ERROR  ")
	out := n.Apply(e)
	if out["level"] != "error" {
		t.Fatalf("expected 'error', got %v", out["level"])
	}
}

func TestApply_RestrictedToFields(t *testing.T) {
	n := normalize.New(normalize.WithFields([]string{"level"}))
	e := makeEntry("level", "  WARN  ", "msg", "  keep spaces  ")
	out := n.Apply(e)
	if out["level"] != "WARN" {
		t.Fatalf("expected 'WARN', got %v", out["level"])
	}
	if out["msg"] != "  keep spaces  " {
		t.Fatalf("expected original msg, got %v", out["msg"])
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	n := normalize.New()
	e := makeEntry("count", 42)
	out := n.Apply(e)
	if out["count"] != 42 {
		t.Fatalf("expected 42, got %v", out["count"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	n := normalize.New(normalize.WithLowercase())
	e := makeEntry("level", "  ERROR  ")
	_ = n.Apply(e)
	if e["level"] != "  ERROR  " {
		t.Fatal("input entry was mutated")
	}
}

func TestParseFields_ValidCSV(t *testing.T) {
	fields := normalize.ParseFields("level, msg ,service")
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
	if fields[1] != "msg" {
		t.Fatalf("expected 'msg', got %q", fields[1])
	}
}

func TestParseFields_Empty_ReturnsNil(t *testing.T) {
	if normalize.ParseFields("") != nil {
		t.Fatal("expected nil for empty input")
	}
}
