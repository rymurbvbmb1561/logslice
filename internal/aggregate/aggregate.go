// Package aggregate provides field-based counting and grouping of log entries.
package aggregate

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/user/logslice/internal/parser"
)

// Counter counts occurrences of distinct values for a given field.
type Counter struct {
	field  string
	counts map[string]int
	total  int
}

// New returns a new Counter that groups entries by the given field name.
func New(field string) *Counter {
	return &Counter{
		field:  field,
		counts: make(map[string]int),
	}
}

// Record observes a log entry and increments the count for its field value.
func (c *Counter) Record(entry parser.Entry) {
	c.total++
	val, ok := entry.Fields[c.field]
	if !ok {
		c.counts["(missing)"]++
		return
	}
	c.counts[fmt.Sprintf("%v", val)]++
}

// Counts returns a copy of the current value→count map.
func (c *Counter) Counts() map[string]int {
	out := make(map[string]int, len(c.counts))
	for k, v := range c.counts {
		out[k] = v
	}
	return out
}

// Total returns the total number of entries recorded.
func (c *Counter) Total() int {
	return c.total
}

// WriteSummary writes a sorted tabular summary of counts to w.
func (c *Counter) WriteSummary(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "VALUE\tCOUNT\n")

	type kv struct {
		key   string
		count int
	}
	pairs := make([]kv, 0, len(c.counts))
	for k, v := range c.counts {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count != pairs[j].count {
			return pairs[i].count > pairs[j].count
		}
		return pairs[i].key < pairs[j].key
	})
	for _, p := range pairs {
		fmt.Fprintf(tw, "%s\t%d\n", p.key, p.count)
	}
	tw.Flush()
}
