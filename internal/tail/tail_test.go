package tail_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/tail"
)

func makeEntry(msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       `{"msg":"` + msg + `"}`,
		Fields:    map[string]any{"msg": msg},
	}
}

func TestApply_ZeroMax_KeepsAll(t *testing.T) {
	tr := tail.New()
	for i := 0; i < 10; i++ {
		tr.Record(makeEntry("x"))
	}
	if tr.Len() != 10 {
		t.Fatalf("expected 10, got %d", tr.Len())
	}
}

func TestApply_LimitsToMax(t *testing.T) {
	tr := tail.New(tail.WithMax(3))
	for i := 0; i < 7; i++ {
		tr.Record(makeEntry("x"))
	}
	if tr.Len() != 3 {
		t.Fatalf("expected 3, got %d", tr.Len())
	}
}

func TestApply_ReturnsLastN(t *testing.T) {
	tr := tail.New(tail.WithMax(2))
	tr.Record(makeEntry("first"))
	tr.Record(makeEntry("second"))
	tr.Record(makeEntry("third"))

	entries := tr.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Fields["msg"] != "second" {
		t.Errorf("expected second, got %v", entries[0].Fields["msg"])
	}
	if entries[1].Fields["msg"] != "third" {
		t.Errorf("expected third, got %v", entries[1].Fields["msg"])
	}
}

func TestEntries_DoesNotMutateBuffer(t *testing.T) {
	tr := tail.New(tail.WithMax(3))
	tr.Record(makeEntry("a"))
	out1 := tr.Entries()
	out1[0].Fields["msg"] = "mutated"
	out2 := tr.Entries()
	if out2[0].Fields["msg"] == "mutated" {
		t.Error("Entries() should return a copy, not a reference")
	}
}

func TestWithMax_Zero_IgnoredDefaultsToUnlimited(t *testing.T) {
	tr := tail.New(tail.WithMax(0))
	for i := 0; i < 5; i++ {
		tr.Record(makeEntry("y"))
	}
	if tr.Len() != 5 {
		t.Fatalf("expected 5, got %d", tr.Len())
	}
}
