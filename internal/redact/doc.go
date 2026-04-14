// Package redact provides field-level and pattern-based redaction of log entries.
//
// A Redactor can be configured to:
//   - Replace specific named fields with a placeholder (e.g. "password", "api_key").
//   - Replace substrings matching regular expressions with a placeholder.
//   - Use built-in named patterns for common sensitive formats such as
//     email addresses, credit card numbers, Bearer tokens, and IPv4 addresses.
//
// Basic usage:
//
//	r := redact.New(
//		redact.WithFields("password", "token"),
//		redact.WithPatterns(redact.PatternEmail),
//		redact.WithPlaceholder("[REDACTED]"),
//	)
//	clean := r.Apply(entry)
//
// Named patterns can be resolved from CLI input using ParsePatterns:
//
//	pats, err := redact.ParsePatterns("email,creditcard")
//
// The original entry is never modified; Apply always returns a shallow copy.
package redact
