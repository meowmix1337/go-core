package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DatabaseAccessLayer represents a database access layer interface
type DatabaseAccessLayer interface {
	// QueryRow executes a query that returns a single row
	QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row

	// QueryRows executes a query that returns multiple rows
	QueryRows(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)

	// Select queries multiple rows and will load the entire result all at once
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Get queries for a single row and will load  the engire result all at once
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Exec executes a query that changes rows
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// ExecTx executes a query that changes rows within a given transaction
	ExecTx(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) (sql.Result, error)

	// BeginTx starts a new transaction
	BeginTx(ctx context.Context) (*sqlx.Tx, error)

	// CommitTx commits the current transaction
	CommitTx(tx *sqlx.Tx) error

	// RollbackTx rolls back the current transaction
	RollbackTx(tx *sqlx.Tx) error

	// Close closes the database connection
	Close() error
}
