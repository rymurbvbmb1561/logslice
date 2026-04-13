package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
)

func makeEntry(ts time.Time, raw string, fields map[string]interface{}) parser.Entry {
	return parser.Entry{
		Timestamp: ts,
		Raw:       raw,
		Fields:    fields,
	}
}

func TestWriter_FormatJSON(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatJSON)

	ts := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	entry := makeEntry(ts, `{"level":"info","msg":"hello"}`, map[string]interface{}{
		"level": "info",
		"msg":   "hello",
	})

	if err := w.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"level\":\"info\"") {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestWriter_FormatRaw(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatRaw)

	rawLine := `{"level":"warn","msg":"test"}`
	ts := time.Now()
	entry := makeEntry(ts, rawLine, nil)

	if err := w.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), rawLine) {
		t.Errorf("expected raw line in output, got: %s", buf.String())
	}
}

func TestWriter_FormatText(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatText)

	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	entry := makeEntry(ts, "", map[string]interface{}{"msg": "hello"})

	if err := w.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "2024-06-01T12:00:00Z") {
		t.Errorf("expected timestamp in text output, got: %s", buf.String())
	}
}

func TestWriter_WriteAll(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.FormatRaw)

	entries := []parser.Entry{
		makeEntry(time.Now(), "line one", nil),
		makeEntry(time.Now(), "line two", nil),
	}

	if err := w.WriteAll(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "line one") || !strings.Contains(out, "line two") {
		t.Errorf("expected both lines in output, got: %s", out)
	}
}

func TestWriter_UnknownFormat(t *testing.T) {
	var buf bytes.Buffer
	w := output.NewWriter(&buf, output.Format("xml"))
	err := w.Write(makeEntry(time.Now(), "", nil))
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}
