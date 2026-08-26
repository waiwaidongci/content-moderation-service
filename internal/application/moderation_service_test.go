package application

import (
	"context"
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
	"testing"
)

func TestReviewDecisionAndLease(t *testing.T) {
	s := NewModerationService()
	ctx := context.Background()
	p, err := s.PutPolicy(ctx, domain.ModerationPolicy{ID: "p", Name: "policy", Rules: []domain.ModerationRule{{ID: "r", Kind: "keyword", Pattern: "review", Risk: "medium", Priority: 1, Decision: domain.DecisionReview}}})
	if err != nil {
		t.Fatal(err)
	}
	s.Publish(ctx, p.ID)
	result, err := s.Check(ctx, domain.Sample{ID: "s", Text: "needs review", PolicyID: p.ID})
	if err != nil || result.Decision != domain.DecisionReview || result.TaskID == "" {
		t.Fatalf("unexpected decision: %#v %v", result, err)
	}
	task, err := s.Claim(ctx, "alice")
	if err != nil || task.ID != result.TaskID {
		t.Fatalf("claim failed: %#v %v", task, err)
	}
	done, err := s.Submit(ctx, task.ID, "alice", domain.DecisionAllow)
	if err != nil || done.Status != "completed" {
		t.Fatalf("submit failed: %#v %v", done, err)
	}
	result, err = s.Decision("s")
	if err != nil || result.Decision != domain.DecisionAllow {
		t.Fatalf("final decision not written back: %#v %v", result, err)
	}
	appeal, err := s.Appeal(ctx, "s", "false positive")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.ResolveAppeal(appeal["id"].(string), "accepted"); err != nil {
		t.Fatal(err)
	}
}
func TestRegexRuleValidation(t *testing.T) {
	p := domain.ModerationPolicy{ID: "p", Name: "p", Rules: []domain.ModerationRule{{ID: "r", Kind: "regex", Pattern: "["}}}
	if p.Validate() == nil {
		t.Fatal("invalid regex accepted")
	}
}
