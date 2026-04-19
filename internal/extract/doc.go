// Package extract provides regex-based field extraction for log entries.
//
// Rules are specified as "field=pattern" where pattern is a regular expression
// containing one or more named capture groups ((?P<name>...)). When the pattern
// matches the value of the source field, each named group is written as a new
// field on the entry.
//
// Example:
//
//	rules, _ := extract.ParseRules([]string{
//	    `message=(?P<level>\w+)\s+(?P<detail>.+)`,
//	})
//	e := extract.New(rules)
//	out := e.Apply(entry)
package extract
