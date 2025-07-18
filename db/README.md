# db

This package provides a robust and flexible database abstraction layer for Go applications, supporting both MySQL and PostgreSQL. It simplifies common database operations, including read/write splitting, transactions, and panic recovery.

## Features

- **`DB` Interface**: Defines a common interface for database operations, allowing for interchangeable database implementations (MySQL, PostgreSQL).
- **Read/Write Splitting**: Supports separate connections for read and write operations, improving performance and scalability. If a read replica is not provided, read operations fall back to the write database.
- **Transactions**: Provides a `Transaction` method that handles the entire transaction lifecycle, including beginning, committing, and rolling back transactions, with built-in panic recovery.
- **`Tx` Interface**: Defines the contract for database transactions.
- **`MySQL` and `PostgreSQL` Implementations**: Concrete implementations for MySQL and PostgreSQL databases.
- **`sqlx` Integration**: Leverages `jmoiron/sqlx` for enhanced database operations, including easy scanning of query results into Go structs.
- **Structured Logging**: Integrates with `zerolog` for informative logging of database connections and transaction events.

## Installation

To use this package, you need to have Go installed. You also need to install the appropriate database drivers and `sqlx`:

```bash
go get github.com/meowmix1337/go-core/db
go get github.com/jmoiron/sqlx
# For MySQL
go get github.com/go-sql-driver/mysql
# For PostgreSQL
go get github.com/lib/pq
```

## Usage

### Initializing a Database Connection

You can initialize either a MySQL or PostgreSQL database connection. It's recommended to load DSNs (Data Source Names) from environment variables or a secure configuration.

#### MySQL

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/meowmix1337/go-core/db"
)

func main() {
	// Example DSNs (replace with your actual database connection strings)
	// For production, load these from environment variables or a secure config.
	writerDSN := os.Getenv("MYSQL_WRITER_DSN") // e.g., "user:password@tcp(127.0.0.1:3306)/dbname?parseTime=true"
	readerDSN := os.Getenv("MYSQL_READER_DSN") // Optional, can be empty if no read replica

	if writerDSN == "" {
		log.Fatal("MYSQL_WRITER_DSN environment variable not set.")
	}

	// Initialize MySQL database
	mysqlDB := db.NewMySQL(writerDSN, readerDSN)

	// Use the database (see examples below)
	_ = mysqlDB // Placeholder to avoid unused variable error
	fmt.Println("MySQL database initialized.")
}
```

#### PostgreSQL

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/meowmix1337/go-core/db"
)

func main() {
	// Example DSNs (replace with your actual database connection strings)
	// For production, load these from environment variables or a secure config.
	writerDSN := os.Getenv("POSTGRES_WRITER_DSN") // e.g., "host=localhost port=5432 user=user password=password dbname=dbname sslmode=disable"
	readerDSN := os.Getenv("POSTGRES_READER_DSN") // Optional, can be empty if no read replica

	if writerDSN == "" {
		log.Fatal("POSTGRES_WRITER_DSN environment variable not set.")
	}

	// Initialize PostgreSQL database
	postgresDB := db.NewPostgres(writerDSN, readerDSN)

	// Use the database (see examples below)
	_ = postgresDB // Placeholder to avoid unused variable error
	fmt.Println("PostgreSQL database initialized.")
}
```

### Performing Database Operations

Assume you have a `User` struct and a `db.DB` instance (e.g., `mysqlDB` or `postgresDB`).

```go
type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}
```

#### Get a Single Row

Use `Get` for write operations (e.g., after an insert/update) or `Get_RO` for read-only operations.

```go
// Example: Get a user by ID
var user User
err := mysqlDB.Get_RO(context.Background(), &user, "SELECT id, name, email FROM users WHERE id = ?", 1)
if err != nil {
	if err == sql.ErrNoRows {
		fmt.Println("User not found.")
	} else {
		log.Fatalf("Error getting user: %v", err)
	}
} else {
	fmt.Printf("Found user: %+v\n", user)
}
```

#### Select Multiple Rows

Use `Select` for write operations or `Select_RO` for read-only operations.

```go
// Example: Select all users
var users []User
err := mysqlDB.Select_RO(context.Background(), &users, "SELECT id, name, email FROM users")
if err != nil {
	log.Fatalf("Error selecting users: %v", err)
}
fmt.Printf("Found %d users:\n", len(users))
for _, u := range users {
	fmt.Printf("  %+v\n", u)
}
```

#### Execute Queries (Insert, Update, Delete)

Use `Exec` for DML (Data Manipulation Language) operations.

```go
// Example: Insert a new user
result, err := mysqlDB.Exec(context.Background(), "INSERT INTO users (name, email) VALUES (?, ?)", "Jane Doe", "jane.doe@example.com")
if err != nil {
	log.Fatalf("Error inserting user: %v", err)
}
lastInsertID, _ := result.LastInsertId()
rowsAffected, _ := result.RowsAffected()
fmt.Printf("Inserted new user with ID: %d, Rows Affected: %d\n", lastInsertID, rowsAffected)

// Example: Update a user
result, err = mysqlDB.Exec(context.Background(), "UPDATE users SET name = ? WHERE id = ?", "Jane Smith", lastInsertID)
if err != nil {
	log.Fatalf("Error updating user: %v", err)
}
rowsAffected, _ = result.RowsAffected()
fmt.Printf("Updated %d rows.\n", rowsAffected)

// Example: Delete a user
result, err = mysqlDB.Exec(context.Background(), "DELETE FROM users WHERE id = ?", lastInsertID)
if err != nil {
	log.Fatalf("Error deleting user: %v", err)
}
rowsAffected, _ = result.RowsAffected()
fmt.Printf("Deleted %d rows.\n", rowsAffected)
```

### Using Transactions

The `Transaction` method provides a safe way to perform a series of database operations within a single transaction. It automatically handles `Commit`, `Rollback` on error, and `Rollback` on panic.

```go
// Example: Transferring money between two accounts
func transferFunds(db db.DB, fromAccountID, toAccountID int, amount float64) error {
	return db.Transaction(context.Background(), func(ctx context.Context, tx db.Tx) error {
		// Deduct from sender
		_, err := tx.Exec(ctx, "UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromAccountID)
		if err != nil {
			return fmt.Errorf("failed to deduct from sender: %w", err)
		}

		// Simulate an error for demonstration
		// if amount > 1000 {
		// 	return errors.New("amount too high, forcing rollback")
		// }

		// Add to receiver
		_, err = tx.Exec(ctx, "UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toAccountID)
		if err != nil {
			return fmt.Errorf("failed to add to receiver: %w", err)
		}

		return nil // If no error, transaction will be committed
	})
}

// In main or another function:
// err := transferFunds(mysqlDB, 1, 2, 500.00)
// if err != nil {
// 	log.Printf("Transaction failed: %v", err)
// } else {
// 	fmt.Println("Funds transferred successfully.")
// }
```

## Contributing

Feel free to open issues or pull requests if you have suggestions or improvements.

```