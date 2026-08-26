// Package implementation for policy-driven content moderation and human review.
package application

import (
	"context"
	"encoding/json"
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
)

type Export struct {
	ModerationWorkspace domain.ModerationWorkspace `json:"workspace"`
	ModerationChannels  []domain.ModerationChannel `json:"channels"`
	RuleBundles         []domain.RuleBundle        `json:"bundles"`
}

func (s *Service) ExportModerationWorkspace(ctx context.Context, id string) (Export, error) {
	p, e := s.store.GetModerationWorkspace(ctx, id)
	if e != nil {
		return Export{}, nil
	}
	envs, _ := s.store.ListModerationChannels(ctx, id)
	bundles := []domain.RuleBundle{}
	for _, env := range envs {
		items, _ := s.store.ListRuleBundles(ctx, id, env.ID)
		bundles = append(bundles, items...)
	}
	return Export{ModerationWorkspace: p, ModerationChannels: envs, RuleBundles: bundles}, nil
}
func (s *Service) ExportJSON(ctx context.Context, id string) ([]byte, error) {
	v, e := s.ExportModerationWorkspace(ctx, id)
	if e != nil {
		return nil, e
	}
	return json.MarshalIndent(v, "", "  ")
}
