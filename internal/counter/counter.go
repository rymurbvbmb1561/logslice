package counter

import "github.com/logslice/logslice/internal/parser"

// Counter counts occurrences of values for a given field across log entries.
// It supports an optional limit to cap the number of tracked distinct values.
type Counter struct {
	field  string
	limit  int
	counts map[string]int
	order  []string
}

// Option is a functional option for configuring a Counter.
type Option func(*Counter)

// WithLimit sets the maximum number of distinct field values to track.
// Once the limit is reached, new distinct values are ignored.
func WithLimit(n int) Option {
	return func(c *Counter) {
		c.limit = n
	}
}

// New creates a Counter that tracks occurrences of values for the given field.
func New(field string, opts ...Option) *Counter {
	c := &Counter{
		field:  field,
		counts: make(map[string]int),
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// Record observes a log entry and increments the count for the field value.
func (c *Counter) Record(entry parser.Entry) {
	v, ok := entry[c.field]
	if !ok {
		return
	}
	key := fmt.Sprint(v)
	if _, seen := c.counts[key]; !seen {
		if c.limit > 0 && len(c.order) >= c.limit {
			return
		}
		c.order = append(c.order, key)
	}
	c.counts[key]++
}

// Counts returns a slice of ValueCount in insertion order.
func (c *Counter) Counts() []ValueCount {
	out := make([]ValueCount, 0, len(c.order))
	for _, k := range c.order {
		out = append(out, ValueCount{Value: k, Count: c.counts[k]})
	}
	return out
}

// Top returns up to n entries sorted by count descending.
func (c *Counter) Top(n int) []ValueCount {
	all := c.Counts()
	sort.Slice(all, func(i, j int) bool {
		return all[i].Count > all[j].Count
	})
	if n > 0 && n < len(all) {
		return all[:n]
	}
	return all
}

// ValueCount pairs a field value with its occurrence count.
type ValueCount struct {
	Value string
	Count int
}
