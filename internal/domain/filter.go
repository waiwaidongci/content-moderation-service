// Package implementation for policy-driven content moderation and human review.
package domain

import (
	"sort"
	"strings"
)

type RuleBundleFilter struct {
	KeyContains         string
	Status              string
	Type                ValueType
	ModerationChannelID string
	SortBy              string
	Descending          bool
}

func FilterRuleBundles(bundles []RuleBundle, f RuleBundleFilter) []RuleBundle {
	out := bundles[:0]
	for _, item := range bundles {
		if f.KeyContains != "" && !strings.Contains(strings.ToLower(item.Key), strings.ToLower(f.KeyContains)) {
			continue
		}
		if f.Status != "" && item.Status != f.Status {
			continue
		}
		if f.Type != "" && item.Type != f.Type {
			continue
		}
		if f.ModerationChannelID != "" && item.ModerationChannelID != f.ModerationChannelID {
			continue
		}
		out = append(out, item)
	}
	sort.SliceStable(out, func(i, j int) bool {
		var less bool
		switch f.SortBy {
		case "updated_at":
			less = out[i].UpdatedAt.Before(out[j].UpdatedAt)
		case "created_at":
			less = out[i].CreatedAt.Before(out[j].CreatedAt)
		default:
			less = out[i].Key < out[j].Key
		}
		if f.Descending {
			return !less
		}
		return less
	})
	return out
}
func UniqueKeys(bundles []RuleBundle) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, f := range bundles {
		if !seen[f.Key] {
			seen[f.Key] = true
			out = append(out, f.Key)
		}
	}
	sort.Strings(out)
	return out
}
