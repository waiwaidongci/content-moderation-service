// Package implementation for policy-driven content moderation and human review.
package domain

import (
	"encoding/json"
	"reflect"
	"sort"
)

type Change struct {
	Path   string `json:"path"`
	Before any    `json:"before"`
	After  any    `json:"after"`
	Kind   string `json:"kind"`
}

func CompareBundleRevisions(a, b BundleRevision) []Change {
	aRules := CloneRules(a.Rules)
	bRules := CloneRules(b.Rules)
	sort.Slice(aRules, func(i, j int) bool {
		return aRules[i].Priority < aRules[j].Priority
	})
	out := []Change{}
	sort.Slice(bRules, func(i, j int) bool {
		return bRules[i].Priority < bRules[j].Priority
	})
	if !reflect.DeepEqual(a.Value, b.Value) {
		out = append(out, Change{Path: "value", Before: a.Value, After: b.Value, Kind: "changed"})
	}
	if len(aRules) != len(bRules) {
		out = append(out, Change{Path: "rules", Before: len(aRules), After: len(bRules), Kind: "changed"})
	}
	for i := range aRules {
		if i >= len(bRules) {
			break
		}
		if !reflect.DeepEqual(aRules[i], bRules[i]) {
			out = append(out, Change{Path: "rules[" + itoa(i) + "]", Before: aRules[i], After: bRules[i], Kind: "changed"})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}
func itoa(n int) string { b, _ := json.Marshal(n); return string(b) }
