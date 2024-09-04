package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type DB interface {
	// Select queries multiple rows and will load the entire result all at once
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Get queries for a single row and will load  the engire result all at once
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Select queries multiple rows and will load the entire result all at once
	Select_RO(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Get queries for a single row and will load  the engire result all at once
	Get_RO(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// Exec executes a query that changes rows
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// BeginTx will create a new transaction using the writer
	BeginTx(ctx context.Context) (Tx, error)

	// Handles transactions start to finish
	Transaction(ctx context.Context, fn func(ctx context.Context, tx Tx) error) error
}

type baseDB struct {
	WriteDB *sqlx.DB
	ReadDB  *sqlx.DB
}

func newBaseDB(driver, writerDSN, readerDSN string) (*baseDB, error) {
	writer, err := sqlx.Open(driver, writerDSN)
	if err != nil {
		log.Err(err).Msg("failed to connect to writer")
		return nil, err
	}

	err = writer.Ping()
	if err != nil {
		log.Err(err).Msg("failed to ping writer")
		return nil, err
	}

	var reader *sqlx.DB
	if readerDSN != "" {
		reader, err = sqlx.Open(driver, readerDSN)
		if err != nil {
			log.Err(err).Msg("failed to connect to reader")
			return nil, err
		}
		err = reader.Ping()
		if err != nil {
			log.Err(err).Msg("failed to ping reader")
			return nil, err
		}
	}

	// by default, always have the writer available
	db := &baseDB{
		WriteDB: writer,
	}

	log.Info().Msg("writer connected and initialized")

	if reader != nil {
		db.ReadDB = reader
		log.Info().Msg("reader connected and initialized")
	} else {
		db.ReadDB = db.WriteDB // Fallback to writer if no reader is provided
		log.Info().Msg("no reader available, falling back to writer")
	}

	return db, nil
}
