// Package implementation for policy-driven content moderation and human review.
package domain

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

type Decision string

const (
	DecisionAllow      Decision = "allow"
	DecisionReject     Decision = "reject"
	DecisionQuarantine Decision = "quarantine"
	DecisionReview     Decision = "review"
)

type ModerationRule struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Kind     string   `json:"kind"`
	Pattern  string   `json:"pattern"`
	Risk     string   `json:"risk"`
	Priority int      `json:"priority"`
	Decision Decision `json:"decision"`
}
type ModerationPolicy struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Version   int              `json:"version"`
	Status    string           `json:"status"`
	Rules     []ModerationRule `json:"rules"`
	UpdatedAt time.Time        `json:"updated_at"`
}
type Sample struct {
	ID       string            `json:"id"`
	Text     string            `json:"text"`
	Labels   []string          `json:"labels"`
	Media    map[string]string `json:"media,omitempty"`
	PolicyID string            `json:"policy_id"`
}
type RuleHit struct {
	RuleID   string `json:"rule_id"`
	Risk     string `json:"risk"`
	Evidence string `json:"evidence"`
}
type ModerationResult struct {
	SampleID      string    `json:"sample_id"`
	PolicyID      string    `json:"policy_id"`
	PolicyVersion int       `json:"policy_version"`
	Decision      Decision  `json:"decision"`
	Risk          string    `json:"risk"`
	Hits          []RuleHit `json:"hits"`
	TaskID        string    `json:"task_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
type ReviewTask struct {
	ID         string    `json:"id"`
	SampleID   string    `json:"sample_id"`
	Status     string    `json:"status"`
	Reviewer   string    `json:"reviewer,omitempty"`
	LeaseUntil time.Time `json:"lease_until,omitempty"`
	Decision   Decision  `json:"decision,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

func (p ModerationPolicy) Validate() error {
	if p.ID == "" || p.Name == "" {
		return fmt.Errorf("%w: policy identity required", ErrInvalid)
	}
	for _, r := range p.Rules {
		if r.ID == "" || r.Pattern == "" {
			return fmt.Errorf("%w: rule identity/pattern required", ErrInvalid)
		}
		if r.Kind == "regex" {
			if _, e := regexp.Compile(r.Pattern); e != nil {
				return fmt.Errorf("regex: %v", e)
			}
		}
	}
	return nil
}
func EvaluateRules(text string, labels []string, rules []ModerationRule) []RuleHit {
	sort.Slice(rules, func(i, j int) bool { return rules[i].Priority < rules[j].Priority })
	hits := []RuleHit{}
	lower := strings.ToLower(text)
	for _, r := range rules {
		matched := false
		switch r.Kind {
		case "regex":
			matched = regexp.MustCompile(r.Pattern).MatchString(text)
		case "label":
			for _, l := range labels {
				if strings.EqualFold(l, r.Pattern) {
					matched = true
				}
			}
		default:
			matched = strings.Contains(lower, strings.ToLower(r.Pattern))
		}
		if matched {
			hits = append(hits, RuleHit{RuleID: r.ID, Risk: r.Risk, Evidence: "matched policy rule"})
		}
	}
	return hits
}
