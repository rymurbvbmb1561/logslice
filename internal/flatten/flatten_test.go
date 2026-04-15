package flatten_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/flatten"
)

func TestApply_FlatEntryUnchanged(t *testing.T) {
	f := flatten.New()
	input := map[string]any{"level": "info", "msg": "hello"}
	out := f.Apply(input)
	if out["level"] != "info" || out["msg"] != "hello" {
		t.Errorf("expected flat entry to be unchanged, got %v", out)
	}
}

func TestApply_NestedMapFlattened(t *testing.T) {
	f := flatten.New()
	input := map[string]any{
		"level": "error",
		"http": map[string]any{
			"method": "GET",
			"status": 200,
		},
	}
	out := f.Apply(input)
	if out["http.method"] != "GET" {
		t.Errorf("expected http.method=GET, got %v", out["http.method"])
	}
	if out["http.status"] != 200 {
		t.Errorf("expected http.status=200, got %v", out["http.status"])
	}
	if _, ok := out["http"]; ok {
		t.Error("expected nested key 'http' to be removed after flattening")
	}
}

func TestApply_DeeplyNestedMap(t *testing.T) {
	f := flatten.New()
	input := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	out := f.Apply(input)
	if out["a.b.c"] != "deep" {
		t.Errorf("expected a.b.c=deep, got %v", out["a.b.c"])
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	f := flatten.New(flatten.WithSeparator("_"))
	input := map[string]any{
		"http": map[string]any{"method": "POST"},
	}
	out := f.Apply(input)
	if out["http_method"] != "POST" {
		t.Errorf("expected http_method=POST, got %v", out["http_method"])
	}
}

func TestApply_MaxDepthLimitsFlattening(t *testing.T) {
	f := flatten.New(flatten.WithMaxDepth(1))
	input := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": "deep",
			},
		},
	}
	out := f.Apply(input)
	// At depth 1, "a.b" should hold the nested map as-is
	if _, ok := out["a.b"]; !ok {
		t.Errorf("expected a.b to exist as a nested map, got keys: %v", out)
	}
	if _, ok := out["a.b.c"]; ok {
		t.Error("expected a.b.c NOT to exist when maxDepth=1")
	}
}

func TestHasPrefix_Match(t *testing.T) {
	entry := map[string]any{"http.method": "GET", "level": "info"}
	if !flatten.HasPrefix(entry, "http") {
		t.Error("expected HasPrefix to return true for 'http'")
	}
}

func TestHasPrefix_NoMatch(t *testing.T) {
	entry := map[string]any{"level": "info", "msg": "ok"}
	if flatten.HasPrefix(entry, "http") {
		t.Error("expected HasPrefix to return false")
	}
}
