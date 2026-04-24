package lookup

// Package lookup provides field value enrichment by looking up values
// in a static key→value map and writing the result to a target field.

import (
	"errors"
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Rule describes a single lookup: read SourceField, find its string value in
// Table, and write the matched value to TargetField.
type Rule struct {
	SourceField string
	TargetField string
	Table       map[string]string
}

// Applier enriches log entries using a set of lookup rules.
type Applier struct {
	rules []Rule
}

// New returns an Applier that will apply the given rules in order.
func New(rules []Rule) *Applier {
	return &Applier{rules: rules}
}

// Apply returns a copy of entry with lookup fields populated.
// If a source field is missing or its value has no match in the table the
// target field is left untouched.
func (a *Applier) Apply(entry parser.Entry) parser.Entry {
	if len(a.rules) == 0 {
		return entry
	}
	out := entry.Clone()
	for _, r := range a.rules {
		v, ok := out.Fields[r.SourceField]
		if !ok {
			continue
		}
		key, ok := v.(string)
		if !ok {
			continue
		}
		if mapped, found := r.Table[key]; found {
			out.Fields[r.TargetField] = mapped
		}
	}
	return out
}

// ParseRules parses a slice of spec strings of the form:
//   source:target=k1->v1,k2->v2
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		r, err := parseRule(spec)
		if err != nil {
			return nil, fmt.Errorf("lookup: invalid spec %q: %w", spec, err)
		}
		rules = append(rules, r)
	}
	return rules, nil
}

func parseRule(spec string) (Rule, error) {
	eqIdx := strings.IndexByte(spec, '=')
	if eqIdx < 0 {
		return Rule{}, errors.New("missing '='")
	}
	header := spec[:eqIdx]
	pairs := spec[eqIdx+1:]

	parts := strings.SplitN(header, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Rule{}, errors.New("header must be 'source:target'")
	}

	table := map[string]string{}
	if pairs != "" {
		for _, pair := range strings.Split(pairs, ",") {
			arrow := strings.SplitN(pair, "->", 2)
			if len(arrow) != 2 || arrow[0] == "" {
				return Rule{}, fmt.Errorf("invalid pair %q, expected 'k->v'", pair)
			}
			table[arrow[0]] = arrow[1]
		}
	}
	return Rule{SourceField: parts[0], TargetField: parts[1], Table: table}, nil
}
