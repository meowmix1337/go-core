package database

import (
	"context"
	"fmt"
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

type MySQLClientTestSuite struct {
	suite.Suite
	mySQLClient *MySQLClient
	mock        sqlmock.Sqlmock
}

func (s *MySQLClientTestSuite) SetupSuite() {
	// Set up the test
	db, mock, err := sqlmock.New()
	if err != nil {
		s.T().Fatal(err)
	}
	sqlxDB := sqlx.NewDb(db, "mock")
	s.mySQLClient = &MySQLClient{db: sqlxDB}
	s.mock = mock
}

func (s *MySQLClientTestSuite) TestMySQLDAL_QueryRow() {

	// Define the expected query and result
	s.mock.ExpectQuery(`SELECT \* FROM users WHERE id = \?`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email"}).
			AddRow(1, "John Doe", "john.doe@example.com"))

	// Call the QueryRow method
	row := s.mySQLClient.QueryRow(context.Background(), "SELECT * FROM users WHERE id = ?", 1)

	// Verify that the result is correct
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		log.Fatal(err)
	}
	if user.ID != 1 || user.Name != "John Doe" || user.Email != "john.doe@example.com" {
		fmt.Errorf("Expected user %+v, got %+v", User{ID: 1, Name: "John Doe", Email: "john.doe@example.com"}, user)
	}
}

func TestMySQLClientTestSuite(t *testing.T) {
	suite.Run(t, new(MySQLClientTestSuite))
}
