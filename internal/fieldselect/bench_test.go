package fieldselect_test

import (
	"testing"

	"github.com/user/logslice/internal/fieldselect"
	"github.com/user/logslice/internal/parser"
)

func BenchmarkApply_Keep(b *testing.B) {
	s := fieldselect.New(fieldselect.WithFields([]string{"time", "level", "msg"}))
	e := parser.Entry{
		Raw: `{}`,
		Fields: map[string]any{
			"time":   "2024-01-01T00:00:00Z",
			"level":  "info",
			"msg":    "benchmark entry",
			"caller": "main.go:42",
			"trace":  "abc123",
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Apply(e)
	}
}

func BenchmarkApply_Drop(b *testing.B) {
	s := fieldselect.New(fieldselect.WithFields([]string{"trace", "caller"}), fieldselect.WithDrop())
	e := parser.Entry{
		Raw: `{}`,
		Fields: map[string]any{
			"time":   "2024-01-01T00:00:00Z",
			"level":  "info",
			"msg":    "benchmark entry",
			"caller": "main.go:42",
			"trace":  "abc123",
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Apply(e)
	}
}
