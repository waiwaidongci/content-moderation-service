package repository

import (
	"context"
	"testing"

	"github.com/ali/go-0821/content-moderation-service/internal/domain"
)

func TestRollbackRunsUndoInReverseOrder(t *testing.T) {
	tx := NewTransaction()
	var order []string
	tx.AddUndo(func() { order = append(order, "first") })
	tx.AddUndo(func() { order = append(order, "second") })
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	if len(order) != 2 || order[0] != "second" || order[1] != "first" {
		t.Fatalf("wrong undo order: %v", order)
	}
}

func TestWithTransactionDoesNotCommitOnError(t *testing.T) {
	var rolledBack bool
	err := WithTransaction(context.Background(), func(ctx context.Context, tx *Transaction) error {
		tx.AddUndo(func() { rolledBack = true })
		return context.DeadlineExceeded
	})
	if err == nil {
		t.Fatal("expected transaction error to be returned")
	}
	if !rolledBack {
		t.Fatal("expected rollback undo to run")
	}
}

func TestValidateRuleBundleRejectsInvalid(t *testing.T) {
	bundle := domain.RuleBundle{ID: "", Type: domain.TypeString, DefaultValue: "default"}
	if err := validateRuleBundle(bundle); err == nil {
		t.Fatal("expected invalid rule bundle to be rejected")
	}
}

func TestValidateBundleRevisionRejectsInvalid(t *testing.T) {
	revision := domain.BundleRevision{Number: 0, Value: "value"}
	if err := validateBundleRevision(revision, domain.TypeString); err == nil {
		t.Fatal("expected invalid bundle revision to be rejected")
	}
}

func TestAddUndoAfterCommitIgnored(t *testing.T) {
	tx := NewTransaction()
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	tx.AddUndo(func() {})
	if len(tx.undo) != 0 {
		t.Fatal("expected undo added after commit to be ignored")
	}
}
