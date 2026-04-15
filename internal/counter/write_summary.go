package counter

import (
	"fmt"
	"io"
)

// WriteSummary writes a human-readable frequency table for the counter's
// tracked field to w. Entries are listed in descending count order.
// If topN is greater than zero, only the top N entries are included.
func (c *Counter) WriteSummary(w io.Writer, topN int) error {
	entries := c.Top(topN)
	if len(entries) == 0 {
		_, err := fmt.Fprintf(w, "no data recorded for field %q\n", c.field)
		return err
	}
	_, err := fmt.Fprintf(w, "field: %s\n", c.field)
	if err != nil {
		return err
	}
	for _, vc := range entries {
		_, err = fmt.Fprintf(w, "  %-30s %d\n", vc.Value, vc.Count)
		if err != nil {
			return err
		}
	}
	return nil
}
