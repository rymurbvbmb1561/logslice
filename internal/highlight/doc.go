// Package highlight provides ANSI terminal colour highlighting for log output.
//
// A Highlighter is created with a set of Rules, each mapping a field name to a
// terminal colour. When enabled, Apply rewrites a raw log line so that every
// occurrence of the named field is wrapped in the corresponding ANSI escape
// sequence, making it easier to spot important fields at a glance.
//
// Usage:
//
//	rules, err := highlight.ParseRules([]string{"level=red", "msg=cyan"})
//	if err != nil { ... }
//	h := highlight.New(true, rules)
//	fmt.Println(h.Apply(line))
//
// Highlighting can be disabled (e.g. when stdout is not a TTY) by passing
// enabled=false to New, in which case Apply returns the line unchanged.
package highlight
