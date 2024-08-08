package database

import (
	"github.com/jmoiron/sqlx"
)

// DBConnector is an interface that will allow users to use a writer or reader
type DBConnector interface {
	WriteDB() *DBWrapper
	ReadDB() *DBWrapper
}

// Database represents the writer and reader DBs
type Database struct {
	writerDB *DBWrapper
	readerDB *DBWrapper
}

func NewDBConnector(driver, writerDSN, readerDSN string) (*Database, error) {
	writer, err := sqlx.Open(driver, writerDSN)
	if err != nil {
		return nil, err
	}

	var reader *sqlx.DB
	if readerDSN != "" {
		reader, err = sqlx.Open(driver, readerDSN)
		if err != nil {
			return nil, err
		}
	}

	// by default, always have the writer available
	db := &Database{
		writerDB: NewDBWrapper(writer),
	}

	if reader != nil {
		db.readerDB = &DBWrapper{db: reader}
	} else {
		db.readerDB = db.writerDB // Fallback to writer if no reader is provided
	}

	return db, nil
}

func (c *Database) WriteDB() *DBWrapper {
	return c.writerDB
}

func (c *Database) ReadDB() *DBWrapper {
	return c.readerDB
}
