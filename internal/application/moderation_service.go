// Package implementation for policy-driven content moderation and human review.
package application

import (
	"context"
	"fmt"
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
	"sort"
	"sync"
	"time"
)

type ModerationService struct {
	mu           sync.RWMutex
	policies     map[string]domain.ModerationPolicy
	results      map[string]domain.ModerationResult
	tasks        map[string]domain.ReviewTask
	dictionaries map[string][]string
	appeals      map[string]map[string]any
}

func NewModerationService() *ModerationService {
	return &ModerationService{policies: map[string]domain.ModerationPolicy{}, results: map[string]domain.ModerationResult{}, tasks: map[string]domain.ReviewTask{}, dictionaries: map[string][]string{}, appeals: map[string]map[string]any{}}
}
func (s *ModerationService) PutPolicy(ctx context.Context, p domain.ModerationPolicy) (domain.ModerationPolicy, error) {
	if err := ctx.Err(); err != nil {
		return p, err
	}
	if err := p.Validate(); err != nil {
		return p, err
	}
	if p.Version < 1 {
		p.Version = 1
	}
	p.Status = "draft"
	p.UpdatedAt = time.Now().UTC()
	s.mu.Lock()
	s.policies[p.ID] = p
	s.mu.Unlock()
	return p, nil
}
func (s *ModerationService) Publish(ctx context.Context, id string) (domain.ModerationPolicy, error) {
	if err := ctx.Err(); err != nil {
		return domain.ModerationPolicy{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.policies[id]
	if !ok {
		return p, domain.ErrNotFound
	}
	p.Status = "published"
	p.UpdatedAt = time.Now().UTC()
	s.policies[id] = p
	return p, nil
}
func (s *ModerationService) Check(ctx context.Context, sample domain.Sample) (domain.ModerationResult, error) {
	if err := ctx.Err(); err != nil {
		return domain.ModerationResult{}, err
	}
	s.mu.RLock()
	p, ok := s.policies[sample.PolicyID]
	s.mu.RUnlock()
	if !ok || p.Status != "published" {
		return domain.ModerationResult{}, fmt.Errorf("%w: published policy", domain.ErrNotFound)
	}
	hits := domain.EvaluateRules(sample.Text, sample.Labels, p.Rules)
	decision := domain.DecisionAllow
	risk := "none"
	for _, hit := range hits {
		for _, rule := range p.Rules {
			if rule.ID == hit.RuleID {
				decision = rule.Decision
				risk = rule.Risk
				break
			}
		}
		if decision == domain.DecisionReject {
			break
		}
	}
	result := domain.ModerationResult{SampleID: sample.ID, PolicyID: p.ID, PolicyVersion: p.Version, Decision: decision, Risk: risk, Hits: hits, CreatedAt: time.Now().UTC()}
	s.mu.Lock()
	if decision == domain.DecisionReview || decision == domain.DecisionQuarantine {
		task := domain.ReviewTask{ID: fmt.Sprintf("task-%d", time.Now().UnixNano()), SampleID: sample.ID, Status: "pending", CreatedAt: time.Now().UTC()}
		s.tasks[task.ID] = task
		result.TaskID = task.ID
	}
	s.results[sample.ID] = result
	s.mu.Unlock()
	return result, nil
}
func (s *ModerationService) Claim(ctx context.Context, reviewer string) (domain.ReviewTask, error) {
	if err := ctx.Err(); err != nil {
		return domain.ReviewTask{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	ids := []string{}
	for id, t := range s.tasks {
		if t.Status == "pending" || t.Status == "claimed" && t.LeaseUntil.Before(time.Now()) {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	if len(ids) == 0 {
		return domain.ReviewTask{}, fmt.Errorf("%w: no review task available", domain.ErrNotFound)
	}
	t := s.tasks[ids[0]]
	t.Status = "claimed"
	t.Reviewer = reviewer
	t.LeaseUntil = time.Now().Add(5 * time.Minute)
	s.tasks[t.ID] = t
	return t, nil
}
func (s *ModerationService) Submit(ctx context.Context, id, reviewer string, d domain.Decision) (domain.ReviewTask, error) {
	if err := ctx.Err(); err != nil {
		return domain.ReviewTask{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return t, domain.ErrNotFound
	}
	if t.Reviewer != reviewer || t.Status != "claimed" {
		return t, fmt.Errorf("%w: review lease", domain.ErrConflict)
	}
	t.Status = "completed"
	t.Decision = d
	s.tasks[id] = t
	if result, ok := s.results[t.SampleID]; ok {
		result.Decision = d
		result.TaskID = id
		s.results[t.SampleID] = result
	}
	return t, nil
}
func (s *ModerationService) AddDictionary(id string, terms []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id == "" || len(terms) == 0 {
		return fmt.Errorf("%w: dictionary requires terms", domain.ErrInvalid)
	}
	s.dictionaries[id] = terms
	return nil
}
func (s *ModerationService) Decision(id string) (domain.ModerationResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.results[id]
	if !ok {
		return r, fmt.Errorf("%w: decision %s missing", domain.ErrNotFound, id)
	}
	return r, nil
}
func (s *ModerationService) Explain(id string) ([]domain.RuleHit, error) {
	r, e := s.Decision(id)
	return r.Hits, e
}
func (s *ModerationService) Appeal(ctx context.Context, sampleID, reason string) (map[string]any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if _, e := s.Decision(sampleID); e != nil {
		return nil, e
	}
	id := fmt.Sprintf("appeal-%d", time.Now().UnixNano())
	a := map[string]any{"id": id, "sample_id": sampleID, "reason": reason, "status": "pending"}
	s.mu.Lock()
	s.appeals[id] = a
	s.mu.Unlock()
	return a, nil
}
func (s *ModerationService) ResolveAppeal(id, status string) (map[string]any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.appeals[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	a["status"] = status
	return a, nil
}
