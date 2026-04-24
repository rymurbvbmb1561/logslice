package numeric

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseRules parses a slice of rule specs into []Rule.
//
// Each spec has the form: field op threshold
// For example: "duration>=200" or "status==404" or "retries!=0".
//
// Supported operators: >, >=, <, <=, ==, !=
func ParseRules(specs []string) ([]Rule, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	ops := []struct {
		sym string
		op  Op
	}{
		{">", OpGT},
		{">=", OpGTE},
		{"<", OpLT},
		{"<=", OpLTE},
		{"==", OpEQ},
		{"!=", OpNEQ},
	}

	var rules []Rule
	for _, spec := range specs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}
		matched := false
		// Try two-char operators first to avoid partial matches.
		for _, candidate := range ops {
			idx := strings.Index(spec, candidate.sym)
			if idx < 0 {
				continue
			}
			// Skip if a longer operator was already matched at the same position.
			field := strings.TrimSpace(spec[:idx])
			raw := strings.TrimSpace(spec[idx+len(candidate.sym):])
			if field == "" {
				return nil, fmt.Errorf("numeric: empty field in spec %q", spec)
			}
			if raw == "" {
				return nil, fmt.Errorf("numeric: missing threshold in spec %q", spec)
			}
			threshold, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return nil, fmt.Errorf("numeric: invalid threshold %q in spec %q", raw, spec)
			}
			rules = append(rules, Rule{Field: field, Op: candidate.op, Threshold: threshold})
			matched = true
			break
		}
		if !matched {
			return nil, fmt.Errorf("numeric: no valid operator found in spec %q", spec)
		}
	}
	return rules, nil
}
