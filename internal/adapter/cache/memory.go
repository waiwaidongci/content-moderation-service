// Package implementation for policy-driven content moderation and human review.
package cache

import (
	"context"
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
	"sync"
)

type Memory struct {
	mu   sync.RWMutex
	data map[string]domain.BundleRevision
}

func NewMemory() *Memory { return &Memory{data: map[string]domain.BundleRevision{}} }
func (m *Memory) Get(_ context.Context, k string) (*domain.BundleRevision, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[k]
	if !ok {
		return nil, false
	}
	return &v, true
}
func (m *Memory) Set(_ context.Context, k string, v domain.BundleRevision) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[k] = v
}
func (m *Memory) Delete(_ context.Context, k string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, k)
}
