package cache

import (
	"context"
	"testing"
	"time"

	"github.com/ali/go-0821/content-moderation-service/internal/domain"
)

func TestTTLSetSkipsRevokedRevision(t *testing.T) {
	c := NewTTL(time.Minute)
	revoked := domain.BundleRevision{Number: 2, RuleBundleID: "bundle-status", Status: "revoked"}
	c.Set(context.Background(), "bundle-status", revoked)
	if c.Size() != 0 {
		t.Fatal("expected revoked revision not to be stored")
	}
}

func TestTTLSetSkipsDraftRevision(t *testing.T) {
	c := NewTTL(time.Minute)
	draft := domain.BundleRevision{Number: 2, RuleBundleID: "bundle-draft", Status: "draft"}
	c.Set(context.Background(), "bundle-draft", draft)
	if c.Size() != 0 {
		t.Fatal("expected draft revision not to be stored")
	}
}

func TestTTLGetSkipsNoneffectiveRevision(t *testing.T) {
	c := NewTTL(time.Minute)
	revoked := domain.BundleRevision{Number: 2, RuleBundleID: "bundle-get", Status: "revoked"}
	c.data["bundle-get"] = entry{revision: revoked, expires: time.Now().Add(time.Minute)}
	if _, ok := c.Get(context.Background(), "bundle-get"); ok {
		t.Fatal("expected non-effective revision not to be returned")
	}
}
