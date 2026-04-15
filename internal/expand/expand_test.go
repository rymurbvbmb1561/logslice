package expand

import (
	"testing"
)

func makeEntry(pairs ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestApply_NoFields_ReturnsOriginal(t *testing.T) {
	e := New(nil)
	in := makeEntry("level", "info", "msg", "hello")
	out, err := e.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["level"] != "info" || out["msg"] != "hello" {
		t.Errorf("entry modified unexpectedly: %v", out)
	}
}

func TestApply_ExpandsJSONStringField(t *testing.T) {
	e := New([]string{"payload"})
	in := makeEntry("level", "info", "payload", `{"user":"alice","action":"login"}`)
	out, err := e.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["payload"]; ok {
		t.Error("expected 'payload' field to be removed")
	}
	if out["user"] != "alice" {
		t.Errorf("expected user=alice, got %v", out["user"])
	}
	if out["action"] != "login" {
		t.Errorf("expected action=login, got %v", out["action"])
	}
}

func TestApply_WithPrefix(t *testing.T) {
	e := New([]string{"meta"}, WithPrefix("meta_"))
	in := makeEntry("msg", "ok", "meta", `{"host":"srv1"}`)
	out, err := e.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["meta_host"] != "srv1" {
		t.Errorf("expected meta_host=srv1, got %v", out["meta_host"])
	}
}

func TestApply_NonJSONStringLeftAlone(t *testing.T) {
	e := New([]string{"payload"})
	in := makeEntry("payload", "plain text")
	out, err := e.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["payload"] != "plain text" {
		t.Errorf("expected payload unchanged, got %v", out["payload"])
	}
}

func TestApply_OverwriteFalse_PreservesExistingKeys(t *testing.T) {
	e := New([]string{"extra"}, WithOverwrite(false))
	in := makeEntry("user", "bob", "extra", `{"user":"alice"}`)
	out, err := e.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["user"] != "bob" {
		t.Errorf("expected user preserved as bob, got %v", out["user"])
	}
}

func TestApply_OverwriteTrue_ReplacesExistingKeys(t *testing.T) {
	e := New([]string{"extra"}, WithOverwrite(true))
	in := makeEntry("user", "bob", "extra", `{"user":"alice"}`)
	out, err := e.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["user"] != "alice" {
		t.Errorf("expected user overwritten to alice, got %v", out["user"])
	}
}

func TestApply_MissingFieldIgnored(t *testing.T) {
	e := New([]string{"nonexistent"})
	in := makeEntry("level", "warn")
	out, err := e.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["level"] != "warn" {
		t.Errorf("unexpected change to entry: %v", out)
	}
}
