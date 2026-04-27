package sort

import (
	"testing"
)

func TestParseSpec_FieldOnly(t *testing.T) {
	spec, err := ParseSpec("level")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.Field != "level" || spec.Order != Ascending {
		t.Errorf("unexpected spec: %+v", spec)
	}
}

func TestParseSpec_ExplicitAsc(t *testing.T) {
	spec, err := ParseSpec("ts:asc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.Order != Ascending {
		t.Errorf("expected ascending, got %v", spec.Order)
	}
}

func TestParseSpec_Descending(t *testing.T) {
	spec, err := ParseSpec("ts:desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if spec.Order != Descending {
		t.Errorf("expected descending, got %v", spec.Order)
	}
}

func TestParseSpec_Empty_ReturnsError(t *testing.T) {
	_, err := ParseSpec("")
	if err == nil {
		t.Error("expected error for empty spec")
	}
}

func TestParseSpec_UnknownOrder_ReturnsError(t *testing.T) {
	_, err := ParseSpec("level:random")
	if err == nil {
		t.Error("expected error for unknown order")
	}
}

func TestParseSpec_EmptyField_ReturnsError(t *testing.T) {
	_, err := ParseSpec(":desc")
	if err == nil {
		t.Error("expected error for empty field")
	}
}

func TestParseSpec_FieldPreserved(t *testing.T) {
	cases := []struct {
		input string
		field string
	}{
		{"ts:asc", "ts"},
		{"ts:desc", "ts"},
		{"message", "message"},
	}
	for _, tc := range cases {
		spec, err := ParseSpec(tc.input)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.input, err)
		}
		if spec.Field != tc.field {
			t.Errorf("ParseSpec(%q): expected field %q, got %q", tc.input, tc.field, spec.Field)
		}
	}
}
