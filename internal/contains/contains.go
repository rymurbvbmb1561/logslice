// Package contains provides a processor that filters log entries
// based on whether a specified field's string value contains a given substring.
package contains

import (
	"fmt"
	"strings"
)

// Entry represents a parsed log line with its raw form and structured fields.
type Entry interface{}

// LogEntry is the concrete entry type used across logslice.
type logEntry struct {
	Fields map[string]interface{}
	Raw    string
}

// Rule defines a single contains check: the source field and the substring to match.
type Rule struct {
	Field     string
	Substring string
	Negate    bool // if true, passes entries that do NOT contain the substring
}

// Processor filters entries based on substring presence in a field.
type Processor struct {
	rules []Rule
}

// New creates a Processor with the given rules.
func New(rules []Rule) *Processor {
	return &Processor{rules: rules}
}

// Apply returns true if the entry passes all contains rules.
// An entry passes a rule when the named field's string value contains
// the specified substring (or does not contain it when Negate is true).
func (p *Processor) Apply(fields map[string]interface{}) bool {
	for _, r := range p.rules {
		val, ok := fields[r.Field]
		if !ok {
			// missing field: treat as empty string
			val = ""
		}
		s, ok := val.(string)
		if !ok {
			s = fmt.Sprintf("%v", val)
		}
		has := strings.Contains(s, r.Substring)
		if r.Negate && has {
			return false
		}
		if !r.Negate && !has {
			return false
		}
	}
	return true
}

// ParseRules parses specs of the form "field=substring" or "!field=substring"
// (leading '!' means negate). Returns nil for an empty slice.
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		negate := false
		if strings.HasPrefix(spec, "!") {
			negate = true
			spec = spec[1:]
		}
		idx := strings.IndexByte(spec, '=')
		if idx < 0 {
			return nil, fmt.Errorf("contains: missing '=' in spec %q", spec)
		}
		field := spec[:idx]
		if field == "" {
			return nil, fmt.Errorf("contains: empty field name in spec %q", spec)
		}
		rules = append(rules, Rule{
			Field:     field,
			Substring: spec[idx+1:],
			Negate:    negate,
		})
	}
	return rules, nil
}
