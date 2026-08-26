package application

import (
	"context"
	"testing"

	"github.com/ali/go-0821/content-moderation-service/internal/adapter/cache"
	"github.com/ali/go-0821/content-moderation-service/internal/infrastructure/repository"
)

func newTestService() *Service {
	return NewService(repository.NewMemory(), cache.NewMemory())
}

func TestImportDoesNotPanicOnNilTags(t *testing.T) {
	s := newTestService()
	data := []byte(`{"workspace":{"id":"workspace-import","name":"workspace-import"},"channels":[],"bundles":[{"id":"bundle-import","workspace_id":"workspace-import","channel_id":"channel-import","key":"risk-score","type":"string","default_value":"default","rules":[{"id":"rule-import","priority":1,"value":"blocked"}]}]}`)
	if err := s.ImportModerationWorkspace(context.Background(), data); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeRuleBundleDoesNotPanicOnNilTags(t *testing.T) {
	bundle, err := DecodeRuleBundle([]byte(`{"id":"bundle-decode","workspace_id":"workspace-decode","channel_id":"channel-decode","key":"risk-score","type":"string","default_value":"default","rules":[{"id":"rule-decode","priority":1,"value":"blocked"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Rules[0].Tags["source"] != "import" {
		t.Fatalf("tags not normalized: %v", bundle.Rules[0].Tags)
	}
}

func TestExportMissingReturnsError(t *testing.T) {
	s := newTestService()
	if _, err := s.ExportModerationWorkspace(context.Background(), "missing-workspace"); err == nil {
		t.Fatal("expected missing workspace to return an error")
	}
}
