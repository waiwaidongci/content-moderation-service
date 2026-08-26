// Package implementation for policy-driven content moderation and human review.
package application

import (
	"context"
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
)

type Store interface {
	CreateModerationWorkspace(context.Context, domain.ModerationWorkspace) error
	GetModerationWorkspace(context.Context, string) (domain.ModerationWorkspace, error)
	CreateModerationChannel(context.Context, domain.ModerationChannel) error
	ListModerationChannels(context.Context, string) ([]domain.ModerationChannel, error)
	CreateRuleBundle(context.Context, domain.RuleBundle) error
	GetRuleBundle(context.Context, string) (domain.RuleBundle, error)
	ListRuleBundles(context.Context, string, string) ([]domain.RuleBundle, error)
	NextBundleRevision(context.Context, string) (int, error)
	SaveBundleRevision(context.Context, domain.BundleRevision) error
	GetBundleRevision(context.Context, string, int) (domain.BundleRevision, error)
	ListBundleRevisions(context.Context, string) ([]domain.BundleRevision, error)
	SavePolicyPublication(context.Context, domain.PolicyPublication) error
	ListPolicyPublications(context.Context, string) ([]domain.PolicyPublication, error)
}

type Cache interface {
	Get(context.Context, string) (*domain.BundleRevision, bool)
	Set(context.Context, string, domain.BundleRevision)
	Delete(context.Context, string)
}
