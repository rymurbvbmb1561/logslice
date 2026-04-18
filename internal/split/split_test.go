package split_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/split"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
		Raw:       "{}",
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	s := split.New(nil)
	e := makeEntry(map[string]any{"msg": "hello world"})
	out := s.Apply(e)
	if out.Fields["msg"] != "hello world" {
		t.Fatalf("expected original entry")
	}
}

func TestApply_SplitsField(t *testing.T) {
	rules := []split.Rule{{Source: "addr", Targets: []string{"host", "port"}, Delim: ":"}}
	s := split.New(rules)
	e := makeEntry(map[string]any{"addr": "localhost:8080"})
	out := s.Apply(e)
	if out.Fields["host"] != "localhost" {
		t.Errorf("host = %v", out.Fields["host"])
	}
	if out.Fields["port"] != "8080" {
		t.Errorf("port = %v", out.Fields["port"])
	}
	if out.Fields["addr"] != "localhost:8080" {
		t.Errorf("source field should be preserved")
	}
}

func TestApply_MissingSourceSkipped(t *testing.T) {
	rules := []split.Rule{{Source: "missing", Targets: []string{"a", "b"}, Delim: ","}}
	s := split.New(rules)
	e := makeEntry(map[string]any{"x": "y"})
	out := s.Apply(e)
	if _, ok := out.Fields["a"]; ok {
		t.Error("should not set target when source missing")
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	rules := []split.Rule{{Source: "count", Targets: []string{"a", "b"}, Delim: ","}}
	s := split.New(rules)
	e := makeEntry(map[string]any{"count": 42})
	out := s.Apply(e)
	if _, ok := out.Fields["a"]; ok {
		t.Error("should not split non-string field")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	rules := []split.Rule{{Source: "kv", Targets: []string{"k", "v"}, Delim: "="}}
	s := split.New(rules)
	e := makeEntry(map[string]any{"kv": "foo=bar"})
	s.Apply(e)
	if _, ok := e.Fields["k"]; ok {
		t.Error("original entry should not be mutated")
	}
}

func TestApply_FewerPartsThanTargets(t *testing.T) {
	rules := []split.Rule{{Source: "data", Targets: []string{"a", "b", "c"}, Delim: ","}}
	s := split.New(rules)
	e := makeEntry(map[string]any{"data": "x,y"})
	out := s.Apply(e)
	if out.Fields["a"] != "x" || out.Fields["b"] != "y" {
		t.Errorf("unexpected values: %v", out.Fields)
	}
	if _, ok := out.Fields["c"]; ok {
		t.Error("c should not be set")
	}
}
