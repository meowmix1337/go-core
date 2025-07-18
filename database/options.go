package database

import "time"

const (
	defaultMaxOpenConns    = 10
	defaultMaxIdleConns    = 5
	defaultConnMaxLifetime = time.Hour
	defaultConnMaxIdleTime = 5 * time.Minute
	defaultQueryTimeout    = 5 * time.Second
)

// Option is a function that sets an option on a Manager.
type Option func(*Manager)

// Config holds the database connection DSNs.
type Config struct {
	Host               string
	Port               int
	User               string
	Password           string
	DBName             string
	SSLMode            string
	// TLS configuration options
	TLSCertPath        string
	TLSKeyPath         string
	TLSCaCertPath      string
	InsecureSkipVerify bool
}

// WithMaxOpenConns sets the maximum number of open connections to the database.
func WithMaxOpenConns(n int) Option {
	return func(m *Manager) {
		m.maxOpenConns = n
	}
}

// WithMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
func WithMaxIdleConns(n int) Option {
	return func(m *Manager) {
		m.maxIdleConns = n
	}
}

// WithConnMaxLifetime sets the maximum amount of time a connection may be reused.
func WithConnMaxLifetime(d time.Duration) Option {
	return func(m *Manager) {
		m.connMaxLifetime = d
	}
}

// WithConnMaxIdleTime sets the maximum amount of time a connection may be idle.
func WithConnMaxIdleTime(d time.Duration) Option {
	return func(m *Manager) {
		m.connMaxIdleTime = d
	}
}

// WithQueryTimeout sets the timeout for queries.
func WithQueryTimeout(d time.Duration) Option {
	return func(m *Manager) {
		m.queryTimeout = d
	}
}
