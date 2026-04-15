// Package head provides a limiter that passes through only the first N log entries.
package head

import "github.com/logslice/logslice/internal/parser"

// Limiter keeps only the first N entries from a stream.
type Limiter struct {
	max   int
	seen  int
	done  bool
}

// Option is a functional option for Limiter.
type Option func(*Limiter)

// WithMax sets the maximum number of entries to pass through.
// A value of zero or less disables limiting (all entries pass).
func WithMax(n int) Option {
	return func(l *Limiter) {
		l.max = n
	}
}

// New creates a new Limiter with the given options.
func New(opts ...Option) *Limiter {
	l := &Limiter{}
	for _, o := range opts {
		o(l)
	}
	return l
}

// Apply returns the entry unchanged if the limit has not been reached,
// and signals completion once the limit is exceeded.
// It returns (entry, true) while under the limit, or (zero, false) once done.
func (l *Limiter) Apply(entry parser.Entry) (parser.Entry, bool) {
	if l.max <= 0 {
		return entry, true
	}
	if l.done {
		return parser.Entry{}, false
	}
	l.seen++
	if l.seen >= l.max {
		l.done = true
	}
	return entry, true
}

// Done reports whether the limiter has reached its maximum.
func (l *Limiter) Done() bool {
	return l.done
}

// Reset resets the limiter so it can be reused from the beginning.
func (l *Limiter) Reset() {
	l.seen = 0
	l.done = false
}
