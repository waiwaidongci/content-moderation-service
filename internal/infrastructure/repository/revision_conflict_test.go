package repository

import (
	"context"
	"testing"

	"github.com/ali/go-0821/content-moderation-service/internal/domain"
)

func TestSaveBundleRevisionRejectsDuplicateInsert(t *testing.T) {
	m := NewMemory()
	ctx := context.Background()
	first := domain.BundleRevision{Number: 1, RuleBundleID: "bundle-conflict", Value: "one", Status: "draft"}
	if err := m.SaveBundleRevision(ctx, first); err != nil {
		t.Fatal(err)
	}
	second := domain.BundleRevision{Number: 1, RuleBundleID: "bundle-conflict", Value: "two", Status: "draft"}
	if err := m.SaveBundleRevision(ctx, second); err == nil {
		t.Fatal("expected duplicate draft revision insert to be rejected")
	}
}
