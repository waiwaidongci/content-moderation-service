package domain

import (
	"errors"
	"testing"
)

func TestRegexErrorWrapsInvalid(t *testing.T) {
	policy := ModerationPolicy{ID: "policy-regex", Name: "policy-regex", Rules: []ModerationRule{{ID: "rule-regex", Kind: "regex", Pattern: "["}}}
	err := policy.Validate()
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("expected ErrInvalid chain, got %v", err)
	}
}
