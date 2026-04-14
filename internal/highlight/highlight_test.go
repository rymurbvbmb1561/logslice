package highlight

import (
	"strings"
	"testing"
)

func TestColorize_WrapsString(t *testing.T) {
	result := Colorize(Red, "hello")
	if !strings.Contains(result, "hello") {
		t.Error("expected original string to be present")
	}
	if !strings.HasPrefix(result, string(Red)) {
		t.Error("expected ANSI prefix")
	}
	if !strings.HasSuffix(result, string(Reset)) {
		t.Error("expected ANSI reset suffix")
	}
}

func TestApply_Disabled_ReturnsOriginal(t *testing.T) {
	h := New(false, []Rule{{Field: "level", Color: Red}})
	line := `{"level":"error","msg":"oops"}`
	if got := h.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	h := New(true, nil)
	line := `{"level":"info"}`
	if got := h.Apply(line); got != line {
		t.Errorf("expected unchanged line, got %q", got)
	}
}

func TestApply_HighlightsField(t *testing.T) {
	h := New(true, []Rule{{Field: "level", Color: Green}})
	line := `{"level":"info","msg":"ok"}`
	result := h.Apply(line)
	if !strings.Contains(result, string(Green)) {
		t.Error("expected green color code in output")
	}
	if !strings.Contains(result, "info") {
		t.Error("expected original value to remain")
	}
}

func TestApply_EmptyFieldSkipped(t *testing.T) {
	h := New(true, []Rule{{Field: "", Color: Red}})
	line := `{"level":"warn"}`
	if got := h.Apply(line); got != line {
		t.Errorf("empty field rule should be skipped, got %q", got)
	}
}

func TestParseColor_KnownColors(t *testing.T) {
	cases := []struct {
		name  string
		want  Color
	}{
		{"red", Red},
		{"green", Green},
		{"yellow", Yellow},
		{"blue", Blue},
		{"cyan", Cyan},
		{"bold", Bold},
	}
	for _, tc := range cases {
		c, ok := ParseColor(tc.name)
		if !ok {
			t.Errorf("ParseColor(%q): expected ok", tc.name)
		}
		if c != tc.want {
			t.Errorf("ParseColor(%q): got %q, want %q", tc.name, c, tc.want)
		}
	}
}

func TestParseColor_Unknown_ReturnsFalse(t *testing.T) {
	_, ok := ParseColor("purple")
	if ok {
		t.Error("expected ok=false for unknown color")
	}
}

func TestParseColor_CaseInsensitive(t *testing.T) {
	c, ok := ParseColor("RED")
	if !ok || c != Red {
		t.Errorf("expected Red, got %q ok=%v", c, ok)
	}
}
