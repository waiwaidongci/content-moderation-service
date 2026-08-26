// Package implementation for policy-driven content moderation and human review.
// Package httpadapter exposes policy-driven moderation and review workflows.
package httpadapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ali/go-0821/content-moderation-service/internal/application"
	"github.com/ali/go-0821/content-moderation-service/internal/domain"
	"github.com/ali/go-0821/content-moderation-service/internal/infrastructure/logging"
	"github.com/ali/go-0821/content-moderation-service/internal/infrastructure/metrics"
	"io"
	"net/http"
	"time"
)

type Server struct {
	bundles    *application.Service
	moderation *application.ModerationService
	log        *logging.Logger
	metrics    *metrics.Metrics
	started    time.Time
}

func New(bundles *application.Service, l *logging.Logger, m *metrics.Metrics) *Server {
	return &Server{bundles: bundles, moderation: application.NewModerationService(), log: l, metrics: m, started: time.Now()}
}
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.health)
	mux.HandleFunc("/readyz", s.ready)
	mux.HandleFunc("/metrics", s.metrics.Handler)
	mux.HandleFunc("/v1/moderation/policies", s.moderationPolicies)
	mux.HandleFunc("/v1/moderation/policies/publish", s.publishModerationPolicy)
	mux.HandleFunc("/v1/moderation/check", s.moderateContent)
	mux.HandleFunc("/v1/moderation/tasks/claim", s.claimReviewTask)
	mux.HandleFunc("/v1/moderation/tasks/submit", s.submitReviewTask)
	mux.HandleFunc("/v1/moderation/decisions/", s.moderationDecision)
	mux.HandleFunc("/v1/moderation/appeals", s.moderationAppeals)
	mux.HandleFunc("/v1/moderation/dictionaries", s.moderationDictionaries)
	mux.HandleFunc("/v1/moderation/workspaces", s.moderationWorkspaces)
	mux.HandleFunc("/v1/moderation/channels", s.moderationChannels)
	mux.HandleFunc("/v1/moderation/bundles", s.ruleBundles)
	return s.middleware(mux)
}
func (s *Server) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.metrics.Inc()
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = fmt.Sprintf("req-%d", time.Now().UnixNano())
		}
		w.Header().Set("X-Request-ID", id)
		start := time.Now()
		ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancel()
		next.ServeHTTP(w, r.WithContext(ctx))
		s.log.Info("http_request", map[string]any{"method": r.Method, "path": r.URL.Path, "duration_ms": time.Since(start).Milliseconds(), "request_id": id})
	})
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]any{"status": "ok", "uptime": time.Since(s.started).String()})
}
func (s *Server) ready(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ready"})
}
func decode(r *http.Request, v any) error {
	b, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}
	if len(b) == 0 {
		return fmt.Errorf("%w: empty body", domain.ErrInvalid)
	}
	if err = json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("%w: invalid JSON", domain.ErrInvalid)
	}
	return nil
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, err error) {
	status := 500
	if errors.Is(err, domain.ErrNotFound) {
		status = 404
	} else if errors.Is(err, domain.ErrConflict) {
		status = 409
	} else if errors.Is(err, domain.ErrInvalid) {
		status = 400
	}
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
func Serve(ctx context.Context, addr string, h http.Handler) error {
	srv := &http.Server{Addr: addr, Handler: h, ReadHeaderTimeout: 5 * time.Second}
	errs := make(chan error, 1)
	go func() { errs <- srv.ListenAndServe() }()
	select {
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdown)
	case err := <-errs:
		if err == http.ErrServerClosed {
			return nil
		}
		return fmt.Errorf("serve HTTP: %w", err)
	}
}
