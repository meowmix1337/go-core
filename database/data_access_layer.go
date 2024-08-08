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

	// CommitTx commits the current transaction
	CommitTx(tx *sqlx.Tx) error

	// RollbackTx rolls back the current transaction
	RollbackTx(tx *sqlx.Tx) error

	// Close closes the database connection
	Close() error
}

// DBWrapper wraps the sqlx functions
// this will help prevent me from using functions I don't care about or need
// I can always expose more functions as needed in the future
type DBWrapper struct {
	db *sqlx.DB
}

func NewDBWrapper(db *sqlx.DB) *DBWrapper {
	return &DBWrapper{
		db: db,
	}
}

// QueryRow executes a query that reeturns a single row
func (m *DBWrapper) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return m.db.QueryRowxContext(ctx, query, args...)
}

// QueryRows exeutes a query that returns multiple rows
func (m *DBWrapper) QueryRows(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return m.db.QueryxContext(ctx, query, args...)
}

// Select queries multiple rows and will load the entire result all at once
func (m *DBWrapper) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.db.SelectContext(ctx, dest, query, args...)
}

// Get queries for a single row and will load the engire result all at once
func (m *DBWrapper) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.db.GetContext(ctx, dest, query, args...)
}

// Exec executes a query that changes rows
func (m *DBWrapper) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

// ExecTx executes a query that changes rows within a given transaction
func (m *DBWrapper) ExecTx(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) (sql.Result, error) {
	return tx.ExecContext(ctx, query, args...)
}

// Commit commits the current transaction
func (m *DBWrapper) CommitTx(tx *sqlx.Tx) error {
	return tx.Commit()
}

// Rollback rolls back the current transaction
func (m *DBWrapper) RollbackTx(tx *sqlx.Tx) error {
	return tx.Rollback()
}

// Close closes the database connection
func (m *DBWrapper) Close() error {
	return m.db.Close()
}
