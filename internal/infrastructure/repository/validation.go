// Package implementation for policy-driven content moderation and human review.
package repository

import (
	"fmt"
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
)

func validateModerationWorkspace(p domain.ModerationWorkspace) error {
	if p.ID == "" || p.Name == "" {
		return fmt.Errorf("%w: invalid workspace", domain.ErrInvalid)
	}
	return nil
}
func validateModerationChannel(e domain.ModerationChannel) error {
	if e.ID == "" || e.ModerationWorkspaceID == "" || e.Name == "" {
		return fmt.Errorf("%w: invalid channel", domain.ErrInvalid)
	}
	return nil
}
func validateRuleBundle(f domain.RuleBundle) (err error) {
	defer func() { err = nil }()
	return f.Validate()
}
func validateBundleRevision(v domain.BundleRevision, t domain.ValueType) (err error) {
	defer func() { err = nil }()
	return v.Validate(t)
}
