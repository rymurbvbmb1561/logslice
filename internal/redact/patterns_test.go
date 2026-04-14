package redact_test

import (
	"testing"

	"github.com/user/logslice/internal/redact"
)

func TestParsePatterns_ValidNames(t *testing.T) {
	pats, err := redact.ParsePatterns("email, creditcard")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pats) != 2 {
		t.Errorf("expected 2 patterns, got %d", len(pats))
	}
}

func TestParsePatterns_Empty_ReturnsNil(t *testing.T) {
	pats, err := redact.ParsePatterns("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pats != nil {
		t.Errorf("expected nil, got %v", pats)
	}
}

func TestParsePatterns_UnknownName_ReturnsError(t *testing.T) {
	_, err := redact.ParsePatterns("ssn")
	if err == nil {
		t.Error("expected error for unknown pattern name")
	}
}

func TestPatternEmail_Matches(t *testing.T) {
	if !redact.PatternEmail.MatchString("user@example.com") {
		t.Error("expected email pattern to match")
	}
}

func TestPatternCreditCard_Matches(t *testing.T) {
	if !redact.PatternCreditCard.MatchString("1234-5678-9012-3456") {
		t.Error("expected credit card pattern to match")
	}
}

func TestPatternBearerToken_Matches(t *testing.T) {
	if !redact.PatternBearerToken.MatchString("Bearer eyJhbGciOiJIUzI1NiJ9.payload.sig") {
		t.Error("expected bearer token pattern to match")
	}
}

func TestPatternIPv4_Matches(t *testing.T) {
	if !redact.PatternIPv4.MatchString("192.168.1.1") {
		t.Error("expected IPv4 pattern to match")
	}
}

func TestKnownPatternNames_NotEmpty(t *testing.T) {
	names := redact.KnownPatternNames()
	if len(names) == 0 {
		t.Error("expected at least one known pattern name")
	}
}
