package highlight

import (
	"fmt"
	"strings"
)

// Color represents an ANSI terminal color code.
type Color string

const (
	Reset  Color = "\033[0m"
	Red    Color = "\033[31m"
	Green  Color = "\033[32m"
	Yellow Color = "\033[33m"
	Blue   Color = "\033[34m"
	Cyan   Color = "\033[36m"
	Bold   Color = "\033[1m"
)

// Highlighter applies color highlighting to field values in log output.
type Highlighter struct {
	enabled bool
	rules   []Rule
}

// Rule maps a field name to a color.
type Rule struct {
	Field string
	Color Color
}

// New creates a Highlighter. If enabled is false, Apply is a no-op.
func New(enabled bool, rules []Rule) *Highlighter {
	return &Highlighter{enabled: enabled, rules: rules}
}

// Colorize wraps s with the given ANSI color codes.
func Colorize(c Color, s string) string {
	return fmt.Sprintf("%s%s%s", c, s, Reset)
}

// Apply highlights occurrences of each rule's field value within line.
// It operates on the raw string representation of a log line.
func (h *Highlighter) Apply(line string) string {
	if !h.enabled || len(h.rules) == 0 {
		return line
	}
	for _, r := range h.rules {
		if r.Field == "" {
			continue
		}
		// Highlight the key name itself.
		colored := Colorize(r.Color, r.Field)
		line = strings.ReplaceAll(line, r.Field, colored)
	}
	return line
}

// ParseColor converts a color name string to a Color constant.
// Returns Reset and false if the name is unrecognised.
func ParseColor(name string) (Color, bool) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "red":
		return Red, true
	case "green":
		return Green, true
	case "yellow":
		return Yellow, true
	case "blue":
		return Blue, true
	case "cyan":
		return Cyan, true
	case "bold":
		return Bold, true
	}
	return Reset, false
}
