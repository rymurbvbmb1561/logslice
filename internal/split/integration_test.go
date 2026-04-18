package split_test

import (
	"testing"

	"github.com/logslice/logslice/internal/split"
)

func TestParseAndApply_RoundTrip(t *testing.T) {
	specs := []string{"endpoint:method,path,proto| "}
	rules, err := split.ParseRules(specs)
	if err != nil {
		t.Fatalf("ParseRules: %v", err)
	}
	s := split.New(rules)
	e := makeEntry(map[string]any{"endpoint": "GET /api/v1/logs HTTP/1.1"})
	out := s.Apply(e)

	if out.Fields["method"] != "GET" {
		t.Errorf("method = %v", out.Fields["method"])
	}
	if out.Fields["path"] != "/api/v1/logs" {
		t.Errorf("path = %v", out.Fields["path"])
	}
	if out.Fields["proto"] != "HTTP/1.1" {
		t.Errorf("proto = %v", out.Fields["proto"])
	}
	if out.Fields["endpoint"] != "GET /api/v1/logs HTTP/1.1" {
		t.Errorf("source field should be preserved")
	}
}

func TestParseAndApply_EmptyRules_EntryUnchanged(t *testing.T) {
	rules, err := split.ParseRules(nil)
	if err != nil {
		t.Fatalf("ParseRules: %v", err)
	}
	s := split.New(rules)
	e := makeEntry(map[string]any{"msg": "hello"})
	out := s.Apply(e)
	if out.Fields["msg"] != "hello" {
		t.Errorf("entry should be unchanged")
	}
}
