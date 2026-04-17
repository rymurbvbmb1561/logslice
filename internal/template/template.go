package template

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Applier renders log entries using a Go-style template string.
// Fields are referenced as {fieldname}.
type Applier struct {
	template string
	fields   []string
}

// New creates an Applier with the given template string.
func New(tmpl string) *Applier {
	return &Applier{template: tmpl, fields: extractFields(tmpl)}
}

// Apply renders the entry using the template, returning the formatted string.
// Missing fields are replaced with "<nil>".
func (a *Applier) Apply(entry parser.Entry) (string, error) {
	if a.template == "" {
		return "", fmt.Errorf("template: empty template")
	}
	var buf bytes.Buffer
	tmpl := a.template
	for _, field := range a.fields {
		placeholder := "{" + field + "}"
		var val string
		if v, ok := entry.Fields[field]; ok {
			val = fmt.Sprintf("%v", v)
		} else {
			val = "<nil>"
		}
		tmpl = strings.ReplaceAll(tmpl, placeholder, val)
	}
	buf.WriteString(tmpl)
	return buf.String(), nil
}

// extractFields parses {field} placeholders from the template.
func extractFields(tmpl string) []string {
	var fields []string
	seen := map[string]bool{}
	rest := tmpl
	for {
		start := strings.Index(rest, "{")
		if start == -1 {
			break
		}
		end := strings.Index(rest[start:], "}")
		if end == -1 {
			break
		}
		field := rest[start+1 : start+end]
		if field != "" && !seen[field] {
			fields = append(fields, field)
			seen[field] = true
		}
		rest = rest[start+end+1:]
	}
	return fields
}
