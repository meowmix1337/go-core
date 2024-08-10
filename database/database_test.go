package database

import (
	"context"
	"errors"
	"log"
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
	mySQLClient *Database
	mock        sqlmock.Sqlmock
}

func TestDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(DatabaseTestSuite))
}

func (s *DatabaseTestSuite) SetupSuite() {
	// Set up the test
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	sqlxDB := sqlx.NewDb(db, "mock")
	s.mySQLClient = &Database{
		writerDB: &DBWrapper{db: sqlxDB},
		writeDB:  sqlxDB,
	}
	s.mock = mock
}

func (s *DatabaseTestSuite) TestDatabase_QueryRow() {

	// Define the expected query and result
	s.mock.ExpectQuery(`SELECT \* FROM users WHERE id = \?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "John Doe", "john.doe@example.com"))

	// Call the QueryRow method
	row := s.mySQLClient.WriteDB().QueryRow(context.Background(), "SELECT * FROM users WHERE id = ?", 1)

	// Verify that the result is correct
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		log.Fatal(err)
	}
	if user.ID != 1 || user.Name != "John Doe" || user.Email != "john.doe@example.com" {
		s.T().Errorf("Expected user %+v, got %+v", User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}, user)
	}
}

func (s *DatabaseTestSuite) TestDatabase_QueryRows() {

	// Define the expected query and result
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users")).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "John Doe", "john.doe@example.com").
			AddRow(2, "Jane Doe", "jane.doe@example.com"))

	// Call the QueryRow method
	rows, err := s.mySQLClient.WriteDB().QueryRows(context.Background(), "SELECT * FROM users")
	s.NoError(err)

	// Verify that the result is correct
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		s.NoError(err)
		users = append(users, user)
	}

	if len(users) != 2 {
		s.T().Errorf("expected 2 users but got %v", len(users))
	}
}

func (s *DatabaseTestSuite) TestDatabase_Get() {

	// Define the expected query and result
	s.mock.ExpectQuery(`SELECT \* FROM users WHERE id = \?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "John Doe", "john.doe@example.com"))

	// Call the QueryRow method
	var user User
	err := s.mySQLClient.WriteDB().Get(context.Background(), &user, "SELECT * FROM users WHERE id = ?", 1)
	if err != nil {
		log.Fatal(err)
	}

	if user.ID != 1 || user.Name != "John Doe" || user.Email != "john.doe@example.com" {
		s.T().Errorf("Expected user %+v, got %+v", User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}, user)
	}
}

func (s *DatabaseTestSuite) TestDatabase_Select() {

	// Define the expected query and result
	s.mock.ExpectQuery(`SELECT \* FROM users`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "John Doe", "john.doe@example.com").
			AddRow(2, "Jane Doe", "jane.doe@example.com"))

	// Call the QueryRow method
	var users []User
	err := s.mySQLClient.WriteDB().Select(context.Background(), &users, "SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
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
	s.mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (name, email) VALUES (?, ?)")).
		WithArgs(userToInsert.Name, userToInsert.Email).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Assuming ID=1 and 1 row affected

	args := []interface{}{userToInsert.Name, userToInsert.Email}
	result, err := s.mySQLClient.WriteDB().Exec(context.Background(), "INSERT INTO users (name, email) VALUES (?, ?)", args...)
	s.NoError(err)

	id, err := result.LastInsertId()
	s.NoError(err)

	if id != 1 {
		s.T().Errorf("expected id to be 1")
	}

	// Ensure all expectations are met
	err = s.mock.ExpectationsWereMet()
	s.NoError(err)
}

func (s *DatabaseTestSuite) TestTransaction_Success() {

	s.mock.ExpectBegin()
	s.mock.ExpectCommit()

	err := Transaction(context.Background(), s.mySQLClient, func(tx *sqlx.Tx) error {
		return nil // Simulate a successful function
	})
	s.NoError(err)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DatabaseTestSuite) TestTransaction_Failure() {

	s.mock.ExpectBegin()
	s.mock.ExpectRollback()

	err := Transaction(context.Background(), s.mySQLClient, func(tx *sqlx.Tx) error {
		return errors.New("failed") // Simulate a failure function
	})
	s.Errorf(err, "failed")
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DatabaseTestSuite) TestTransaction_Panic() {

	s.mock.ExpectBegin()
	s.mock.ExpectRollback()

	defer func() {
		if r := recover(); r != nil {
			s.Equal("unexpected panic", r)
			s.NoError(s.mock.ExpectationsWereMet())
		} else {
			s.T().Errorf("expected panic, but code did not panic")
		}
	}()

	err := Transaction(context.Background(), s.mySQLClient, func(tx *sqlx.Tx) error {
		panic("unexpected panic") // Simulate a panic
	})
	s.NoError(err)
}

func (s *DatabaseTestSuite) TestTransaction_BeginTxFail() {

	s.mock.ExpectBegin().WillReturnError(errors.New("begin tx failed"))

	err := Transaction(context.Background(), s.mySQLClient, func(tx *sqlx.Tx) error {
		return nil // This should not be executed
	})
	s.Error(err)
	s.Equal("begin tx failed", err.Error())
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *DatabaseTestSuite) TestTransaction_CommitFail() {

	commitErr := errors.New("commit failed")
	s.mock.ExpectBegin()
	s.mock.ExpectCommit().WillReturnError(commitErr)

	err := Transaction(context.Background(), s.mySQLClient, func(tx *sqlx.Tx) error {
		return nil // Simulate a successful function
	})
	s.Error(err)
	s.Equal("commit failed", err.Error())
	s.NoError(s.mock.ExpectationsWereMet())
}
