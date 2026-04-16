// Package rename provides field renaming for structured log entries.
//
// Rules are expressed as "old=new" specs and applied in order. If the source
// field does not exist in an entry the rule is silently skipped. The original
// entry is never mutated; Apply always returns a shallow copy.
//
// Example:
//
//	rules, err := rename.ParseRules([]string{"msg=message", "ts=timestamp"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	r := rename.New(rules)
//	outEntry := r.Apply(inEntry)
package rename
