package mask_test

import (
	"testing"

	"github.com/user/logslice/internal/mask"
	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	e := make(parser.Entry)
	for k, v := range fields {
		e[k] = v
	}
	return e
}

func TestApply_NoFields_ReturnsOriginal(t *testing.T) {
	m := mask.New()
	entry := makeEntry(map[string]interface{}{"token": "abc123"})
	out := m.Apply(entry)
	if out["token"] != "abc123" {
		t.Errorf("expected unchanged, got %v", out["token"])
	}
}

func TestApply_MasksEntireField(t *testing.T) {
	m := mask.New(mask.WithFields([]string{"token"}))
	entry := makeEntry(map[string]interface{}{"token": "abc123"})
	out := m.Apply(entry)
	if out["token"] != "******" {
		t.Errorf("expected '******', got %v", out["token"])
	}
}

func TestApply_VisibleLeft(t *testing.T) {
	m := mask.New(mask.WithFields([]string{"token"}), mask.WithVisibleLeft(2))
	entry := makeEntry(map[string]interface{}{"token": "abc123"})
	out := m.Apply(entry)
	if out["token"] != "ab****" {
		t.Errorf("expected 'ab****', got %v", out["token"])
	}
}

func TestApply_VisibleRight(t *testing.T) {
	m := mask.New(mask.WithFields([]string{"token"}), mask.WithVisibleRight(2))
	entry := makeEntry(map[string]interface{}{"token": "abc123"})
	out := m.Apply(entry)
	if out["token"] != "****23" {
		t.Errorf("expected '****23', got %v", out["token"])
	}
}

func TestApply_VisibleLeftAndRight(t *testing.T) {
	m := mask.New(mask.WithFields([]string{"card"}), mask.WithVisibleLeft(1), mask.WithVisibleRight(2))
	entry := makeEntry(map[string]interface{}{"card": "1234567890"})
	out := m.Apply(entry)
	if out["card"] != "1*******90" {
		t.Errorf("expected '1*******90', got %v", out["card"])
	}
}

func TestApply_CustomChar(t *testing.T) {
	m := mask.New(mask.WithFields([]string{"pass"}), mask.WithChar('#'))
	entry := makeEntry(map[string]interface{}{"pass": "secret"})
	out := m.Apply(entry)
	if out["pass"] != "######" {
		t.Errorf("expected '######', got %v", out["pass"])
	}
}

func TestApply_NonStringFieldUnchanged(t *testing.T) {
	m := mask.New(mask.WithFields([]string{"count"}))
	entry := makeEntry(map[string]interface{}{"count": 42})
	out := m.Apply(entry)
	if out["count"] != 42 {
		t.Errorf("expected 42, got %v", out["count"])
	}
}

func TestApply_ShortStringNotMaskedWhenVisibleExceedsLength(t *testing.T) {
	m := mask.New(mask.WithFields([]string{"tok"}), mask.WithVisibleLeft(3), mask.WithVisibleRight(3))
	entry := makeEntry(map[string]interface{}{"tok": "ab"})
	out := m.Apply(entry)
	// visible left+right >= len, so original is returned
	if out["tok"] != "ab" {
		t.Errorf("expected 'ab' unchanged, got %v", out["tok"])
	}
}
