package template_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
	"github.com/logslice/logslice/internal/template"
)

func makeEntry(fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
		Raw:       "{}",
	}
}

func TestApply_BasicTemplate(t *testing.T) {
	a := template.New("{level} - {msg}")
	e := makeEntry(map[string]any{"level": "info", "msg": "hello"})
	out, err := a.Apply(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "info - hello" {
		t.Errorf("got %q, want %q", out, "info - hello")
	}
}

func TestApply_MissingFieldRendersNil(t *testing.T) {
	a := template.New("{level} - {msg}")
	e := makeEntry(map[string]any{"level": "warn"})
	out, err := a.Apply(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "warn - <nil>" {
		t.Errorf("got %q", out)
	}
}

func TestApply_EmptyTemplate_ReturnsError(t *testing.T) {
	a := template.New("")
	e := makeEntry(map[string]any{})
	_, err := a.Apply(e)
	if err == nil {
		t.Fatal("expected error for empty template")
	}
}

func TestApply_NoPlaceholders_ReturnsLiteral(t *testing.T) {
	a := template.New("static text")
	e := makeEntry(map[string]any{"level": "info"})
	out, err := a.Apply(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "static text" {
		t.Errorf("got %q", out)
	}
}

func TestApply_DuplicatePlaceholder_ReplacedOnce(t *testing.T) {
	a := template.New("{level} {level}")
	e := makeEntry(map[string]any{"level": "error"})
	out, err := a.Apply(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "error error" {
		t.Errorf("got %q", out)
	}
}

func TestApply_NumericField(t *testing.T) {
	a := template.New("code={code}")
	e := makeEntry(map[string]any{"code": 200})
	out, err := a.Apply(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "code=200" {
		t.Errorf("got %q", out)
	}
}
