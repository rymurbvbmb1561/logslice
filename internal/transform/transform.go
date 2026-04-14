package transform

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Transformer applies a set of field transformations to a log entry.
type Transformer struct {
	rules []Rule
}

// Rule describes a single field transformation.
type Rule struct {
	Field  string
	Action Action
	Value  string
}

// Action represents the kind of transformation to apply.
type Action int

const (
	ActionRename Action = iota
	ActionSet
	ActionDelete
)

// New returns a Transformer configured with the given rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply applies all transformation rules to the entry, returning a modified copy.
func (t *Transformer) Apply(entry parser.Entry) parser.Entry {
	out := make(parser.Entry, len(entry))
	for k, v := range entry {
		out[k] = v
	}
	for _, r := range t.rules {
		switch r.Action {
		case ActionDelete:
			delete(out, r.Field)
		case ActionSet:
			out[r.Field] = r.Value
		case ActionRename:
			if val, ok := out[r.Field]; ok {
				out[r.Value] = val
				delete(out, r.Field)
			}
		}
	}
	return out
}

// ParseRules parses transformation specs of the form:
//
//	delete:field
//	set:field=value
//	rename:old=new
func ParseRules(specs []string) ([]Rule, error) {
	var rules []Rule
	for _, spec := range specs {
		parts := strings.SplitN(spec, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid transform spec %q: expected action:args", spec)
		}
		kind, args := parts[0], parts[1]
		switch kind {
		case "delete":
			if args == "" {
				return nil, fmt.Errorf("delete transform requires a field name")
			}
			rules = append(rules, Rule{Field: args, Action: ActionDelete})
		case "set":
			kv := strings.SplitN(args, "=", 2)
			if len(kv) != 2 || kv[0] == "" {
				return nil, fmt.Errorf("set transform requires field=value, got %q", args)
			}
			rules = append(rules, Rule{Field: kv[0], Action: ActionSet, Value: kv[1]})
		case "rename":
			kv := strings.SplitN(args, "=", 2)
			if len(kv) != 2 || kv[0] == "" || kv[1] == "" {
				return nil, fmt.Errorf("rename transform requires old=new, got %q", args)
			}
			rules = append(rules, Rule{Field: kv[0], Action: ActionRename, Value: kv[1]})
		default:
			return nil, fmt.Errorf("unknown transform action %q", kind)
		}
	}
	return rules, nil
}
