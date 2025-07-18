package database

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestNewManager(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	// Test with a writer and no reader
	m, err := NewManagerWithDB("sqlmock", sqlxDB, nil)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	if m == nil {
		t.Fatal("manager is nil")
	}
	if m.Reader() == nil {
		t.Fatal("reader is nil")
	}
	if m.Writer() == nil {
		t.Fatal("writer is nil")
	}
	if m.Reader() != m.Writer() {
		t.Fatal("reader and writer should be the same")
	}
	m.Close()

	// Test with a writer and a reader
	m, err = NewManagerWithDB("sqlmock", sqlxDB, sqlxDB)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	if m == nil {
		t.Fatal("manager is nil")
	}
	if m.Reader() == nil {
		t.Fatal("reader is nil")
	}
	if m.Writer() == nil {
		t.Fatal("writer is nil")
	}
	m.Close()
}

func TestWithTransaction(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	sqlxDB := sqlx.NewDb(mockDB, "sqlmock")

	m := &Manager{
		writer: sqlxDB,
	}

	// Test a successful transaction
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO test").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = m.WithTransaction(context.Background(), func(tx Queryer) error {
		_, err := tx.ExecContext(context.Background(), "INSERT INTO test (name) VALUES ('test')")
		return err
	})
	if err != nil {
		t.Fatalf("failed to run transaction: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test a failed transaction
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO test").WillReturnError(fmt.Errorf("some error"))
	mock.ExpectRollback()

	err = m.WithTransaction(context.Background(), func(tx Queryer) error {
		_, err := tx.ExecContext(context.Background(), "INSERT INTO test (name) VALUES ('test')")
		return err
	})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// Test a panic
	mock.ExpectBegin()
	mock.ExpectRollback()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	err = m.WithTransaction(context.Background(), func(tx Queryer) error {
		panic("test panic")
	})
	if err != nil {
		t.Fatalf("expected panic, got error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
