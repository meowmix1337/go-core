package database

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type MySQLClient struct {
	db *sqlx.DB
}

func NewMySQLClient(dsn string) (*MySQLClient, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Err(err).Msg("failed to connect to mysql database")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Err(err).Msg("ping to mysql database failed")
		return nil, err
	}

	return &MySQLClient{
		db: db,
	}, nil
}

// QueryRow executes a query that reeturns a single row
func (m *MySQLClient) QueryRow(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	return m.db.QueryRowxContext(ctx, query, args...)
}

// QueryRows exeutes a query that returns multiple rows
func (m *MySQLClient) QueryRows(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	return m.db.QueryxContext(ctx, query, args...)
}

// Select queries multiple rows and will load the entire result all at once
func (m *MySQLClient) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.db.SelectContext(ctx, dest, query, args...)
}

// Get queries for a single row and will load the engire result all at once
func (m *MySQLClient) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.db.GetContext(ctx, dest, query, args...)
}

// Exec executes a query that changes rows
func (m *MySQLClient) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

// ExecTx executes a query that changes rows within a given transaction
func (m *MySQLClient) ExecTx(ctx context.Context, tx *sqlx.Tx, query string, args ...interface{}) (sql.Result, error) {
	return tx.ExecContext(ctx, query, args...)
}

// Begin starts a new transaction
func (m *MySQLClient) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return m.db.BeginTxx(ctx, nil)
}

// Commit commits the current transaction
func (m *MySQLClient) CommitTx(tx *sqlx.Tx) error {
	return tx.Commit()
}

// Rollback rolls back the current transaction
func (m *MySQLClient) RollbackTx(tx *sqlx.Tx) error {
	return tx.Rollback()
}

// Close closes the database connection
func (m *MySQLClient) Close() error {
	return m.db.Close()
}
