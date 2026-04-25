package ceiling_test

import (
	"testing"

	"github.com/logslice/logslice/internal/ceiling"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	a := ceiling.New(nil)
	e := makeEntry(map[string]any{"latency": float64(123)})
	out := a.Apply(e)
	if out.Fields["latency"] != float64(123) {
		t.Fatalf("expected 123, got %v", out.Fields["latency"])
	}
}

func TestApply_CeilsToMultiple(t *testing.T) {
	a := ceiling.New([]ceiling.Rule{{Field: "latency", Multiple: 100}})
	e := makeEntry(map[string]any{"latency": float64(123)})
	out := a.Apply(e)
	if out.Fields["latency"] != float64(200) {
		t.Fatalf("expected 200, got %v", out.Fields["latency"])
	}
}

func TestApply_AlreadyOnBoundary_Unchanged(t *testing.T) {
	a := ceiling.New([]ceiling.Rule{{Field: "latency", Multiple: 100}})
	e := makeEntry(map[string]any{"latency": float64(200)})
	out := a.Apply(e)
	if out.Fields["latency"] != float64(200) {
		t.Fatalf("expected 200, got %v", out.Fields["latency"])
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	a := ceiling.New([]ceiling.Rule{{Field: "missing", Multiple: 10}})
	e := makeEntry(map[string]any{"other": float64(5)})
	out := a.Apply(e)
	if _, ok := out.Fields["missing"]; ok {
		t.Fatal("expected missing field to remain absent")
	}
}

func TestApply_NonNumericFieldUnchanged(t *testing.T) {
	a := ceiling.New([]ceiling.Rule{{Field: "level", Multiple: 10}})
	e := makeEntry(map[string]any{"level": "info"})
	out := a.Apply(e)
	if out.Fields["level"] != "info" {
		t.Fatalf("expected 'info', got %v", out.Fields["level"])
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	a := ceiling.New([]ceiling.Rule{{Field: "val", Multiple: 50}})
	e := makeEntry(map[string]any{"val": float64(33)})
	_ = a.Apply(e)
	if e.Fields["val"] != float64(33) {
		t.Fatal("input entry was mutated")
	}
}

func TestApply_StringNumericField(t *testing.T) {
	a := ceiling.New([]ceiling.Rule{{Field: "score", Multiple: 5}})
	e := makeEntry(map[string]any{"score": "13"})
	out := a.Apply(e)
	if out.Fields["score"] != float64(15) {
		t.Fatalf("expected 15, got %v", out.Fields["score"])
	}
}
