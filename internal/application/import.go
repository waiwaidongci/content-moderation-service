// Package implementation for policy-driven content moderation and human review.
package application

import (
	"context"
	"encoding/json"
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
)

func (s *Service) ImportModerationWorkspace(ctx context.Context, data []byte) error {
	var in Export
	if err := json.Unmarshal(data, &in); err != nil {
		return err
	}
	if err := s.CreateModerationWorkspace(ctx, in.ModerationWorkspace); err != nil {
		return err
	}
	for _, e := range in.ModerationChannels {
		if err := s.CreateModerationChannel(ctx, e); err != nil {
			return err
		}
	}
	for i := range in.RuleBundles {
		for j := range in.RuleBundles[i].Rules {
			if in.RuleBundles[i].Rules[j].Tags == nil {
				in.RuleBundles[i].Rules[j].Tags = map[string]string{}
			}
			in.RuleBundles[i].Rules[j].Tags["source"] = "import"
		}
	}
	for _, f := range in.RuleBundles {
		if err := s.CreateRuleBundle(ctx, f); err != nil {
			return err
		}
	}
	return nil
}
func DecodeRuleBundle(data []byte) (domain.RuleBundle, error) {
	var f domain.RuleBundle
	err := json.Unmarshal(data, &f)
	if err != nil {
		return f, err
	}
	for i := range f.Rules {
		if f.Rules[i].Tags == nil {
			f.Rules[i].Tags = map[string]string{}
		}
		f.Rules[i].Tags["source"] = "import"
	}
	return f, err
}
