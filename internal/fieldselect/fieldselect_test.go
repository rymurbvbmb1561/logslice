package fieldselect_test

import (
	"testing"

	"github.com/user/logslice/internal/fieldselect"
	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Raw: "{}", Fields: fields}
}

func TestApply_NoFields_ReturnsOriginal(t *testing.T) {
	s := fieldselect.New()
	e := makeEntry(map[string]any{"a": 1, "b": 2})
	out := s.Apply(e)
	if len(out.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out.Fields))
	}
}

func TestApply_KeepsOnlySelectedFields(t *testing.T) {
	s := fieldselect.New(fieldselect.WithFields([]string{"a", "c"}))
	e := makeEntry(map[string]any{"a": 1, "b": 2, "c": 3})
	out := s.Apply(e)
	if len(out.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out.Fields))
	}
	if _, ok := out.Fields["b"]; ok {
		t.Error("field 'b' should have been removed")
	}
}

func TestApply_DropMode_RemovesListedFields(t *testing.T) {
	s := fieldselect.New(fieldselect.WithFields([]string{"secret"}), fieldselect.WithDrop())
	e := makeEntry(map[string]any{"msg": "hello", "secret": "token"})
	out := s.Apply(e)
	if _, ok := out.Fields["secret"]; ok {
		t.Error("field 'secret' should have been dropped")
	}
	if _, ok := out.Fields["msg"]; !ok {
		t.Error("field 'msg' should be present")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	s := fieldselect.New(fieldselect.WithFields([]string{"a"}))
	e := makeEntry(map[string]any{"a": 1, "b": 2})
	s.Apply(e)
	if len(e.Fields) != 2 {
		t.Error("original entry was mutated")
	}
}

func TestApply_MissingFieldIgnored(t *testing.T) {
	s := fieldselect.New(fieldselect.WithFields([]string{"a", "z"}))
	e := makeEntry(map[string]any{"a": 1, "b": 2})
	out := s.Apply(e)
	if len(out.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(out.Fields))
	}
}

func TestParseFields_ValidSpec(t *testing.T) {
	fields := fieldselect.ParseFields("a,b,c")
	if len(fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(fields))
	}
}

func TestParseFields_Empty_ReturnsNil(t *testing.T) {
	fields := fieldselect.ParseFields("")
	if fields != nil {
		t.Error("expected nil for empty spec")
	}
}

func TestParseFields_SkipsEmptySegments(t *testing.T) {
	fields := fieldselect.ParseFields("a,,b")
	if len(fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(fields))
	}
}
