package httpadapter

import (
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

func TestDecisionMissingReturns404(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/moderation/decisions/missing-sample", nil)
	rec := httptest.NewRecorder()
	newTestHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestClaimReturns404(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/moderation/tasks/claim", strings.NewReader(`{"reviewer":"alice"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	newTestHandler().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
