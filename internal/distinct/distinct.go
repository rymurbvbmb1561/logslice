// Package distinct provides a processor that filters log entries,
// keeping only those with unique values for a given field.
package distinct

import "github.com/logslice/logslice/internal/parser"

// Option configures a Processor.
type Option func(*Processor)

// WithField sets the field used for distinctness comparison.
func WithField(field string) Option {
	return func(p *Processor) {
		p.field = field
	}
}

// Processor filters entries to unique values of a field.
type Processor struct {
	field string
	seen  map[string]struct{}
}

// New creates a new distinct Processor.
// If no field is set via options, every entry is treated as distinct (all pass).
func New(opts ...Option) *Processor {
	p := &Processor{
		seen: make(map[string]struct{}),
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Apply returns true if the entry should be kept (i.e. the field value is new).
func (p *Processor) Apply(entry parser.Entry) bool {
	if p.field == "" {
		return true
	}
	v, ok := entry.Fields[p.field]
	if !ok {
		return true
	}
	key, ok := v.(string)
	if !ok {
		// For non-string values use fmt representation via type assertion chain.
		key = stringify(v)
	}
	if _, seen := p.seen[key]; seen {
		return false
	}
	p.seen[key] = struct{}{}
	return true
}

// Reset clears the seen set, allowing the processor to be reused.
func (p *Processor) Reset() {
	p.seen = make(map[string]struct{})
}

func stringify(v any) string {
	switch val := v.(type) {
	case float64:
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", val)
	}
}
