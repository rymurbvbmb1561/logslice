package truncator_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/truncator"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Raw:       `{}`,
		Fields:    fields,
	}
}

func TestApply_NoTruncationWhenMaxLenZero(t *testing.T) {
	tr := truncator.New(truncator.Options{MaxLen: 0})
	e := makeEntry(map[string]interface{}{"msg": "hello world"})
	out := tr.Apply(e)
	if got := out.Fields["msg"]; got != "hello world" {
		t.Fatalf("expected unchanged value, got %q", got)
	}
}

func TestApply_TruncatesLongString(t *testing.T) {
	tr := truncator.New(truncator.Options{MaxLen: 5, Ellipsis: "..."})
	e := makeEntry(map[string]interface{}{"msg": "hello world"})
	out := tr.Apply(e)
	if got := out.Fields["msg"]; got != "hello..." {
		t.Fatalf("expected %q, got %q", "hello...", got)
	}
}

func TestApply_ShortStringUnchanged(t *testing.T) {
	tr := truncator.New(truncator.Options{MaxLen: 20})
	e := makeEntry(map[string]interface{}{"msg": "hi"})
	out := tr.Apply(e)
	if got := out.Fields["msg"]; got != "hi" {
		t.Fatalf("expected %q, got %q", "hi", got)
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	tr := truncator.New(truncator.Options{MaxLen: 3})
	e := makeEntry(map[string]interface{}{"count": 42})
	out := tr.Apply(e)
	if got := out.Fields["count"]; got != 42 {
		t.Fatalf("expected 42, got %v", got)
	}
}

func TestApply_DefaultEllipsis(t *testing.T) {
	tr := truncator.New(truncator.Options{MaxLen: 4}) // Ellipsis left empty → "..."
	e := makeEntry(map[string]interface{}{"msg": "abcdefgh"})
	out := tr.Apply(e)
	expected := "abcd..."
	if got := out.Fields["msg"]; got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}

func TestApply_PreservesTimestampAndRaw(t *testing.T) {
	ts := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	tr := truncator.New(truncator.Options{MaxLen: 3})
	e := parser.Entry{Timestamp: ts, Raw: `{"msg":"hi"}`, Fields: map[string]interface{}{"msg": "hi"}}
	out := tr.Apply(e)
	if !out.Timestamp.Equal(ts) {
		t.Fatalf("timestamp mismatch: %v", out.Timestamp)
	}
	if out.Raw != e.Raw {
		t.Fatalf("raw mismatch: %q", out.Raw)
	}
}
