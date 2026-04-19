// Package lowercase provides an applier that converts string field values
// in a log entry to lowercase.
//
// Usage:
//
//	fields, err := lowercase.ParseFields("msg,level")
//	if err != nil {
//		log.Fatal(err)
//	}
//	a := lowercase.New(fields)
//	outEntry := a.Apply(entry)
//
// Only fields whose values are strings are affected; numeric or other typed
// values are passed through unchanged.
package lowercase
