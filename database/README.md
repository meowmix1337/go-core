# database
This is a interface to allow me to use any relational database as needed. This will have mysql and postgres clients

## Add the dependency
```
go get github.com/meowmix1337/go-core
```

## Usage
MySQL
```go
import "github.com/meowmix1337/go-core/derror"

db, err := database.NewMySQLClient(dsn)
```

Postgress
```go
import "github.com/meowmix1337/go-core/derror"

db, err := database.NewMyPostgresClient(dsn)
```

Transactions
```go
tx, err := db.BeginTx(ctx)

result, err := db.ExecTx(ctx, tx, query, 1, 2, 3)

if err := db.CommitTx(tx); err != nil {
    // log something here
}
```
