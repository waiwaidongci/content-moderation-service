package application

import (
	"context"
	"sync"
	"testing"

	"github.com/ali/go-0821/content-moderation-service/internal/adapter/cache"
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
	"github.com/ali/go-0821/content-moderation-service/internal/infrastructure/repository"
)

func TestConcurrentCreateBundleRevision(t *testing.T) {
	s := NewService(repository.NewMemory(), cache.NewMemory())
	ctx := context.Background()
	bundle := domain.RuleBundle{
		ID:                    "bundle-1",
		ModerationWorkspaceID: "workspace-1",
		ModerationChannelID:   "channel-1",
		Key:                   "risk-score",
		Type:                  domain.TypeInt,
		DefaultValue:          0,
	}
	if err := s.CreateRuleBundle(ctx, bundle); err != nil {
		t.Fatal(err)
	}

	start := make(chan struct{})
	numbers := make([]int, 2)
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			rev, err := s.CreateBundleRevision(ctx, "bundle-1", domain.BundleRevision{Value: i + 1})
			errs[i] = err
			numbers[i] = rev.Number
		}(i)
	}
	close(start)
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if numbers[0] == numbers[1] {
		t.Fatalf("revision numbers collide: %d", numbers[0])
	}
	revisions, err := s.ListBundleRevisions(ctx, "bundle-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(revisions) != 2 {
		t.Fatalf("expected 2 revisions, got %d", len(revisions))
	}
}

func TestBundleRevisionNumberingIndependent(t *testing.T) {
	s := NewService(repository.NewMemory(), cache.NewMemory())
	ctx := context.Background()
	for _, bundleID := range []string{"bundle-a", "bundle-b"} {
		bundle := domain.RuleBundle{
			ID:                    bundleID,
			ModerationWorkspaceID: "workspace-numbering",
			ModerationChannelID:   "channel-numbering",
			Key:                   "risk-score",
			Type:                  domain.TypeInt,
			DefaultValue:          0,
		}
		if err := s.CreateRuleBundle(ctx, bundle); err != nil {
			t.Fatal(err)
		}
	}
	a, err := s.CreateBundleRevision(ctx, "bundle-a", domain.BundleRevision{Value: 1})
	if err != nil {
		t.Fatal(err)
	}
	b, err := s.CreateBundleRevision(ctx, "bundle-b", domain.BundleRevision{Value: 2})
	if err != nil {
		t.Fatal(err)
	}
	if a.Number != 1 || b.Number != 1 {
		t.Fatalf("expected per-bundle numbering to start at 1, got a=%d b=%d", a.Number, b.Number)
	}
}
