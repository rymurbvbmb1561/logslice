// Package replace provides a Replacer that applies regex find-and-replace
// operations to string fields within log entries.
//
// Rules are specified as "field/pattern/replacement" where pattern is a
// regular expression. Non-string fields and missing fields are left unchanged.
//
// Example usage:
//
//	rules, err := replace.ParseRules([]string{"message/error/ERR", "host/\\.internal$""})
//	if err != nil {
//		log.Fatal(err)
//	}
//	r := replace.New(rules)
//	out := r.Apply(entry)
package replace
