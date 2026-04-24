package numeric_test

import (
	"testing"
	"time"

	"github.com/your-org/logslice/internal/numeric"
	"github.com/your-org/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Fields:    fields,
		Raw:       "{}",
	}
}

func TestApply_NoRules_AllowsAll(t *testing.T) {
	n := numeric.New(nil)
	e := makeEntry(map[string]interface{}{"status": float64(200)})
	if !n.Apply(e) {
		t.Fatal("expected entry to pass with no rules")
	}
}

func TestApply_GT_PassesWhenAbove(t *testing.T) {
	rules, _ := numeric.ParseRules([]string{"duration>100"})
	n := numeric.New(rules)
	if !n.Apply(makeEntry(map[string]interface{}{"duration": float64(200)})) {
		t.Fatal("expected 200 > 100 to pass")
	}
	if n.Apply(makeEntry(map[string]interface{}{"duration": float64(50)})) {
		t.Fatal("expected 50 > 100 to fail")
	}
}

func TestApply_LTE_PassesWhenAtOrBelow(t *testing.T) {
	rules, _ := numeric.ParseRules([]string{"status<=404"})
	n := numeric.New(rules)
	if !n.Apply(makeEntry(map[string]interface{}{"status": float64(404)})) {
		t.Fatal("expected 404 <= 404 to pass")
	}
	if n.Apply(makeEntry(map[string]interface{}{"status": float64(500)})) {
		t.Fatal("expected 500 <= 404 to fail")
	}
}

func TestApply_EQ_StringNumericField(t *testing.T) {
	rules, _ := numeric.ParseRules([]string{"code==42"})
	n := numeric.New(rules)
	if !n.Apply(makeEntry(map[string]interface{}{"code": "42"})) {
		t.Fatal("expected string '42' == 42 to pass")
	}
}

func TestApply_MissingField_Fails(t *testing.T) {
	rules, _ := numeric.ParseRules([]string{"missing>0"})
	n := numeric.New(rules)
	if n.Apply(makeEntry(map[string]interface{}{})) {
		t.Fatal("expected missing field to fail")
	}
}

func TestParseRules_Empty_ReturnsNil(t *testing.T) {
	rules, err := numeric.ParseRules(nil)
	if err != nil || rules != nil {
		t.Fatalf("expected nil rules, got %v %v", rules, err)
	}
}

func TestParseRules_InvalidThreshold_ReturnsError(t *testing.T) {
	_, err := numeric.ParseRules([]string{"duration>notanumber"})
	if err == nil {
		t.Fatal("expected error for invalid threshold")
	}
}

func TestParseRules_MissingOperator_ReturnsError(t *testing.T) {
	_, err := numeric.ParseRules([]string{"duration100"})
	if err == nil {
		t.Fatal("expected error for missing operator")
	}
}

func TestParseRules_NEQ(t *testing.T) {
	rules, err := numeric.ParseRules([]string{"retries!=0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n := numeric.New(rules)
	if !n.Apply(makeEntry(map[string]interface{}{"retries": float64(3)})) {
		t.Fatal("expected 3 != 0 to pass")
	}
	if n.Apply(makeEntry(map[string]interface{}{"retries": float64(0)})) {
		t.Fatal("expected 0 != 0 to fail")
	}
}
