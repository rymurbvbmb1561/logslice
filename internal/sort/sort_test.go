package sort

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{
		Timestamp: time.Time{},
		Fields:    fields,
		Raw:       "{}",
	}
}

func TestEntries_SortAscending(t *testing.T) {
	s := New("level")
	s.Add(makeEntry(map[string]interface{}{"level": "warn"}))
	s.Add(makeEntry(map[string]interface{}{"level": "error"}))
	s.Add(makeEntry(map[string]interface{}{"level": "info"}))

	out := s.Entries()
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
	if out[0].Fields["level"] != "error" || out[1].Fields["level"] != "info" || out[2].Fields["level"] != "warn" {
		t.Errorf("unexpected order: %v %v %v", out[0].Fields["level"], out[1].Fields["level"], out[2].Fields["level"])
	}
}

func TestEntries_SortDescending(t *testing.T) {
	s := New("level", WithOrder(Descending))
	s.Add(makeEntry(map[string]interface{}{"level": "info"}))
	s.Add(makeEntry(map[string]interface{}{"level": "warn"}))
	s.Add(makeEntry(map[string]interface{}{"level": "error"}))

	out := s.Entries()
	if out[0].Fields["level"] != "warn" {
		t.Errorf("expected warn first, got %v", out[0].Fields["level"])
	}
}

func TestEntries_MissingFieldSortsFirst(t *testing.T) {
	s := New("code")
	s.Add(makeEntry(map[string]interface{}{"code": "500"}))
	s.Add(makeEntry(map[string]interface{}{}))
	s.Add(makeEntry(map[string]interface{}{"code": "200"}))

	out := s.Entries()
	if stringify(out[0].Fields["code"]) != "" {
		t.Errorf("expected missing field first, got %v", out[0].Fields["code"])
	}
}

func TestEntries_DoesNotMutateBuffer(t *testing.T) {
	s := New("msg")
	s.Add(makeEntry(map[string]interface{}{"msg": "b"}))
	s.Add(makeEntry(map[string]interface{}{"msg": "a"}))

	_ = s.Entries()
	if s.entries[0].Fields["msg"] != "b" {
		t.Error("original buffer was mutated")
	}
}

func TestReset_ClearsBuffer(t *testing.T) {
	s := New("msg")
	s.Add(makeEntry(map[string]interface{}{"msg": "hello"}))
	s.Reset()
	if len(s.Entries()) != 0 {
		t.Error("expected empty after reset")
	}
}
