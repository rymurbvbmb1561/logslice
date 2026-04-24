// Package wrap provides a processor that wraps a log entry's fields
// under a new top-level key, nesting the selected fields inside a map.
package wrap

import (
	"fmt"
	"strings"
)

// Entry represents a parsed log line.
type Entry struct {
	Fields map[string]any
	Raw    string
}

// Rule describes a single wrap operation: which fields to nest and
// the destination key to place them under.
type Rule struct {
	// Fields is the list of source field names to nest.
	Fields []string
	// Target is the key under which the nested map will be stored.
	Target string
	// Drop removes the original top-level fields after wrapping when true.
	Drop bool
}

// Wrapper applies wrap rules to log entries.
type Wrapper struct {
	rules []Rule
}

// New returns a Wrapper configured with the provided rules.
func New(rules []Rule) *Wrapper {
	return &Wrapper{rules: rules}
}

// Apply returns a copy of e with the wrap rules applied.
// If no rules are configured the original entry is returned unchanged.
func (w *Wrapper) Apply(e Entry) Entry {
	if len(w.rules) == 0 {
		return e
	}
	out := make(map[string]any, len(e.Fields))
	for k, v := range e.Fields {
		out[k] = v
	}
	for _, r := range w.rules {
		nested := make(map[string]any, len(r.Fields))
		for _, f := range r.Fields {
			if v, ok := out[f]; ok {
				nested[f] = v
			}
		}
		if len(nested) == 0 {
			continue
		}
		out[r.Target] = nested
		if r.Drop {
			for _, f := range r.Fields {
				delete(out, f)
			}
		}
	}
	return Entry{Fields: out, Raw: e.Raw}
}

// ParseRules parses wrap rule specifications of the form:
//
//	"target=field1,field2[+drop]"
//
// The optional "+drop" suffix causes source fields to be removed
// from the top level after wrapping.
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		eq := strings.IndexByte(s, '=')
		if eq < 1 {
			return nil, fmt.Errorf("wrap: missing '=' in rule %q", s)
		}
		target := strings.TrimSpace(s[:eq])
		if target == "" {
			return nil, fmt.Errorf("wrap: empty target in rule %q", s)
		}
		rest := s[eq+1:]
		drop := false
		if strings.HasSuffix(rest, "+drop") {
			drop = true
			rest = strings.TrimSuffix(rest, "+drop")
		}
		parts := strings.Split(rest, ",")
		var fields []string
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			fields = append(fields, p)
		}
		if len(fields) == 0 {
			return nil, fmt.Errorf("wrap: no fields specified in rule %q", s)
		}
		rules = append(rules, Rule{Target: target, Fields: fields, Drop: drop})
	}
	return rules, nil
}
