package fieldselect_test

import (
	"testing"

	"github.com/user/logslice/internal/fieldselect"
	"github.com/user/logslice/internal/parser"
)

func TestParseAndApply_KeepRoundTrip(t *testing.T) {
	line := `{"time":"2024-01-01T00:00:00Z","level":"info","msg":"hello","secret":"abc"}`
	e, err := parser.ParseLine(line)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	fields := fieldselect.ParseFields("time,level,msg")
	s := fieldselect.New(fieldselect.WithFields(fields))
	out := s.Apply(e)

	if _, ok := out.Fields["secret"]; ok {
		t.Error("secret should have been removed")
	}
	for _, f := range []string{"time", "level", "msg"} {
		if _, ok := out.Fields[f]; !ok {
			t.Errorf("field %q should be present", f)
		}
	}
}

func TestParseAndApply_DropRoundTrip(t *testing.T) {
	line := `{"time":"2024-01-01T00:00:00Z","level":"info","msg":"hello","secret":"abc"}`
	e, err := parser.ParseLine(line)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	fields := fieldselect.ParseFields("secret")
	s := fieldselect.New(fieldselect.WithFields(fields), fieldselect.WithDrop())
	out := s.Apply(e)

	if _, ok := out.Fields["secret"]; ok {
		t.Error("secret should have been dropped")
	}
	if len(out.Fields) != len(e.Fields)-1 {
		t.Errorf("expected %d fields, got %d", len(e.Fields)-1, len(out.Fields))
	}
}

func TestApply_PreservesRaw(t *testing.T) {
	e := parser.Entry{Raw: `{"a":1}`, Fields: map[string]any{"a": 1, "b": 2}}
	s := fieldselect.New(fieldselect.WithFields([]string{"a"}))
	out := s.Apply(e)
	if out.Raw != e.Raw {
		t.Errorf("Raw changed: got %q want %q", out.Raw, e.Raw)
	}
}
