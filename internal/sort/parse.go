package sort

import (
	"fmt"
	"strings"
)

// Spec holds parsed sort configuration.
type Spec struct {
	Field string
	Order Order
}

// ParseSpec parses a sort spec of the form "field" or "field:asc" / "field:desc".
func ParseSpec(s string) (Spec, error) {
	if s == "" {
		return Spec{}, fmt.Errorf("sort spec must not be empty")
	}
	parts := strings.SplitN(s, ":", 2)
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return Spec{}, fmt.Errorf("sort field must not be empty")
	}
	order := Ascending
	if len(parts) == 2 {
		switch strings.ToLower(strings.TrimSpace(parts[1])) {
		case "asc", "":
			order = Ascending
		case "desc":
			order = Descending
		default:
			return Spec{}, fmt.Errorf("unknown sort order %q: want asc or desc", parts[1])
		}
	}
	return Spec{Field: field, Order: order}, nil
}
