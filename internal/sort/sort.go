package sort

import (
	"fmt"
	"sort"

	"github.com/user/logslice/internal/parser"
)

// Order defines the sort direction.
type Order int

const (
	Ascending Order = iota
	Descending
)

// Sorter buffers log entries and emits them in sorted order by a field.
type Sorter struct {
	field   string
	order   Order
	entries []parser.Entry
}

// Option is a functional option for Sorter.
type Option func(*Sorter)

// WithOrder sets the sort order.
func WithOrder(o Order) Option {
	return func(s *Sorter) {
		s.order = o
	}
}

// New creates a Sorter that sorts by the given field.
func New(field string, opts ...Option) *Sorter {
	s := &Sorter{field: field, order: Ascending}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Add buffers an entry.
func (s *Sorter) Add(e parser.Entry) {
	s.entries = append(s.entries, e)
}

// Entries returns all buffered entries sorted by the configured field.
func (s *Sorter) Entries() []parser.Entry {
	out := make([]parser.Entry, len(s.entries))
	copy(out, s.entries)
	sort.SliceStable(out, func(i, j int) bool {
		vi := stringify(out[i].Fields[s.field])
		vj := stringify(out[j].Fields[s.field])
		if s.order == Descending {
			return vi > vj
		}
		return vi < vj
	})
	return out
}

// Reset clears the buffer.
func (s *Sorter) Reset() {
	s.entries = s.entries[:0]
}

func stringify(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}
