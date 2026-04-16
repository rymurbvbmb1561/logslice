// Package fieldselect provides a Selector that keeps or drops specific fields
// from parsed log entries.
//
// In keep mode (default) only the listed fields are retained in the output
// entry. In drop mode (WithDrop) the listed fields are removed and all others
// are kept.
//
// Example – keep only "time", "level", and "msg":
//
//	s := fieldselect.New(fieldselect.WithFields([]string{"time", "level", "msg"}))
//	out := s.Apply(entry)
//
// Example – drop "password" and "token":
//
//	s := fieldselect.New(
//		fieldselect.WithFields([]string{"password", "token"}),
//		fieldselect.WithDrop(),
//	)
//	out := s.Apply(entry)
package fieldselect
