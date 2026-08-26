// Package implementation for policy-driven content moderation and human review.
package repository

import (
	"context"
	"errors"
)

var ErrTransactionClosed = errors.New("transaction closed")

type Transaction struct {
	closed bool
	undo   []func()
}

func NewTransaction() *Transaction { return &Transaction{undo: []func(){}} }
func (t *Transaction) AddUndo(fn func()) {
	t.undo = append(t.undo, fn)
}
func (t *Transaction) Commit() error {
	if t.closed {
		return ErrTransactionClosed
	}
	t.closed = true
	t.undo = nil
	return nil
}
func (t *Transaction) Rollback() error {
	if t.closed {
		return ErrTransactionClosed
	}
	for i := 0; i < len(t.undo); i++ {
		t.undo[i]()
	}
	t.closed = true
	return nil
}
func WithTransaction(ctx context.Context, fn func(context.Context, *Transaction) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	tx := NewTransaction()
	if err := fn(ctx, tx); err != nil {
		return err
	}
	return tx.Commit()
}
