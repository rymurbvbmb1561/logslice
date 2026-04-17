package distinct_test

import (
	"testing"

	"github.com/logslice/logslice/internal/distinct"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_NoField_AllowsAll(t *testing.T) {
	p := distinct.New()
	for i := 0; i < 5; i++ {
		if !p.Apply(makeEntry(map[string]any{"msg": "hello"})) {
			t.Fatal("expected all entries to pass when no field configured")
		}
	}
}

func TestApply_UniqueValues_AllPass(t *testing.T) {
	p := distinct.New(distinct.WithField("host"))
	values := []string{"a", "b", "c"}
	for _, v := range values {
		if !p.Apply(makeEntry(map[string]any{"host": v})) {
			t.Fatalf("expected entry with host=%q to pass", v)
		}
	}
}

func TestApply_DuplicateValue_Dropped(t *testing.T) {
	p := distinct.New(distinct.WithField("host"))
	e := makeEntry(map[string]any{"host": "web-01"})
	if !p.Apply(e) {
		t.Fatal("first occurrence should pass")
	}
	if p.Apply(e) {
		t.Fatal("second occurrence should be dropped")
	}
}

func TestApply_MissingField_Passes(t *testing.T) {
	p := distinct.New(distinct.WithField("host"))
	e := makeEntry(map[string]any{"msg": "no host"})
	if !p.Apply(e) {
		t.Fatal("entry missing the field should always pass")
	}
	if !p.Apply(e) {
		t.Fatal("repeated entry missing the field should still pass")
	}
}

func TestApply_NumericField_Deduplicated(t *testing.T) {
	p := distinct.New(distinct.WithField("code"))
	e := makeEntry(map[string]any{"code": float64(200)})
	if !p.Apply(e) {
		t.Fatal("first numeric value should pass")
	}
	if p.Apply(e) {
		t.Fatal("duplicate numeric value should be dropped")
	}
}

func TestReset_ClearsSeen(t *testing.T) {
	p := distinct.New(distinct.WithField("id"))
	e := makeEntry(map[string]any{"id": "x"})
	p.Apply(e)
	p.Reset()
	if !p.Apply(e) {
		t.Fatal("entry should pass again after Reset")
	}
}
