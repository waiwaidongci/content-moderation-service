package application

import (
	"context"
	"errors"
	"testing"

	"github.com/ali/go-0821/content-moderation-service/internal/domain"
)

func TestDecisionMissingWrapsSentinel(t *testing.T) {
	s := NewModerationService()
	_, err := s.Decision("missing-sample")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound chain, got %v", err)
	}
}

func TestClaimMissingWrapsSentinel(t *testing.T) {
	s := NewModerationService()
	_, err := s.Claim(context.Background(), "reviewer-1")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound chain, got %v", err)
	}
}
