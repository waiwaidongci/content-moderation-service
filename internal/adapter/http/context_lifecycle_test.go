package httpadapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ali/go-0821/content-moderation-service/internal/adapter/cache"
	"github.com/ali/go-0821/content-moderation-service/internal/application"
	"github.com/ali/go-0821/content-moderation-service/internal/infrastructure/logging"
	"github.com/ali/go-0821/content-moderation-service/internal/infrastructure/metrics"
	"github.com/ali/go-0821/content-moderation-service/internal/infrastructure/repository"
)

func newTestHandler() http.Handler {
	bundles := application.NewService(repository.NewMemory(), cache.NewMemory())
	return New(bundles, logging.New(), metrics.New()).Handler()
}

func canceledRequest(method, path, body string) *http.Request {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	return req.WithContext(ctx)
}

func TestCanceledWorkspaceCreateDoesNotProceed(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestHandler().ServeHTTP(rec, canceledRequest(http.MethodPost, "/v1/moderation/workspaces", `{"id":"workspace-cancel","name":"workspace-cancel"}`))
	if rec.Code == http.StatusCreated {
		t.Fatalf("canceled workspace create should not return 201, got %d", rec.Code)
	}
}

func TestCanceledChannelCreateDoesNotProceed(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestHandler().ServeHTTP(rec, canceledRequest(http.MethodPost, "/v1/moderation/channels", `{"id":"channel-cancel","workspace_id":"workspace-cancel","name":"channel-cancel"}`))
	if rec.Code == http.StatusCreated {
		t.Fatalf("canceled channel create should not return 201, got %d", rec.Code)
	}
}

func TestCanceledBundleCreateDoesNotProceed(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestHandler().ServeHTTP(rec, canceledRequest(http.MethodPost, "/v1/moderation/bundles", `{"id":"bundle-cancel","workspace_id":"workspace-cancel","channel_id":"channel-cancel","key":"risk-score","type":"string","default_value":"default"}`))
	if rec.Code == http.StatusCreated {
		t.Fatalf("canceled bundle create should not return 201, got %d", rec.Code)
	}
}
