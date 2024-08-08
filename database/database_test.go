package database

import (
	"context"
	"log"
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
	s.mySQLClient = &Database{writerDB: &DBWrapper{db: sqlxDB}}
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
