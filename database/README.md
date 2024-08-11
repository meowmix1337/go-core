# database
This is a interface to allow me to use any relational database as needed. This will have mysql and postgres clients

## Add the dependency
```
go get github.com/meowmix1337/go-core
```

## Purpose
1. Can connect to different databases (specifically MySQL or Postgres)
2. Be explicit on what you doing. i.e. A repo function is specifically for a transaction, suffix with Tx `CreateTx`, `DeleteTx`. If no transaction is needed for some of these functions then they should have a non transactional one `Create`, `Delete`
   1. This will mean more code and maybe duplicate code (you should be using a constant for a query)
   2. Less mistakes when using transactions (I've seen people not pass in a tx when it was required so data integrity breaks)
3. Use the `WithTransaction` function if you don't want to explicitly Start, commit or rollback.
4. The Repository struct should have the `Database`, never the service.

## Usage
MySQL
```go
import "github.com/meowmix1337/go-core/database"

// readerDSN is optional and can be set to "".
db, err := database.NewDBConnector("mysql", writerDSN, readerDSN)
```

Postgress
```go
import "github.com/meowmix1337/go-core/database"

// readerDSN is optional and can be set to "".
db, err := database.NewDBConnector("postgres", writerDSN, readerDSN)
```

Repository
```go
// some repo you have
func (r *UserRepoistory) Create(ctx context.Context, user User) (int, error) {
    // ...

    // WriteDB to use writer, ReadDB to use reader
    result, err := r.db.WriteDB().ExecContext(ctx, "INSERT INTO user (name, email) VALUES (?, ?)", user.Name, user.Email)

    // ...
}

// transaction specific query
func (r *UserRepoistory) CreateTx(ctx context.Context, tx *sqlx.Tx, user User) (int, error) {
    // ...

    result, err := tx.ExecContext(ctx, "INSERT INTO user (name, email) VALUES (?, ?)", user.Name, user.Email)

    // ...
}
```

Transactions
```go
err := Transaction(context.Background(), s.mySQLClient, func(ctx context.Context, tx *sqlx.Tx) error {

    result, err := repo.UserRepo.CreateTx(ctx, tx, userInput)
    if err != nil {
        return err
    }

    result, err := repo.UserRepo.ByIDTx(ctx, tx, id)
    if err != nil {
        return err
    }

    result, err := repo.UserRepo.DeleteTx(ctx, tx, someID)
    if err != nil {
        return err
    }

    // return nil if everything worked
    return nil
})

if err != nil {
    // commit/rollback error probably happened
}
```
