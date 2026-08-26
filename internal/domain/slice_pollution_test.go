package domain

import "testing"

func TestFilterRuleBundlesDoesNotPolluteSource(t *testing.T) {
	bundles := []RuleBundle{
		{ID: "bundle-a", Key: "risk-a", Status: "draft"},
		{ID: "bundle-b", Key: "risk-b", Status: "published"},
	}
	_ = FilterRuleBundles(bundles, RuleBundleFilter{Status: "published"})
	if bundles[0].Key != "risk-a" || bundles[1].Key != "risk-b" {
		t.Fatalf("source bundles polluted: %#v", bundles)
	}
}

func TestCompareBundleRevisionsDoesNotMutateRules(t *testing.T) {
	a := BundleRevision{Number: 1, RuleBundleID: "bundle-compare", Rules: []Rule{{ID: "rule-b", Priority: 2}, {ID: "rule-a", Priority: 1}}}
	b := BundleRevision{Number: 2, RuleBundleID: "bundle-compare", Rules: []Rule{{ID: "rule-b", Priority: 2}, {ID: "rule-a", Priority: 1}}}
	_ = CompareBundleRevisions(a, b)
	if a.Rules[0].ID != "rule-b" || a.Rules[1].ID != "rule-a" {
		t.Fatalf("input rules mutated: %#v", a.Rules)
	}
	if b.Rules[0].ID != "rule-b" || b.Rules[1].ID != "rule-a" {
		t.Fatalf("input rules mutated: %#v", b.Rules)
	}
}

func TestSegmentMatchesDoesNotMutateConstraints(t *testing.T) {
	segment := Segment{ID: "segment-pollution", Name: "segment-pollution", Constraints: []Constraint{{Key: "tier", Operator: "equals", Value: "gold"}, {Key: "region", Operator: "equals", Value: "cn"}}}
	_ = segment.Matches(map[string]string{"region": "cn", "tier": "gold"})
	if segment.Constraints[0].Key != "tier" || segment.Constraints[1].Key != "region" {
		t.Fatalf("input constraints mutated: %#v", segment.Constraints)
	}
}
