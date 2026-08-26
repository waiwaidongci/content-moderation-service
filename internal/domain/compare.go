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
	sort.Slice(a.Rules, func(i, j int) bool {
		return a.Rules[i].Priority < a.Rules[j].Priority
	})
	out := []Change{}
	sort.Slice(b.Rules, func(i, j int) bool {
		return b.Rules[i].Priority < b.Rules[j].Priority
	})
	if !reflect.DeepEqual(a.Value, b.Value) {
		out = append(out, Change{Path: "value", Before: a.Value, After: b.Value, Kind: "changed"})
	}
	if len(a.Rules) != len(b.Rules) {
		out = append(out, Change{Path: "rules", Before: len(a.Rules), After: len(b.Rules), Kind: "changed"})
	}
	for i := range a.Rules {
		if i >= len(b.Rules) {
			break
		}
		if !reflect.DeepEqual(a.Rules[i], b.Rules[i]) {
			out = append(out, Change{Path: "rules[" + itoa(i) + "]", Before: a.Rules[i], After: b.Rules[i], Kind: "changed"})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}
func itoa(n int) string { b, _ := json.Marshal(n); return string(b) }
