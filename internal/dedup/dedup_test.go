package dedup

import (
	"testing"
)

func makeEntry(fields map[string]interface{}) Entry {
	e := make(Entry)
	for k, v := range fields {
		e[k] = v
	}
	return e
}

func TestIsDuplicate_FullEntry_FirstSeenNotDuplicate(t *testing.T) {
	d := New("")
	e := makeEntry(map[string]interface{}{"msg": "hello", "level": "info"})
	dup, err := d.IsDuplicate(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dup {
		t.Error("expected not duplicate on first occurrence")
	}
}

func TestIsDuplicate_FullEntry_SecondSeenIsDuplicate(t *testing.T) {
	d := New("")
	e := makeEntry(map[string]interface{}{"msg": "hello", "level": "info"})
	_, _ = d.IsDuplicate(e)
	dup, err := d.IsDuplicate(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dup {
		t.Error("expected duplicate on second occurrence")
	}
}

func TestIsDuplicate_ByField_DeduplicatesOnFieldValue(t *testing.T) {
	d := New("request_id")
	e1 := makeEntry(map[string]interface{}{"request_id": "abc-123", "msg": "first"})
	e2 := makeEntry(map[string]interface{}{"request_id": "abc-123", "msg": "second"})

	_, _ = d.IsDuplicate(e1)
	dup, err := d.IsDuplicate(e2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !dup {
		t.Error("expected duplicate when field value matches")
	}
}

func TestIsDuplicate_ByField_DifferentValueNotDuplicate(t *testing.T) {
	d := New("request_id")
	e1 := makeEntry(map[string]interface{}{"request_id": "abc-123"})
	e2 := makeEntry(map[string]interface{}{"request_id": "xyz-999"})

	_, _ = d.IsDuplicate(e1)
	dup, err := d.IsDuplicate(e2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dup {
		t.Error("expected not duplicate when field value differs")
	}
}

func TestIsDuplicate_ByField_MissingFieldReturnsError(t *testing.T) {
	d := New("request_id")
	e := makeEntry(map[string]interface{}{"msg": "no id here"})
	_, err := d.IsDuplicate(e)
	if err == nil {
		t.Error("expected error when field is missing")
	}
}

func TestSeen_CountsUniqueEntries(t *testing.T) {
	d := New("")
	for i := 0; i < 5; i++ {
		e := makeEntry(map[string]interface{}{"n": i})
		_, _ = d.IsDuplicate(e)
	}
	if d.Seen() != 5 {
		t.Errorf("expected 5 unique entries, got %d", d.Seen())
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := New("")
	e := makeEntry(map[string]interface{}{"msg": "hello"})
	_, _ = d.IsDuplicate(e)
	d.Reset()
	if d.Seen() != 0 {
		t.Errorf("expected 0 after reset, got %d", d.Seen())
	}
	dup, _ := d.IsDuplicate(e)
	if dup {
		t.Error("expected not duplicate after reset")
	}
}
