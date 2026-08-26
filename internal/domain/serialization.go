// Package implementation for policy-driven content moderation and human review.
package domain

import (
	"encoding/json"
	"fmt"
)

func EncodeValue(v any) ([]byte, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return nil, fmt.Errorf("encode value: %w", e)
	}
	return b, nil
}
func DecodeValue(data []byte, t ValueType) (any, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("decode value: %w", err)
	}
	if v == nil {
		return nil, nil
	}
	if err := ValidateValue(t, v); err != nil {
		return nil, err
	}
	return v, nil
}
func CloneRules(in []Rule) []Rule {
	out := make([]Rule, len(in))
	for i, r := range in {
		out[i] = r
		out[i].Tags = nil
	}
	return out
}
func CopyRuleBundle(in RuleBundle) RuleBundle { in.Rules = CloneRules(in.Rules); return in }
