package uppercase_test

import (
	"testing"

	"github.com/your-org/logslice/internal/parser"
	"github.com/your-org/logslice/internal/uppercase"
)

func makeEntry(kvs ...interface{}) parser.Entry {
	e := make(parser.Entry)
	for i := 0; i+1 < len(kvs); i += 2 {
		e[kvs[i].(string)] = kvs[i+1]
	}
	return e
}

func TestApply_NoFields_ReturnsOriginal(t *testing.T) {
	tr := uppercase.New()
	e := makeEntry("msg", "hello world")
	out := tr.Apply(e)
	if out["msg"] != "hello world" {
		t.Errorf("expected unchanged, got %v", out["msg"])
	}
}

func TestApply_UppercasesField(t *testing.T) {
	tr := uppercase.New(uppercase.WithFields([]string{"msg"}))
	e := makeEntry("msg", "hello world")
	out := tr.Apply(e)
	if out["msg"] != "HELLO WORLD" {
		t.Errorf("expected HELLO WORLD, got %v", out["msg"])
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	tr := uppercase.New(uppercase.WithFields([]string{"missing"}))
	e := makeEntry("msg", "hello")
	out := tr.Apply(e)
	if _, ok := out["missing"]; ok {
		t.Error("expected missing field to remain absent")
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	tr := uppercase.New(uppercase.WithFields([]string{"count"}))
	e := makeEntry("count", 42)
	out := tr.Apply(e)
	if out["count"] != 42 {
		t.Errorf("expected 42, got %v", out["count"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	tr := uppercase.New(uppercase.WithFields([]string{"msg"}))
	e := makeEntry("msg", "original")
	tr.Apply(e)
	if e["msg"] != "original" {
		t.Error("input entry was mutated")
	}
}

func TestParseFields_Valid(t *testing.T) {
	fields, err := uppercase.ParseFields("msg, level, host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 3 || fields[0] != "msg" || fields[1] != "level" || fields[2] != "host" {
		t.Errorf("unexpected fields: %v", fields)
	}
}

func TestParseFields_Empty_ReturnsNil(t *testing.T) {
	fields, err := uppercase.ParseFields("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields != nil {
		t.Errorf("expected nil, got %v", fields)
	}
}

func TestParseFields_EmptySegment_ReturnsError(t *testing.T) {
	_, err := uppercase.ParseFields("msg,,level")
	if err == nil {
		t.Error("expected error for empty segment")
	}
}
