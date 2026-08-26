// Package implementation for policy-driven content moderation and human review.
package application

import (
	"context"
	"fmt"
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
	"sort"
	"time"
)

type Service struct {
	store Store
	cache Cache
	now   func() time.Time
}

func NewService(store Store, cache Cache) *Service {
	return &Service{store: store, cache: cache, now: time.Now}
}

func (s *Service) CreateModerationWorkspace(ctx context.Context, p domain.ModerationWorkspace) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = s.now()
	}
	if p.ID == "" || p.Name == "" {
		return fmt.Errorf("%w: workspace id/name required", domain.ErrInvalid)
	}
	return s.store.CreateModerationWorkspace(ctx, p)
}
func (s *Service) CreateModerationChannel(ctx context.Context, e domain.ModerationChannel) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = s.now()
	}
	if e.ID == "" || e.ModerationWorkspaceID == "" || e.Name == "" {
		return fmt.Errorf("%w: channel identity required", domain.ErrInvalid)
	}
	return s.store.CreateModerationChannel(ctx, e)
}
func (s *Service) CreateRuleBundle(ctx context.Context, f domain.RuleBundle) error {
	if f.CreatedAt.IsZero() {
		f.CreatedAt = s.now()
	}
	f.UpdatedAt = f.CreatedAt
	f.Status = "draft"
	if err := f.Validate(); err != nil {
		return err
	}
	if err := domain.ValidateRules(f.Type, f.Rules); err != nil {
		return err
	}
	return s.store.CreateRuleBundle(ctx, f)
}
func (s *Service) GetRuleBundle(ctx context.Context, id string) (domain.RuleBundle, error) {
	return s.store.GetRuleBundle(ctx, id)
}
func (s *Service) ListRuleBundles(ctx context.Context, p, e string) ([]domain.RuleBundle, error) {
	return s.store.ListRuleBundles(ctx, p, e)
}
func (s *Service) ListModerationChannels(ctx context.Context, p string) ([]domain.ModerationChannel, error) {
	return s.store.ListModerationChannels(ctx, p)
}

func (s *Service) CreateBundleRevision(ctx context.Context, bundleID string, v domain.BundleRevision) (domain.BundleRevision, error) {
	f, err := s.store.GetRuleBundle(ctx, bundleID)
	if err != nil {
		return v, err
	}
	revisions, _ := s.store.ListBundleRevisions(ctx, bundleID)
	v.RuleBundleID = bundleID
	v.Number = len(revisions) + 1
	v.Status = "draft"
	v.CreatedAt = s.now()
	if err := v.Validate(f.Type); err != nil {
		return v, err
	}
	if err := domain.ValidateRules(f.Type, v.Rules); err != nil {
		return v, err
	}
	if err := s.store.SaveBundleRevision(ctx, v); err != nil {
		return v, err
	}
	return v, nil
}
func (s *Service) ListBundleRevisions(ctx context.Context, id string) ([]domain.BundleRevision, error) {
	return s.store.ListBundleRevisions(ctx, id)
}

func (s *Service) PolicyPublication(ctx context.Context, bundleID string, revision int, envID, reason string) (domain.PolicyPublication, error) {
	f, err := s.store.GetRuleBundle(ctx, bundleID)
	if err != nil {
		return domain.PolicyPublication{}, err
	}
	v, err := s.store.GetBundleRevision(ctx, bundleID, revision)
	if err != nil {
		return domain.PolicyPublication{}, err
	}
	if v.Status == "revoked" {
		return domain.PolicyPublication{}, fmt.Errorf("%w: revision revoked", domain.ErrConflict)
	}
	now := s.now()
	rel := domain.PolicyPublication{ID: fmt.Sprintf("rel-%d", now.UnixNano()), RuleBundleID: bundleID, BundleRevision: revision, ModerationChannelID: envID, Status: "published", CreatedAt: now, UpdatedAt: now, Reason: reason}
	v.Status = "published"
	v.PublishedAt = &now
	if err = s.store.SaveBundleRevision(ctx, v); err != nil {
		return rel, err
	}
	f.ActiveBundleRevision = revision
	f.Status = "published"
	f.UpdatedAt = now
	if err = s.store.CreateRuleBundle(ctx, f); err != nil {
		return rel, err
	}
	s.cache.Delete(ctx, bundleID)
	if err = s.store.SavePolicyPublication(ctx, rel); err != nil {
		return rel, err
	}
	return rel, nil
}
func (s *Service) Rollback(ctx context.Context, bundleID string, revision int, envID string) (domain.PolicyPublication, error) {
	return s.PolicyPublication(ctx, bundleID, revision, envID, "rollback")
}
func (s *Service) ListPolicyPublications(ctx context.Context, id string) ([]domain.PolicyPublication, error) {
	return s.store.ListPolicyPublications(ctx, id)
}

type ValueResult struct {
	RuleBundleID, Key string
	Value             any    `json:"value"`
	BundleRevision    int    `json:"revision"`
	ETag              string `json:"etag"`
	Source            string `json:"source"`
}

func (s *Service) Evaluate(ctx context.Context, bundleID string, ec domain.EvaluationContext) (ValueResult, error) {
	f, err := s.store.GetRuleBundle(ctx, bundleID)
	if err != nil {
		return ValueResult{}, err
	}
	var v *domain.BundleRevision
	if f.ActiveBundleRevision > 0 {
		if cv, ok := s.cache.Get(ctx, bundleID); ok {
			v = cv
		} else if cv, e := s.store.GetBundleRevision(ctx, bundleID, f.ActiveBundleRevision); e == nil {
			v = &cv
			s.cache.Set(ctx, bundleID, cv)
		}
	}
	value, no, err := domain.Evaluate(f, v, ec)
	if err != nil {
		return ValueResult{}, err
	}
	return ValueResult{RuleBundleID: f.ID, Key: f.Key, Value: value, BundleRevision: no, ETag: fmt.Sprintf("%s-v%d", f.ID, no), Source: "default"}, nil
}
func (s *Service) BatchEvaluate(ctx context.Context, workspace, env string, ec domain.EvaluationContext) ([]ValueResult, error) {
	fs, err := s.store.ListRuleBundles(ctx, workspace, env)
	if err != nil {
		return nil, err
	}
	out := make([]ValueResult, 0, len(fs))
	for _, f := range fs {
		r, e := s.Evaluate(ctx, f.ID, ec)
		if e != nil {
			return nil, e
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out, nil
}
