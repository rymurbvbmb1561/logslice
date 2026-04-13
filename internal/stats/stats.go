// Package stats provides utilities for collecting and reporting
// processing statistics during log slicing operations.
package stats

import (
	"fmt"
	"io"
	"time"
)

// Collector tracks counts and timing for a log processing run.
type Collector struct {
	StartTime   time.Time
	LinesRead   int
	LinesParsed int
	LinesMatched int
	ParseErrors int
}

// New returns a new Collector with the start time set to now.
func New() *Collector {
	return &Collector{
		StartTime: time.Now(),
	}
}

// RecordRead increments the lines-read counter.
func (c *Collector) RecordRead() {
	c.LinesRead++
}

// RecordParsed increments the successfully-parsed counter.
func (c *Collector) RecordParsed() {
	c.LinesParsed++
}

// RecordMatched increments the matched (output) counter.
func (c *Collector) RecordMatched() {
	c.LinesMatched++
}

// RecordParseError increments the parse-error counter.
func (c *Collector) RecordParseError() {
	c.ParseErrors++
}

// Elapsed returns the duration since the collector was created.
func (c *Collector) Elapsed() time.Duration {
	return time.Since(c.StartTime)
}

// Print writes a human-readable summary to w.
func (c *Collector) Print(w io.Writer) {
	fmt.Fprintf(w, "lines read:    %d\n", c.LinesRead)
	fmt.Fprintf(w, "lines parsed:  %d\n", c.LinesParsed)
	fmt.Fprintf(w, "lines matched: %d\n", c.LinesMatched)
	if c.ParseErrors > 0 {
		fmt.Fprintf(w, "parse errors:  %d\n", c.ParseErrors)
	}
	fmt.Fprintf(w, "elapsed:       %s\n", c.Elapsed().Round(time.Millisecond))
}
