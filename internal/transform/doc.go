// Package transform provides field-level transformations for log entries.
//
// A [Transformer] applies an ordered list of [Rule] values to each
// [parser.Entry], producing a modified copy without mutating the original.
//
// Three actions are supported:
//
//   - delete  — removes a field from the entry
//   - set     — adds or overwrites a field with a fixed value
//   - rename  — renames an existing field, leaving others untouched
//
// Rules are parsed from string specs via [ParseRules], which accepts the
// following formats:
//
//	delete:<field>
//	set:<field>=<value>
//	rename:<old>=<new>
//
// Example:
//
//	rules, err := transform.ParseRules([]string{"delete:secret", "rename:message=msg"})
//	tr := transform.New(rules)
//	modified := tr.Apply(entry)
package transform
