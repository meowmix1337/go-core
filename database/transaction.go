package database

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

// WithTransaction runs a function within a database transaction.
// If the function returns an error, the transaction is rolled back.
// Otherwise, the transaction is committed.
func (m *Manager) WithTransaction(ctx context.Context, fn func(tx Queryer) error) (err error) {
	tx, err := m.writer.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				zerolog.Ctx(ctx).Error().Err(rbErr).Msg("rollback failed after panic")
			}
			panic(p) // re-panic after rollback
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("transaction failed: %v, and rollback failed: %w", err, rbErr)
			} else {
				err = fmt.Errorf("transaction failed: %w", err)
			}
		} else {
			if cErr := tx.Commit(); cErr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", cErr)
			}
		}
	}()

	err = fn(tx)
	return err
}
