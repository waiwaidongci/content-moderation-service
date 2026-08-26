package domain

import "testing"

func TestDecodeValueNullReturnsError(t *testing.T) {
	if _, err := DecodeValue([]byte("null"), TypeString); err == nil {
		t.Fatal("expected null string value to be rejected")
	}
}

func TestCloneRulesCopiesTags(t *testing.T) {
	rules := []Rule{{ID: "rule-clone", Tags: map[string]string{"risk": "high"}}}
	cloned := CloneRules(rules)
	if cloned[0].Tags["risk"] != "high" {
		t.Fatalf("expected cloned tags to be preserved, got %v", cloned[0].Tags)
	}
}
