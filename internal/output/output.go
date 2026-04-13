// Package output handles writing filtered log entries to various destinations.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/user/logslice/internal/parser"
)

// Format represents the output format for log entries.
type Format string

const (
	// FormatJSON outputs each entry as a JSON object.
	FormatJSON Format = "json"
	// FormatText outputs each entry as a plain text line.
	FormatText Format = "text"
	// FormatRaw outputs the original raw line as-is.
	FormatRaw Format = "raw"
)

// Writer writes log entries to an io.Writer in the specified format.
type Writer struct {
	w      io.Writer
	format Format
}

// NewWriter creates a new Writer targeting the given io.Writer and format.
func NewWriter(w io.Writer, format Format) *Writer {
	return &Writer{w: w, format: format}
}

// NewStdoutWriter creates a Writer that writes to stdout.
func NewStdoutWriter(format Format) *Writer {
	return NewWriter(os.Stdout, format)
}

// Write outputs a single log entry according to the configured format.
func (w *Writer) Write(entry parser.Entry) error {
	switch w.format {
	case FormatJSON:
		return w.writeJSON(entry)
	case FormatText:
		return w.writeText(entry)
	case FormatRaw:
		return w.writeRaw(entry)
	default:
		return fmt.Errorf("unknown format: %s", w.format)
	}
}

// WriteAll outputs all entries in the slice.
func (w *Writer) WriteAll(entries []parser.Entry) error {
	for _, e := range entries {
		if err := w.Write(e); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeJSON(entry parser.Entry) error {
	b, err := json.Marshal(entry.Fields)
	if err != nil {
		return fmt.Errorf("marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(w.w, "%s\n", b)
	return err
}

func (w *Writer) writeText(entry parser.Entry) error {
	_, err := fmt.Fprintf(w.w, "[%s] %v\n", entry.Timestamp.Format("2006-01-02T15:04:05Z07:00"), entry.Fields)
	return err
}

func (w *Writer) writeRaw(entry parser.Entry) error {
	_, err := fmt.Fprintf(w.w, "%s\n", entry.Raw)
	return err
}
