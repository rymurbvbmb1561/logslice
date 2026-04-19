package lowercase_test

import (
	"testing"
	"time"

	"github.com/your-org/logslice/internal/lowercase"
	"github.com/your-org/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Timestamp: time.Time{}, Fields: fields, Raw: ""}
}

func TestApply_NoFields_ReturnsOriginal(t *testing.T) {
	a := lowercase.New(nil)
	e := makeEntry(map[string]any{"msg": "HELLO"})
	out := a.Apply(e)
	if out.Fields["msg"] != "HELLO" {
		t.Fatalf("expected HELLO, got %v", out.Fields["msg"])
	}
}

func TestApply_LowercasesField(t *testing.T) {
	a := lowercase.New([]string{"msg"})
	e := makeEntry(map[string]any{"msg": "HELLO World"})
	out := a.Apply(e)
	if out.Fields["msg"] != "hello world" {
		t.Fatalf("expected 'hello world', got %v", out.Fields["msg"])
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	a := lowercase.New([]string{"missing"})
	e := makeEntry(map[string]any{"msg": "HELLO"})
	out := a.Apply(e)
	if out.Fields["msg"] != "HELLO" {
		t.Fatalf("expected HELLO unchanged, got %v", out.Fields["msg"])
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	a := lowercase.New([]string{"count"})
	e := makeEntry(map[string]any{"count": 42})
	out := a.Apply(e)
	if out.Fields["count"] != 42 {
		t.Fatalf("expected 42 unchanged, got %v", out.Fields["count"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	a := lowercase.New([]string{"msg"})
	fields := map[string]any{"msg": "UPPER"}
	e := makeEntry(fields)
	a.Apply(e)
	if fields["msg"] != "UPPER" {
		t.Fatal("original entry was mutated")
	}
}

func TestParseFields_ValidSpec(t *testing.T) {
	fields, err := lowercase.ParseFields("msg, level, host")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
}

func TestParseFields_Empty_ReturnsNil(t *testing.T) {
	fields, err := lowercase.ParseFields("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fields != nil {
		t.Fatalf("expected nil, got %v", fields)
	}
}

func TestParseFields_EmptySegment_ReturnsError(t *testing.T) {
	_, err := lowercase.ParseFields("msg,,level")
	if err == nil {
		t.Fatal("expected error for empty segment")
	}
}
