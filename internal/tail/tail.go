// Package tail provides a processor that keeps only the last N log entries.
package tail

import "github.com/logslice/logslice/internal/parser"

// Option configures a Tail processor.
type Option func(*Tail)

// WithMax sets the maximum number of entries to retain.
func WithMax(n int) Option {
	return func(t *Tail) {
		if n > 0 {
			t.max = n
		}
	}
}

// Tail keeps a rolling buffer of the last N entries.
type Tail struct {
	max    int
	buffer []parser.Entry
}

// New creates a Tail processor. With no options all entries are kept.
func New(opts ...Option) *Tail {
	t := &Tail{}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Record adds an entry to the internal buffer, evicting the oldest when full.
func (t *Tail) Record(e parser.Entry) {
	if t.max <= 0 {
		t.buffer = append(t.buffer, e)
		return
	}
	if len(t.buffer) >= t.max {
		t.buffer = t.buffer[1:]
	}
	t.buffer = append(t.buffer, e)
}

// Entries returns the retained entries in order.
func (t *Tail) Entries() []parser.Entry {
	out := make([]parser.Entry, len(t.buffer))
	copy(out, t.buffer)
	return out
}

// Len returns the current number of buffered entries.
func (t *Tail) Len() int {
	return len(t.buffer)
}
