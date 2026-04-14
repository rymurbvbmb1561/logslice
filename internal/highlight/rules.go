package highlight

import (
	"fmt"
	"strings"
)

// ParseRules parses a slice of "field=color" strings into Rule values.
// Returns an error if any entry is malformed or the color is unknown.
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("highlight: invalid rule %q, expected field=color", s)
		}
		field := strings.TrimSpace(parts[0])
		colorName := strings.TrimSpace(parts[1])
		c, ok := ParseColor(colorName)
		if !ok {
			return nil, fmt.Errorf("highlight: unknown color %q in rule %q", colorName, s)
		}
		rules = append(rules, Rule{Field: field, Color: c})
	}
	return rules, nil
}

// DefaultRules returns a sensible default set of highlight rules.
func DefaultRules() []Rule {
	return []Rule{
		{Field: "error", Color: Red},
		{Field: "warn", Color: Yellow},
		{Field: "info", Color: Green},
		{Field: "debug", Color: Cyan},
	}
}
