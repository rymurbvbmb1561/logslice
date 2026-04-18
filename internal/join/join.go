package join

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule defines how to join multiple fields into a single target field.
type Rule struct {
	Target    string
	Sources   []string
	Separator string
}

// Joiner combines multiple log entry fields into a single field.
type Joiner struct {
	rules []Rule
}

// New returns a Joiner with the given rules.
func New(rules []Rule) *Joiner {
	return &Joiner{rules: rules}
}

// Apply returns a new entry with join rules applied.
func (j *Joiner) Apply(entry parser.Entry) parser.Entry {
	if len(j.rules) == 0 {
		return entry
	}
	out := entry.Clone()
	for _, r := range j.rules {
		parts := make([]string, 0, len(r.Sources))
		for _, src := range r.Sources {
			if v, ok := out.Fields[src]; ok {
				parts = append(parts, fmt.Sprintf("%v", v))
			}
		}
		out.Fields[r.Target] = strings.Join(parts, r.Separator)
	}
	return out
}

// ParseRules parses specs of the form "target=src1,src2[|sep]".
// Example: "full_name=first,last| "
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		eqIdx := strings.Index(spec, "=")
		if eqIdx < 1 {
			return nil, fmt.Errorf("join: invalid spec %q: missing '='" , spec)
		}
		target := spec[:eqIdx]
		rest := spec[eqIdx+1:]
		sep := " "
		if pipeIdx := strings.LastIndex(rest, "|"); pipeIdx >= 0 {
			sep = rest[pipeIdx+1:]
			rest = rest[:pipeIdx]
		}
		sources := strings.Split(rest, ",")
		if len(sources) < 2 {
			return nil, fmt.Errorf("join: spec %q must have at least two source fields", spec)
		}
		for _, s := range sources {
			if s == "" {
				return nil, fmt.Errorf("join: spec %q contains empty source field", spec)
			}
		}
		rules = append(rules, Rule{Target: target, Sources: sources, Separator: sep})
	}
	return rules, nil
}
