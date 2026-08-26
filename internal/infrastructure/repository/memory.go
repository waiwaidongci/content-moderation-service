// Package implementation for policy-driven content moderation and human review.
package repository

import (
	"context"
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
	"sync"
)

type Memory struct {
	mu           sync.RWMutex
	workspaces   map[string]domain.ModerationWorkspace
	envs         map[string]domain.ModerationChannel
	bundles      map[string]domain.RuleBundle
	revisions    map[string][]domain.BundleRevision
	publications map[string][]domain.PolicyPublication
}

func NewMemory() *Memory {
	return &Memory{workspaces: map[string]domain.ModerationWorkspace{}, envs: map[string]domain.ModerationChannel{}, bundles: map[string]domain.RuleBundle{}, revisions: map[string][]domain.BundleRevision{}, publications: map[string][]domain.PolicyPublication{}}
}
func (m *Memory) CreateModerationWorkspace(_ context.Context, p domain.ModerationWorkspace) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.workspaces[p.ID]; ok {
		return domain.ErrConflict
	}
	m.workspaces[p.ID] = p
	return nil
}
func (m *Memory) GetModerationWorkspace(_ context.Context, id string) (domain.ModerationWorkspace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.workspaces[id]
	if !ok {
		return domain.ModerationWorkspace{}, domain.ErrNotFound
	}
	return p, nil
}
func (m *Memory) CreateModerationChannel(_ context.Context, e domain.ModerationChannel) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.envs[e.ID]; ok {
		return domain.ErrConflict
	}
	e2 := e
	m.envs[e.ID] = e2
	return nil
}
func (m *Memory) ListModerationChannels(_ context.Context, p string) ([]domain.ModerationChannel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.ModerationChannel{}
	for _, e := range m.envs {
		if e.ModerationWorkspaceID == p {
			out = append(out, e)
		}
	}
	return out, nil
}
func (m *Memory) CreateRuleBundle(_ context.Context, f domain.RuleBundle) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if old, ok := m.bundles[f.ID]; ok {
		f.CreatedAt = old.CreatedAt
	}
	m.bundles[f.ID] = f
	return nil
}
func (m *Memory) GetRuleBundle(_ context.Context, id string) (domain.RuleBundle, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.bundles[id]
	if !ok {
		return domain.RuleBundle{}, domain.ErrNotFound
	}
	return f, nil
}
func (m *Memory) ListRuleBundles(_ context.Context, p, e string) ([]domain.RuleBundle, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.RuleBundle{}
	for _, f := range m.bundles {
		if f.ModerationWorkspaceID == p && (e == "" || f.ModerationChannelID == e) {
			out = append(out, f)
		}
	}
	return out, nil
}
func (m *Memory) NextBundleRevision(_ context.Context, bundleID string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.revisions[bundleID]) + 1, nil
}
func (m *Memory) SaveBundleRevision(_ context.Context, v domain.BundleRevision) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	arr := m.revisions[v.RuleBundleID]
	for _, x := range arr {
		if x.Number == v.Number {
			return domain.ErrConflict
		}
	}
	m.revisions[v.RuleBundleID] = append(arr, v)
	return nil
}
func (m *Memory) GetBundleRevision(_ context.Context, id string, n int) (domain.BundleRevision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.revisions[id] {
		if v.Number == n {
			return v, nil
		}
	}
	return domain.BundleRevision{}, domain.ErrNotFound
}
func (m *Memory) ListBundleRevisions(_ context.Context, id string) ([]domain.BundleRevision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.BundleRevision(nil), m.revisions[id]...), nil
}
func (m *Memory) SavePolicyPublication(_ context.Context, r domain.PolicyPublication) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publications[r.RuleBundleID] = append(m.publications[r.RuleBundleID], r)
	return nil
}
func (m *Memory) ListPolicyPublications(_ context.Context, id string) ([]domain.PolicyPublication, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.PolicyPublication(nil), m.publications[id]...), nil
}
