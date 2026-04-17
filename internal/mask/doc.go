// Package mask provides a Masker that partially obscures string field values
// in log entries, useful for hiding sensitive data such as tokens, passwords,
// or credit card numbers while preserving partial visibility for debugging.
//
// Example usage:
//
//	m := mask.New(
//		mask.WithFields([]string{"token", "password"}),
//		mask.WithVisibleLeft(2),
//		mask.WithVisibleRight(2),
//		mask.WithChar('*'),
//	)
//	masked := m.Apply(entry)
//
// If visibleLeft + visibleRight >= len(value), the value is returned unchanged.
//
// The default mask character is '*', visibleLeft and visibleRight both default
// to 0, meaning the entire value is obscured unless overridden via options.
package mask
