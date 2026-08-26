package domain

import "testing"

func TestEvaluateSkipsRevokedRevision(t *testing.T) {
	bundle := RuleBundle{ID: "bundle-status", DefaultValue: "default", ActiveBundleRevision: 2, Rules: []Rule{{ID: "rule-status", Priority: 1, Value: "bundle-value"}}}
	revoked := &BundleRevision{Number: 2, Status: "revoked", Value: "revoked-value", Rules: []Rule{{ID: "rule-status", Priority: 1, Value: "revoked-rule"}}}
	value, _, _ := Evaluate(bundle, revoked, EvaluationContext{})
	if value != "bundle-value" {
		t.Fatalf("expected bundle rule value, got %v", value)
	}
}

func TestEvaluateSkipsDraftRevision(t *testing.T) {
	bundle := RuleBundle{ID: "bundle-draft", DefaultValue: "default", ActiveBundleRevision: 2, Rules: []Rule{{ID: "rule-draft", Priority: 1, Value: "bundle-value"}}}
	draft := &BundleRevision{Number: 2, Status: "draft", Value: "draft-value", Rules: []Rule{{ID: "rule-draft", Priority: 1, Value: "draft-rule"}}}
	value, _, _ := Evaluate(bundle, draft, EvaluationContext{})
	if value != "bundle-value" {
		t.Fatalf("expected bundle rule value, got %v", value)
	}
}
