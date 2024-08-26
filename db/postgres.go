package db

import (
	"context"
	"database/sql"

	"github.com/meowmix1337/go-core/derror"
	"github.com/rs/zerolog/log"
)

type postgres struct {
	*baseDB
}

func NewPostgres(writerDSN, readerDSN string) *postgres {
	baseDB, err := newBaseDB("postgres", writerDSN, readerDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to DB")
	}

	return &postgres{
		baseDB: baseDB,
	}
}

// Select queries multiple rows and will load the entire result all at once
func (m *postgres) Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.WriteDB.SelectContext(ctx, dest, query, args...)
}

// Get queries for a single row and will load the engire result all at once
func (m *postgres) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.WriteDB.GetContext(ctx, dest, query, args...)
}

// Select queries multiple rows and will load the entire result all at once
func (m *postgres) Select_RO(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.ReadDB.SelectContext(ctx, dest, query, args...)
}

// Get queries for a single row and will load the engire result all at once
func (m *postgres) Get_RO(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return m.ReadDB.GetContext(ctx, dest, query, args...)
}

// Exec executes a query that changes rows
func (m *postgres) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.WriteDB.ExecContext(ctx, query, args...)
}

// BeginTx will create a new transaction using the writer
func (m *postgres) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := m.WriteDB.BeginTxx(ctx, nil)
	if err != nil {
		log.Err(err).Msg("Failed to start transaction")
		return nil, err
	}
	return NewMySQLTx(tx), nil
}

func (m *postgres) Transaction(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error {

	tx, err := m.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		// handle panic
		if p := recover(); p != nil {
			log.Err(err).Msg("panic happened, rolling back")
			_ = tx.Rollback()
			panic(p) // Re-throw panic after rollback
		}
	}()

	err = fn(ctx, tx)
	if err != nil {
		log.Err(err).Msg("query failed, attempt rolling back")
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Err(rbErr).Msg("failed to roll back transaction")
			return derror.New(ctx, derror.InternalServerCode, derror.InternalType, "error rolling back", rbErr).Wrap(err)
		}
		log.Info().Msg("roll back successful")
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		log.Err(commitErr).Msg("failed to commit transaction")
		return commitErr
	}

	return nil
}
