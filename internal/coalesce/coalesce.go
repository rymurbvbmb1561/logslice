// Package coalesce provides a processor that merges a list of candidate
// fields into a single target field, using the first non-empty value found.
package coalesce

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Rule describes a single coalesce operation: pick the first non-empty value
// from Sources and write it to Target.
type Rule struct {
	Target  string
	Sources []string
}

// Coalescer applies a set of coalesce rules to log entries.
type Coalescer struct {
	rules []Rule
}

// New returns a Coalescer configured with the given rules.
func New(rules []Rule) *Coalescer {
	return &Coalescer{rules: rules}
}

// Apply walks each rule, finds the first non-empty source field value, and
// writes it to the target field. The original entry is not modified; a shallow
// copy of its Fields map is returned.
func (c *Coalescer) Apply(entry parser.Entry) parser.Entry {
	if len(c.rules) == 0 {
		return entry
	}

	fields := make(map[string]interface{}, len(entry.Fields))
	for k, v := range entry.Fields {
		fields[k] = v
	}

	for _, rule := range c.rules {
		for _, src := range rule.Sources {
			v, ok := fields[src]
			if !ok {
				continue
			}
			s, ok := v.(string)
			if !ok || strings.TrimSpace(s) == "" {
				continue
			}
			fields[rule.Target] = s
			break
		}
	}

	return parser.Entry{Timestamp: entry.Timestamp, Raw: entry.Raw, Fields: fields}
}

// ParseRules parses a slice of spec strings of the form
// "target=src1,src2,..." into Rule values.
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("coalesce: invalid spec %q: expected target=src1,src2", spec)
		}
		target := strings.TrimSpace(parts[0])
		if target == "" {
			return nil, fmt.Errorf("coalesce: invalid spec %q: target must not be empty", spec)
		}
		rawSrcs := strings.Split(parts[1], ",")
		var sources []string
		for _, s := range rawSrcs {
			s = strings.TrimSpace(s)
			if s != "" {
				sources = append(sources, s)
			}
		}
		if len(sources) == 0 {
			return nil, fmt.Errorf("coalesce: invalid spec %q: at least one source required", spec)
		}
		rules = append(rules, Rule{Target: target, Sources: sources})
	}
	return rules, nil
}
