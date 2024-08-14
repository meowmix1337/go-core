package db

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

type DatabaseTestSuite struct {
	suite.Suite
	mysqlDB    *mySQL
	mockWriter sqlmock.Sqlmock
	mockReader sqlmock.Sqlmock
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (s *DatabaseTestSuite) SetupSuite() {
	// Set up the test
	writeDB, mockWriter, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	writer := sqlx.NewDb(writeDB, "mock")

	readDB, mockReader, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	reader := sqlx.NewDb(readDB, "mock")

	s.mysqlDB = &mySQL{
		baseDB: &baseDB{WriteDB: writer, ReadDB: reader},
	}
	s.mockWriter = mockWriter
	s.mockReader = mockReader
}

func (s *DatabaseTestSuite) TestDatabase_Get() {

	// Define the expected query and result
	s.mockWriter.ExpectQuery(`SELECT \* FROM users WHERE id = \?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "John Doe", "john.doe@example.com"))

	// Call the QueryRow method
	var user User
	err := s.mysqlDB.Get(context.Background(), &user, "SELECT * FROM users WHERE id = ?", 1)
	if err != nil {
		s.T().Error(err)
	}

	if user.ID != 1 || user.Name != "John Doe" || user.Email != "john.doe@example.com" {
		s.T().Errorf("Expected user %+v, got %+v", User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}, user)
	}
}

func (s *DatabaseTestSuite) TestDatabase_Get_RO() {

	// Define the expected query and result
	s.mockReader.ExpectQuery(`SELECT \* FROM users WHERE id = \?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "John Doe", "john.doe@example.com"))

	// Call the QueryRow method
	var user User
	err := s.mysqlDB.Get_RO(context.Background(), &user, "SELECT * FROM users WHERE id = ?", 1)
	if err != nil {
		s.T().Error(err)
	}

	if user.ID != 1 || user.Name != "John Doe" || user.Email != "john.doe@example.com" {
		s.T().Errorf("Expected user %+v, got %+v", User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}, user)
	}
}

func (s *DatabaseTestSuite) TestDatabase_Select() {

	// Define the expected query and result
	s.mockWriter.ExpectQuery(`SELECT \* FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "John Doe", "john.doe@example.com").
			AddRow(2, "Jane Doe", "jane.doe@example.com"))

	// Call the QueryRow method
	var users []User
	err := s.mysqlDB.Select(context.Background(), &users, "SELECT * FROM users")
	if err != nil {
		s.T().Error(err)
	}

	if len(users) != 2 {
		s.T().Errorf("Expected 2 users, got %+v", len(users))
	}
}

func (s *DatabaseTestSuite) TestDatabase_Select_RO() {

	// Define the expected query and result
	s.mockReader.ExpectQuery(`SELECT \* FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "John Doe", "john.doe@example.com").
			AddRow(2, "Jane Doe", "jane.doe@example.com"))

	// Call the QueryRow method
	var users []User
	err := s.mysqlDB.Select_RO(context.Background(), &users, "SELECT * FROM users")
	if err != nil {
		s.T().Error(err)
	}

	if len(users) != 2 {
		s.T().Errorf("Expected 2 users, got %+v", len(users))
	}
}

func (s *DatabaseTestSuite) TestDatabase_Exec() {

	userToInsert := User{
		Email: "john.doe@example.com",
		Name:  "John Doe",
	}

	// Expectation: The INSERT statement will be called with specific arguments
	s.mockWriter.ExpectExec(regexp.QuoteMeta("INSERT INTO users (name, email) VALUES (?, ?)")).
		WithArgs(userToInsert.Name, userToInsert.Email).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Assuming ID=1 and 1 row affected

	args := []interface{}{userToInsert.Name, userToInsert.Email}
	result, err := s.mysqlDB.Exec(context.Background(), "INSERT INTO users (name, email) VALUES (?, ?)", args...)
	s.NoError(err)

	id, err := result.LastInsertId()
	s.NoError(err)

	if id != 1 {
		s.T().Errorf("expected id to be 1")
	}

	// Ensure all expectations are met
	err = s.mockWriter.ExpectationsWereMet()
	s.NoError(err)
}

func (s *DatabaseTestSuite) TestTransaction_Success() {

	s.mockWriter.ExpectBegin()
	s.mockWriter.ExpectCommit()

	err := s.mysqlDB.Transaction(context.Background(), func(ctx context.Context, tx Tx) error {
		return nil // Simulate a successful function
	})
	s.NoError(err)
	s.NoError(s.mockWriter.ExpectationsWereMet())
}

func (s *DatabaseTestSuite) TestTransaction_Failure() {

	s.mockWriter.ExpectBegin()
	s.mockWriter.ExpectRollback()

	err := s.mysqlDB.Transaction(context.Background(), func(ctx context.Context, tx Tx) error {
		return errors.New("failed") // Simulate a failure function
	})
	s.Errorf(err, "failed")
	s.NoError(s.mockWriter.ExpectationsWereMet())
}

func (s *DatabaseTestSuite) TestTransaction_Panic() {

	s.mockWriter.ExpectBegin()
	s.mockWriter.ExpectRollback()

	defer func() {
		if r := recover(); r != nil {
			s.Equal("unexpected panic", r)
			s.NoError(s.mockWriter.ExpectationsWereMet())
		} else {
			s.T().Errorf("expected panic, but code did not panic")
		}
	}()

	err := s.mysqlDB.Transaction(context.Background(), func(ctx context.Context, tx Tx) error {
		panic("unexpected panic") // Simulate a panic
	})
	s.NoError(err)
}

func (s *DatabaseTestSuite) TestTransaction_BeginTxFail() {

	s.mockWriter.ExpectBegin().WillReturnError(errors.New("begin tx failed"))

	err := s.mysqlDB.Transaction(context.Background(), func(ctx context.Context, tx Tx) error {
		return nil // This should not be executed
	})
	s.Error(err)
	s.Equal("begin tx failed", err.Error())
	s.NoError(s.mockWriter.ExpectationsWereMet())
}

func (s *DatabaseTestSuite) TestTransaction_CommitFail() {

	commitErr := errors.New("commit failed")
	s.mockWriter.ExpectBegin()
	s.mockWriter.ExpectCommit().WillReturnError(commitErr)

	err := s.mysqlDB.Transaction(context.Background(), func(ctx context.Context, tx Tx) error {
		return nil // Simulate a successful function
	})
	s.Error(err)
	s.Equal("commit failed", err.Error())
	s.NoError(s.mockWriter.ExpectationsWereMet())
}
