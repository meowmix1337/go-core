package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type mySQLTx struct {
	tx *sqlx.Tx
}

func NewMySQLTx(tx *sqlx.Tx) Tx {
	return &mySQLTx{tx: tx}
}

// Select queries multiple rows and will load the entire result all at once
func (m *mySQLTx) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.tx.SelectContext(ctx, dest, query, args...)
}

// Get queries for a single row and will load the engire result all at once
func (m *mySQLTx) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.tx.GetContext(ctx, dest, query, args...)
}

// Exec executes a query that changes rows
func (m *mySQLTx) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.tx.ExecContext(ctx, query, args...)
}

// Commit commits the current transaction
func (m *mySQLTx) Commit() error {
	return m.tx.Commit()
}

// Rollback rolls back the current transaction
func (m *mySQLTx) Rollback() error {
	return m.tx.Rollback()
}
