package redact_test

import (
	"regexp"
	"testing"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/redact"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	e := make(parser.Entry)
	for k, v := range fields {
		e[k] = v
	}
	return e
}

func TestApply_RedactsNamedField(t *testing.T) {
	r := redact.New(redact.WithFields("password"))
	entry := makeEntry(map[string]interface{}{"user": "alice", "password": "secret"})
	out := r.Apply(entry)
	if out["password"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", out["password"])
	}
	if out["user"] != "alice" {
		t.Errorf("user field should be unchanged")
	}
}

func TestApply_MissingFieldIgnored(t *testing.T) {
	r := redact.New(redact.WithFields("token"))
	entry := makeEntry(map[string]interface{}{"msg": "hello"})
	out := r.Apply(entry)
	if _, ok := out["token"]; ok {
		t.Error("token field should not be added")
	}
}

func TestApply_PatternRedactsMatchingValue(t *testing.T) {
	pat := regexp.MustCompile(`\b\d{4}-\d{4}-\d{4}-\d{4}\b`)
	r := redact.New(redact.WithPatterns(pat))
	entry := makeEntry(map[string]interface{}{"msg": "card 1234-5678-9012-3456 used"})
	out := r.Apply(entry)
	expected := "card [REDACTED] used"
	if out["msg"] != expected {
		t.Errorf("expected %q, got %q", expected, out["msg"])
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	r := redact.New(redact.WithFields("secret"), redact.WithPlaceholder("***"))
	entry := makeEntry(map[string]interface{}{"secret": "topsecret"})
	out := r.Apply(entry)
	if out["secret"] != "***" {
		t.Errorf("expected ***, got %v", out["secret"])
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	pat := regexp.MustCompile(`\d+`)
	r := redact.New(redact.WithPatterns(pat))
	entry := makeEntry(map[string]interface{}{"count": 42})
	out := r.Apply(entry)
	if out["count"] != 42 {
		t.Errorf("numeric field should not be altered")
	}
}

func TestApply_OriginalEntryUnmodified(t *testing.T) {
	r := redact.New(redact.WithFields("api_key"))
	entry := makeEntry(map[string]interface{}{"api_key": "abc123"})
	r.Apply(entry)
	if entry["api_key"] != "abc123" {
		t.Error("original entry should not be modified")
	}
}

func TestParseFields_CommaSeparated(t *testing.T) {
	fields := redact.ParseFields("user, password , token")
	expected := []string{"user", "password", "token"}
	if len(fields) != len(expected) {
		t.Fatalf("expected %d fields, got %d", len(expected), len(fields))
	}
	for i, f := range fields {
		if f != expected[i] {
			t.Errorf("field[%d]: expected %q, got %q", i, expected[i], f)
		}
	}
}

func TestParseFields_Empty_ReturnsNil(t *testing.T) {
	fields := redact.ParseFields("")
	if fields != nil {
		t.Errorf("expected nil, got %v", fields)
	}
}
