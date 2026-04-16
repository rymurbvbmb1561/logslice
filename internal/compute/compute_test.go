package compute

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Fields: fields, Raw: "{}", Timestamp: time.Time{}}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	a := New(nil)
	e := makeEntry(map[string]any{"x": 1.0})
	out := a.Apply(e)
	if len(out.Fields) != 1 {
		t.Fatalf("expected 1 field, got %d", len(out.Fields))
	}
}

func TestApply_Multiply(t *testing.T) {
	rules, _ := ParseRules([]string{"ms=latency*1000"})
	a := New(rules)
	e := makeEntry(map[string]any{"latency": 1.5})
	out := a.Apply(e)
	if out.Fields["ms"] != 1500.0 {
		t.Fatalf("expected 1500, got %v", out.Fields["ms"])
	}
}

func TestApply_Divide(t *testing.T) {
	rules, _ := ParseRules([]string{"s=ms/1000"})
	a := New(rules)
	e := makeEntry(map[string]any{"ms": 2000.0})
	out := a.Apply(e)
	if out.Fields["s"] != 2.0 {
		t.Fatalf("expected 2, got %v", out.Fields["s"])
	}
}

func TestApply_Add(t *testing.T) {
	rules, _ := ParseRules([]string{"total=base+10"})
	a := New(rules)
	e := makeEntry(map[string]any{"base": 5.0})
	out := a.Apply(e)
	if out.Fields["total"] != 15.0 {
		t.Fatalf("expected 15, got %v", out.Fields["total"])
	}
}

func TestApply_MissingSourceSkipped(t *testing.T) {
	rules, _ := ParseRules([]string{"out=missing*2"})
	a := New(rules)
	e := makeEntry(map[string]any{"x": 1.0})
	out := a.Apply(e)
	if _, ok := out.Fields["out"]; ok {
		t.Fatal("expected field to be absent")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	rules, _ := ParseRules([]string{"y=x*2"})
	a := New(rules)
	e := makeEntry(map[string]any{"x": 3.0})
	a.Apply(e)
	if _, ok := e.Fields["y"]; ok {
		t.Fatal("input entry was mutated")
	}
}

func TestApply_DivideByZeroSkipped(t *testing.T) {
	rules, _ := ParseRules([]string{"r=x/0"})
	a := New(rules)
	e := makeEntry(map[string]any{"x": 5.0})
	out := a.Apply(e)
	if _, ok := out.Fields["r"]; ok {
		t.Fatal("expected field to be absent on divide by zero")
	}
}
