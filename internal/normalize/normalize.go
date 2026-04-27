// Package normalize provides field value normalization for log entries.
// It trims whitespace, collapses internal runs of whitespace, and optionally
// lowercases string fields before further processing.
package normalize

import (
	"strings"

	"logslice/internal/parser"
)

// Option is a functional option for the Normalizer.
type Option func(*Normalizer)

// Normalizer applies normalization rules to log entry fields.
type Normalizer struct {
	fields    []string
	lowercase bool
}

// WithFields restricts normalization to the specified fields.
// If empty, all string fields are normalized.
func WithFields(fields []string) Option {
	return func(n *Normalizer) {
		n.fields = fields
	}
}

// WithLowercase enables lowercasing of normalized string values.
func WithLowercase() Option {
	return func(n *Normalizer) {
		n.lowercase = true
	}
}

// New creates a Normalizer with the given options.
func New(opts ...Option) *Normalizer {
	n := &Normalizer{}
	for _, o := range opts {
		o(n)
	}
	return n
}

// Apply returns a new entry with normalized field values.
// The original entry is never mutated.
func (n *Normalizer) Apply(e parser.Entry) parser.Entry {
	out := make(parser.Entry, len(e))
	for k, v := range e {
		out[k] = v
	}

	applyTo := func(key string) {
		val, ok := out[key]
		if !ok {
			return
		}
		s, ok := val.(string)
		if !ok {
			return
		}
		s = strings.Join(strings.Fields(s), " ")
		if n.lowercase {
			s = strings.ToLower(s)
		}
		out[key] = s
	}

	if len(n.fields) == 0 {
		for k := range out {
			applyTo(k)
		}
	} else {
		for _, f := range n.fields {
			applyTo(f)
		}
	}
	return out
}

// ParseFields splits a comma-separated list of field names.
func ParseFields(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
