package database

import (
	"context"
	"fmt"
)

// RunInTransaction runs a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (m *Manager) RunInTransaction(ctx context.Context, fn func(tx Queryer) error) error {
	tx, err := m.writer.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-panic after rollback
		} else if err != nil {
			tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error, update err
		}
	}()

	err = fn(tx)
	return err
}
