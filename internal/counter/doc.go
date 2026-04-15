// Package counter provides a field-value frequency counter for structured log
// entries. It tracks how many times each distinct value appears for a named
// field, and can return results sorted by frequency via Top.
//
// Example usage:
//
//	c := counter.New("level", counter.WithLimit(100))
//	for _, entry := range entries {
//		c.Record(entry)
//	}
//	for _, vc := range c.Top(10) {
//		fmt.Printf("%s: %d\n", vc.Value, vc.Count)
//	}
package counter
