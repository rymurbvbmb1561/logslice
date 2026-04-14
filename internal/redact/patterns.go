package redact

import "regexp"

// Common pre-compiled patterns for well-known sensitive data formats.
var (
	// PatternEmail matches a simple email address.
	PatternEmail = regexp.MustCompile(`[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}`)

	// PatternCreditCard matches 16-digit card numbers in groups of 4.
	PatternCreditCard = regexp.MustCompile(`\b\d{4}[\-\s]\d{4}[\-\s]\d{4}[\-\s]\d{4}\b`)

	// PatternBearerToken matches Bearer authorization header values.
	PatternBearerToken = regexp.MustCompile(`(?i)bearer\s+[A-Za-z0-9\-._~+/]+=*`)

	// PatternIPv4 matches IPv4 addresses.
	PatternIPv4 = regexp.MustCompile(`\b(?:\d{1,3}\.){3}\d{1,3}\b`)
)

// knownPatternNames maps CLI-friendly names to compiled patterns.
var knownPatternNames = map[string]*regexp.Regexp{
	"email":       PatternEmail,
	"creditcard":  PatternCreditCard,
	"bearertoken": PatternBearerToken,
	"ipv4":        PatternIPv4,
}

// ParsePatterns parses a comma-separated list of named patterns.
// Returns an error for unknown pattern names.
func ParsePatterns(s string) ([]*regexp.Regexp, error) {
	if s == "" {
		return nil, nil
	}
	parts := splitTrim(s)
	out := make([]*regexp.Regexp, 0, len(parts))
	for _, name := range parts {
		pat, ok := knownPatternNames[name]
		if !ok {
			return nil, &unknownPatternError{name: name}
		}
		out = append(out, pat)
	}
	return out, nil
}

// KnownPatternNames returns the list of built-in pattern names.
func KnownPatternNames() []string {
	names := make([]string, 0, len(knownPatternNames))
	for k := range knownPatternNames {
		names = append(names, k)
	}
	return names
}

type unknownPatternError struct{ name string }

func (e *unknownPatternError) Error() string {
	return "redact: unknown pattern name: " + e.name
}

func splitTrim(s string) []string {
	var out []string
	for _, p := range splitComma(s) {
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func splitComma(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == ',' {
			parts = append(parts, trimSpace(s[start:i]))
			start = i + 1
		}
	}
	parts = append(parts, trimSpace(s[start:]))
	return parts
}

func trimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t') {
		s = s[1:]
	}
	for len(s) > 0 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t') {
		s = s[:len(s)-1]
	}
	return s
}
