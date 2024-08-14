package db

import (
	"context"
	"database/sql"
)

type Tx interface {
	// Select queries multiple rows and will load the entire result all at once
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Get queries for a single row and will load  the engire result all at once
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Exec executes a query that changes rows
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// CommitTx commits the current transaction
	Commit() error

	// RollbackTx rolls back the current transaction
	Rollback() error
}
