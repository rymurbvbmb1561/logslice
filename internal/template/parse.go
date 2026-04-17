package template

import (
	"fmt"
	"strings"
)

// ParseTemplate validates and returns an Applier from a template string.
// Returns an error if the template has unclosed braces.
func ParseTemplate(tmpl string) (*Applier, error) {
	if err := validateBraces(tmpl); err != nil {
		return nil, err
	}
	return New(tmpl), nil
}

func validateBraces(tmpl string) error {
	depth := 0
	for i, ch := range tmpl {
		switch ch {
		case '{':
			depth++
			if depth > 1 {
				return fmt.Errorf("template: nested braces not supported at position %d", i)
			}
		case '}':
			if depth == 0 {
				return fmt.Errorf("template: unexpected '}' at position %d", i)
			}
			depth--
		}
	}
	if depth != 0 {
		return fmt.Errorf("template: unclosed '{' in template")
	}
	_ = strings.TrimSpace(tmpl) // no-op, just to use strings import
	return nil
}
