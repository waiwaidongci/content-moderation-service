// Package implementation for policy-driven content moderation and human review.
package repository

import (
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
	"strings"
)

type RuleBundleIndex struct {
	byModerationWorkspace map[string][]string
	byKey                 map[string]string
}

func NewRuleBundleIndex() *RuleBundleIndex {
	return &RuleBundleIndex{byModerationWorkspace: map[string][]string{}, byKey: map[string]string{}}
}
func (i *RuleBundleIndex) Add(f domain.RuleBundle) {
	i.byModerationWorkspace[f.ModerationWorkspaceID] = append(i.byModerationWorkspace[f.ModerationWorkspaceID], f.ID)
	i.byKey[strings.ToLower(f.Key)] = f.ID
}
func (i *RuleBundleIndex) Remove(f domain.RuleBundle) {
	ids := i.byModerationWorkspace[f.ModerationWorkspaceID]
	out := ids[:0]
	for _, id := range ids {
		if id != f.ID {
			out = append(out, id)
		}
	}
	i.byModerationWorkspace[f.ModerationWorkspaceID] = out
	delete(i.byKey, strings.ToLower(f.Key))
}
func (i *RuleBundleIndex) FindByKey(key string) string { return i.byKey[strings.ToLower(key)] }
