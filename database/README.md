# database

This package provides a flexible and configurable database manager for Go applications, built on top of `jmoiron/sqlx`. It supports read/write splitting, connection pooling, and robust transaction management with automatic rollback on errors or panics.

## Features

- **`Manager` Struct**: Centralizes database connection management, supporting separate connections for read and write operations.
- **Configurable Connection Pool**: Allows setting maximum open connections, idle connections, and connection lifetime.
- **Transaction Management**: The `WithTransaction` function provides a safe way to execute database operations within a transaction, ensuring atomicity. It handles `BEGIN`, `COMMIT`, `ROLLBACK`, and panic recovery automatically.
- **`Queryer` Interface**: Defines common database query operations (`GetContext`, `SelectContext`, `ExecContext`, `NamedExecContext`), making it easy to work with both `sqlx.DB` and `sqlx.Tx` instances.
- **DSN Building**: Helper function to construct DSN (Data Source Name) strings from a structured configuration.
- **Testability**: Designed to be easily testable with mock database connections.

## Installation

To use this package, you need to have Go installed. You also need to install `sqlx` and your chosen database driver (e.g., `github.com/go-sql-driver/mysql` or `github.com/lib/pq`).

```bash
go get github.com/meowmix1337/go-core/database
go get github.com/jmoiron/sqlx
# Install your database driver, e.g., for MySQL:
go get github.com/go-sql-driver/mysql
# Or for PostgreSQL:
go get github.com/lib/pq
```

## Usage

### Configuration

Define your database connection details using the `Config` struct.

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/meowmix1337/go-core/database"
)

func main() {
	// Example configuration for a MySQL database
	writerConfig := database.Config{
		Host:     os.Getenv("DB_WRITER_HOST"),
		Port:     3306,
		User:     os.Getenv("DB_WRITER_USER"),
		Password: os.Getenv("DB_WRITER_PASSWORD"),
		DBName:   os.Getenv("DB_WRITER_NAME"),
		// InsecureSkipVerify: true, // Use with caution, only for development/testing
	}

	// Optional: Configuration for a read replica
	var readerConfig *database.Config
	if os.Getenv("DB_READER_HOST") != "" {
		readerConfig = &database.Config{
			Host:     os.Getenv("DB_READER_HOST"),
			Port:     3306,
			User:     os.Getenv("DB_READER_USER"),
			Password: os.Getenv("DB_READER_PASSWORD"),
			DBName:   os.Getenv("DB_READER_NAME"),
		}
	}

	// Initialize the database manager
	manager, err := database.NewManager(
		"mysql", // or "postgres"
		writerConfig,
		readerConfig,
		database.WithMaxOpenConns(20),
		database.WithMaxIdleConns(10),
		database.WithConnMaxLifetime(5*time.Minute),
		database.WithQueryTimeout(3*time.Second),
	)
	if err != nil {
		log.Fatalf("Failed to initialize database manager: %v", err)
	}
	defer manager.Close() // Ensure connections are closed when the application exits

	fmt.Println("Database manager initialized successfully.")

	// You can now use manager.Reader() and manager.Writer() for queries
	// See examples below.
}
```

### Performing Queries

The `Manager` provides `Reader()` and `Writer()` methods that return `*sqlx.DB` instances, which implement the `Queryer` interface. You can use these to perform standard `sqlx` operations.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/meowmix1337/go-core/database"
)

type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`	
	Email string `db:"email"`
}

func main() {
	// Assume manager is initialized as in the Configuration example
	writerConfig := database.Config{
		Host:     os.Getenv("DB_WRITER_HOST"),
		Port:     3306,
		User:     os.Getenv("DB_WRITER_USER"),
		Password: os.Getenv("DB_WRITER_PASSWORD"),
		DBName:   os.Getenv("DB_WRITER_NAME"),
	}
	manager, err := database.NewManager("mysql", writerConfig, nil)
	if err != nil {
		log.Fatalf("Failed to initialize database manager: %v", err)
	}
	defer manager.Close();

	ctx := context.Background()

	// Example: Get a single user from the reader (read-only)
	var user User
	err = manager.Reader().GetContext(ctx, &user, "SELECT id, name, email FROM users WHERE id = ?", 1)
	if err != nil {
		log.Printf("Error getting user: %v", err)
	} else {
		fmt.Printf("Found user (read-only): %+v\n", user)
	}

	// Example: Select multiple users from the reader
	var users []User
	err = manager.Reader().SelectContext(ctx, &users, "SELECT id, name, email FROM users LIMIT 10")
	if err != nil {
		log.Printf("Error selecting users: %v", err)
	} else {
		fmt.Printf("Found %d users (read-only).\n", len(users))
	}

	// Example: Insert a new user using the writer (write operations)
	result, err := manager.Writer().ExecContext(ctx, "INSERT INTO users (name, email) VALUES (?, ?)", "John Doe", "john.doe@example.com")
	if err != nil {
		log.Printf("Error inserting user: %v", err)
	} else {
		lastID, _ := result.LastInsertId()
		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("Inserted new user with ID %d, rows affected: %d\n", lastID, rowsAffected)
	}

	// Example: Update a user using the writer
	result, err = manager.Writer().ExecContext(ctx, "UPDATE users SET name = ? WHERE id = ?", "Jane Doe", 1)
	if err != nil {
		log.Printf("Error updating user: %v", err)
	} else {
		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("Updated %d rows.\n", rowsAffected)
	}
}
```

### Using Transactions

The `WithTransaction` function ensures that a series of database operations are executed atomically.

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/meowmix1337/go-core/database"
)

func main() {
	// Assume manager is initialized as in the Configuration example
	writerConfig := database.Config{
		Host:     os.Getenv("DB_WRITER_HOST"),
		Port:     3306,
		User:     os.Getenv("DB_WRITER_USER"),
		Password: os.Getenv("DB_WRITER_PASSWORD"),
		DBName:   os.Getenv("DB_WRITER_NAME"),
	}
	manager, err := database.NewManager("mysql", writerConfig, nil)
	if err != nil {
		log.Fatalf("Failed to initialize database manager: %v", err)
	}
	defer manager.Close()

	ctx := context.Background()

	// Example: Transfer funds in a transaction
	err = manager.WithTransaction(ctx, func(tx database.Queryer) error {
		// Deduct from sender
		_, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - ? WHERE id = ?", 100, 1)
		if err != nil {
			return fmt.Errorf("failed to deduct from account 1: %w", err)
		}

		// Simulate a condition that might cause a rollback
		// if someCondition {
		// 	return errors.New("simulated error, transaction will rollback")
		// }

		// Add to receiver
		_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + ? WHERE id = ?", 100, 2)
		if err != nil {
			return fmt.Errorf("failed to add to account 2: %w", err)
		}

		return nil // If no error, transaction commits
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
	} else {
		fmt.Println("Transaction completed successfully.")
	}

	// Example: Transaction with a panic (will also rollback)
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic in transaction: %v\n", r)
			}
		}()

		err = manager.WithTransaction(ctx, func(tx database.Queryer) error {
			log.Println("Inside transaction, about to panic...")
			panic("something went terribly wrong!")
		})
		if err != nil {
			log.Printf("Transaction with panic returned error: %v", err)
		}
	}()
}
```

## Contributing

Feel free to open issues or pull requests if you have suggestions or improvements.
