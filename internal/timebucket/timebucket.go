// Package timebucket groups log entries into fixed-size time buckets
// and counts occurrences per bucket.
package timebucket

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/user/logslice/internal/parser"
)

// Bucketer groups entries by truncated time intervals.
type Bucketer struct {
	interval time.Duration
	buckets  map[time.Time]int
}

// WithInterval sets the bucket interval. Defaults to 1 minute if zero.
func WithInterval(d time.Duration) func(*Bucketer) {
	return func(b *Bucketer) {
		if d > 0 {
			b.interval = d
		}
	}
}

// New creates a Bucketer with optional configuration.
func New(opts ...func(*Bucketer)) *Bucketer {
	b := &Bucketer{
		interval: time.Minute,
		buckets:  make(map[time.Time]int),
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

// Record adds an entry to the appropriate time bucket.
// Entries without a timestamp are ignored.
func (b *Bucketer) Record(e parser.Entry) {
	if e.Timestamp.IsZero() {
		return
	}
	key := e.Timestamp.Truncate(b.interval)
	b.buckets[key]++
}

// WriteSummary writes a sorted table of bucket counts to w.
func (b *Bucketer) WriteSummary(w io.Writer) {
	type row struct {
		t time.Time
		n int
	}
	rows := make([]row, 0, len(b.buckets))
	for t, n := range b.buckets {
		rows = append(rows, row{t, n})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].t.Before(rows[j].t)
	})
	fmt.Fprintln(w, "bucket\t\t\t\tcount")
	for _, r := range rows {
		fmt.Fprintf(w, "%s\t%d\n", r.t.UTC().Format(time.RFC3339), r.n)
	}
}

// Buckets returns a copy of the internal bucket map.
func (b *Bucketer) Buckets() map[time.Time]int {
	out := make(map[time.Time]int, len(b.buckets))
	for k, v := range b.buckets {
		out[k] = v
	}
	return out
}
