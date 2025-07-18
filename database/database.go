package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Manager is a database manager that provides reader and writer capabilities.
// It uses sqlx for database operations.
type Manager struct {
	driverName      string
	reader          *sqlx.DB
	writer          *sqlx.DB
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
	queryTimeout    time.Duration
}

func buildDSN(cfg Config) string {
	// [user[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	if cfg.InsecureSkipVerify {
		dsn += "&tls=skip-verify"
	} else if cfg.SSLMode != "" && cfg.SSLMode != "disable" {
		dsn += "&tls=true"
	}

	return dsn
}

// NewManager creates a new database manager.
func NewManager(driverName string, writerCfg Config, readerCfg *Config, opts ...Option) (*Manager, error) {
	m := &Manager{
		driverName:      driverName,
		maxOpenConns:    defaultMaxOpenConns,
		maxIdleConns:    defaultMaxIdleConns,
		connMaxLifetime: defaultConnMaxLifetime,
		connMaxIdleTime: defaultConnMaxIdleTime,
		queryTimeout:    defaultQueryTimeout,
	}

	for _, opt := range opts {
		opt(m)
	}

	writerDSN := buildDSN(writerCfg)
	writerDB, err := sqlx.Connect(driverName, writerDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to writer database: %w", err)
	}
	m.writer = writerDB
	m.writer.SetMaxOpenConns(m.maxOpenConns)
	m.writer.SetMaxIdleConns(m.maxIdleConns)
	m.writer.SetConnMaxLifetime(m.connMaxLifetime)
	m.writer.SetConnMaxIdleTime(m.connMaxIdleTime)

	if readerCfg == nil {
		m.reader = m.writer
	} else {
		readerDSN := buildDSN(*readerCfg)
		readerDB, err := sqlx.Connect(driverName, readerDSN)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to reader database: %w", err)
		}
		m.reader = readerDB
		m.reader.SetMaxOpenConns(m.maxOpenConns)
		m.reader.SetMaxIdleConns(m.maxIdleConns)
		m.reader.SetConnMaxLifetime(m.connMaxLifetime)
		m.reader.SetConnMaxIdleTime(m.connMaxIdleTime)
	}

	return m, nil
}

// Reader returns the reader database connection.
func (m *Manager) Reader() *sqlx.DB {
	return m.reader
}

// Writer returns the writer database connection.
func (m *Manager) Writer() *sqlx.DB {
	return m.writer
}

// Close closes the database connections.
func (m *Manager) Close() error {
	if err := m.writer.Close(); err != nil {
		return fmt.Errorf("failed to close writer database: %w", err)
	}
	if err := m.reader.Close(); err != nil {
		return fmt.Errorf("failed to close reader database: %w", err)
	}
	return nil
}

// Queryer is an interface for query operations.
type Queryer interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}
