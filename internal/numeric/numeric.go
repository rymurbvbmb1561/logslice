// Package numeric provides a processor that filters log entries based on
// numeric field comparisons (greater than, less than, equal to, etc.).
package numeric

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/your-org/logslice/internal/parser"
)

// Op represents a numeric comparison operator.
type Op int

const (
	OpGT  Op = iota // >
	OpGTE            // >=
	OpLT             // <
	OpLTE            // <=
	OpEQ             // ==
	OpNEQ            // !=
)

// Rule defines a single numeric filter: field op threshold.
type Rule struct {
	Field     string
	Op        Op
	Threshold float64
}

// Numeric filters entries by evaluating numeric rules against field values.
type Numeric struct {
	rules []Rule
}

// New returns a Numeric processor with the given rules.
func New(rules []Rule) *Numeric {
	return &Numeric{rules: rules}
}

// Apply returns true if the entry satisfies all numeric rules.
func (n *Numeric) Apply(entry parser.Entry) bool {
	for _, r := range n.rules {
		v, ok := entry.Fields[r.Field]
		if !ok {
			return false
		}
		f, err := toFloat(v)
		if err != nil {
			return false
		}
		if !compare(f, r.Op, r.Threshold) {
			return false
		}
	}
	return true
}

func compare(val float64, op Op, threshold float64) bool {
	switch op {
	case OpGT:
		return val > threshold
	case OpGTE:
		return val >= threshold
	case OpLT:
		return val < threshold
	case OpLTE:
		return val <= threshold
	case OpEQ:
		return val == threshold
	case OpNEQ:
		return val != threshold
	}
	return false
}

func toFloat(v interface{}) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case int:
		return float64(x), nil
	case string:
		return strconv.ParseFloat(strings.TrimSpace(x), 64)
	}
	return 0, fmt.Errorf("unsupported type %T", v)
}
