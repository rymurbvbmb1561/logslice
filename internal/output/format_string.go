package output

import "fmt"

// ParseFormat converts a string into a Format constant.
// Returns an error if the string does not match a known format.
func ParseFormat(s string) (Format, error) {
	switch Format(s) {
	case FormatJSON:
		return FormatJSON, nil
	case FormatText:
		return FormatText, nil
	case FormatRaw:
		return FormatRaw, nil
	default:
		return "", fmt.Errorf("unknown output format %q: must be one of json, text, raw", s)
	}
}

// String returns the string representation of a Format.
func (f Format) String() string {
	return string(f)
}

// KnownFormats returns all supported format identifiers as a slice of strings.
func KnownFormats() []string {
	return []string{
		string(FormatJSON),
		string(FormatText),
		string(FormatRaw),
	}
}
