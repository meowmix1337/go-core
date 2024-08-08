package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

// DBConnector is an interface that will allow users to use a writer or reader
type DBConnector interface {
	// WriteDB returns the DBWrapper that implements the DAL interface
	WriteDB() *DBWrapper

	// ReadDB returns the DBWrapper that implements the DAL interface
	ReadDB() *DBWrapper

	// BeginTx will create a new transaction using the writer
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
}

// Database represents the writer and reader DBs
type Database struct {
	writeDB  *sqlx.DB
	writerDB *DBWrapper
	readerDB *DBWrapper
}

func NewDBConnector(driver, writerDSN, readerDSN string) (*Database, error) {
	writer, err := sqlx.Open(driver, writerDSN)
	if err != nil {
		log.Err(err).Msg("failed to connect to writer")
		return nil, err
	}

	var reader *sqlx.DB
	if readerDSN != "" {
		reader, err = sqlx.Open(driver, readerDSN)
		if err != nil {
			log.Err(err).Msg("failed to connect to reader")
			return nil, err
		}
	}

	// by default, always have the writer available
	db := &Database{
		writerDB: NewDBWrapper(writer),
		writeDB:  writer,
	}

	log.Info().Msg("writer connected and initialized")

	if reader != nil {
		db.readerDB = &DBWrapper{db: reader}
		log.Info().Msg("reader connected and initialized")
	} else {
		db.readerDB = db.writerDB // Fallback to writer if no reader is provided
		log.Info().Msg("no reader available, falling back to writer")
	}

	return db, nil
}

// WriteDB returns the writer
func (d *Database) WriteDB() *DBWrapper {
	return d.writerDB
}

// ReadDB returns the reader, if no reader is available, this will point to the writer
func (d *Database) ReadDB() *DBWrapper {
	return d.readerDB
}

// Begin starts a new transaction
func (d *Database) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return d.writeDB.BeginTxx(ctx, nil)
}
