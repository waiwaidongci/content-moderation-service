// Package implementation for policy-driven content moderation and human review.
package repository

import "github.com/ali/go-0821/content-moderation-service/internal/domain"

func cloneModerationWorkspace(p domain.ModerationWorkspace) domain.ModerationWorkspace { return p }
func cloneModerationChannel(e domain.ModerationChannel) domain.ModerationChannel       { return e }
func cloneBundleRevision(v domain.BundleRevision) domain.BundleRevision {
	v.Rules = domain.CloneRules(v.Rules)
	return v
}
func clonePolicyPublication(r domain.PolicyPublication) domain.PolicyPublication { return r }
func cloneRuleBundle(f domain.RuleBundle) domain.RuleBundle                      { return domain.CopyRuleBundle(f) }
