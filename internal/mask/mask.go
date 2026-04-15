package mask

import (
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Masker partially obscures field values, revealing only a prefix or suffix.
type Masker struct {
	fields      map[string]struct{}
	visibleLeft int
	visibleRight int
	char        rune
}

// Option configures a Masker.
type Option func(*Masker)

// WithFields sets the field names whose values will be masked.
func WithFields(fields []string) Option {
	return func(m *Masker) {
		for _, f := range fields {
			if f != "" {
				m.fields[f] = struct{}{}
			}
		}
	}
}

// WithVisibleLeft keeps n characters visible at the start of the value.
func WithVisibleLeft(n int) Option {
	return func(m *Masker) { m.visibleLeft = n }
}

// WithVisibleRight keeps n characters visible at the end of the value.
func WithVisibleRight(n int) Option {
	return func(m *Masker) { m.visibleRight = n }
}

// WithChar sets the masking character (default '*').
func WithChar(c rune) Option {
	return func(m *Masker) { m.char = c }
}

// New creates a Masker with the given options.
func New(opts ...Option) *Masker {
	m := &Masker{
		fields: make(map[string]struct{}),
		char:   '*',
	}
	for _, o := range opts {
		o(m)
	}
	return m
}

// Apply returns a copy of entry with configured fields masked.
func (m *Masker) Apply(entry parser.Entry) parser.Entry {
	if len(m.fields) == 0 {
		return entry
	}
	out := make(parser.Entry, len(entry))
	for k, v := range entry {
		if _, ok := m.fields[k]; ok {
			if s, isStr := v.(string); isStr {
				out[k] = m.maskString(s)
				continue
			}
		}
		out[k] = v
	}
	return out
}

func (m *Masker) maskString(s string) string {
	runes := []rune(s)
	n := len(runes)
	left := m.visibleLeft
	right := m.visibleRight
	if left+right >= n {
		return s
	}
	masked := n - left - right
	var b strings.Builder
	b.WriteString(string(runes[:left]))
	for i := 0; i < masked; i++ {
		b.WriteRune(m.char)
	}
	if right > 0 {
		b.WriteString(string(runes[n-right:]))
	}
	return b.String()
}
