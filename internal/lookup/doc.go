// Package lookup enriches log entries by replacing or supplementing field
// values using a static lookup table.
//
// Each Rule reads a value from a source field, looks it up in a
// user-supplied key→value map, and writes the matched result to a target
// field. If the source field is absent or the value has no match the entry
// is left unchanged.
//
// Specs are parsed with ParseRules from strings of the form:
//
//	source:target=key1->value1,key2->value2
//
// Example:
//
//	rules, _ := lookup.ParseRules([]string{"status:label=200->OK,404->Not Found"})
//	applier := lookup.New(rules)
//	enriched := applier.Apply(entry)
package lookup
