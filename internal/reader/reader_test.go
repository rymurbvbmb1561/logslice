package reader

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return path
}

func TestNewFileReader_ValidFile(t *testing.T) {
	path := writeTempFile(t, "line1\nline2\n")
	r, err := NewFileReader(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer r.Close()
	if r.Name() != path {
		t.Errorf("Name() = %q, want %q", r.Name(), path)
	}
}

func TestNewFileReader_MissingFile(t *testing.T) {
	_, err := NewFileReader("/nonexistent/path/file.log")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLines_ReturnsAllLines(t *testing.T) {
	content := "{\"level\":\"info\"}\n{\"level\":\"error\"}\n{\"level\":\"debug\"}\n"
	path := writeTempFile(t, content)
	r, err := NewFileReader(path)
	if err != nil {
		t.Fatalf("NewFileReader: %v", err)
	}
	defer r.Close()

	lines, err := r.Lines()
	if err != nil {
		t.Fatalf("Lines(): %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("got %d lines, want 3", len(lines))
	}
}

func TestLines_EmptyFile(t *testing.T) {
	path := writeTempFile(t, "")
	r, err := NewFileReader(path)
	if err != nil {
		t.Fatalf("NewFileReader: %v", err)
	}
	defer r.Close()

	lines, err := r.Lines()
	if err != nil {
		t.Fatalf("Lines(): %v", err)
	}
	if len(lines) != 0 {
		t.Errorf("got %d lines, want 0", len(lines))
	}
}

func TestStream_SendsAllLines(t *testing.T) {
	content := "alpha\nbeta\ngamma\n"
	path := writeTempFile(t, content)
	r, err := NewFileReader(path)
	if err != nil {
		t.Fatalf("NewFileReader: %v", err)
	}
	defer r.Close()

	ch := make(chan string, 10)
	if err := r.Stream(ch); err != nil {
		t.Fatalf("Stream(): %v", err)
	}

	var got []string
	for line := range ch {
		got = append(got, line)
	}
	if len(got) != 3 {
		t.Errorf("got %d lines from stream, want 3", len(got))
	}
}

func TestStdinReader_Name(t *testing.T) {
	r := NewStdinReader()
	if r.Name() != "<stdin>" {
		t.Errorf("Name() = %q, want \"<stdin>\"", r.Name())
	}
}
