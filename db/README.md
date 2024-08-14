# database
This is a interface to allow me to use any relational database as needed. This will have mysql and postgres clients

## Add the dependency
```
go get github.com/meowmix1337/go-core
```

## Purpose
1. Can connect to different databases (specifically MySQL or Postgres) and swap between each one with ease
2. Support for Writer and Reader.
   1. If no Reader is configured, it will default to Writer even if you call the Reader
3. Easy transaction handling. Use the `Transaction` function to execute a sequence in a transaction

## Usage
MySQL
```go
import "github.com/meowmix1337/go-core/database"

// readerDSN is optional and can be set to "".
db, err := db.NewMySQL(writerDSN, readerDSN)
```

Postgress
```go
import "github.com/meowmix1337/go-core/database"

// readerDSN is optional and can be set to "".
db, err := database.NewPostgress(writerDSN, readerDSN)
```

Example Repo
```go
type UserRepo struct {
    db *db.DB // DB interface
}

func NewUserRepo(db *db.DB) *UserRepo {
    return &UserRepo{
        db: db
    }
}

func (r *UserRepo) DeleteUser(ctx context.Context, userID int) error {
    err := r.db.Transaction(context.Background(), func(ctx context.Context, tx db.Tx) error {

        // execute your queries here
        result, err := tx.Exec(ctx, "INSERT ...", 1, 2, 3)
        if err != nil {
            return err
        }

        var user User
        result, err := tx.Get(ctx, &user, "SELECT ...", userID)
        if err != nil {
            return err
        }

        // you can call additional repo functions here too
        // but it is up to you to allow the Tx to be passed through as a function param
        result, err := r.Delete(ctx, tx, someID)
        if err != nil {
            return err
        }

        // return nil if everything worked
        return nil
    })

    if err != nil {
        // commit/rollback error probably happened
    }
}
```

Example Service
```go

type UserService interface {
    ...
}

// Having the DB in the service allows for more flexiblity when creating transactions
// since the repos will always be on the service layer, passing the DB through the repo
// can be easier and cleaner to understand how to complete a transaction
type userService struct {
    db db.DB // DB interface
    userRepo UserRepo
    orderRepo OrderRepo
}

func NewUserService(db db.DB, userRepo UserRepo, orderRepo OrderRepo) UserService {
    return &userService{
        db db.DB // DB interface
        userRepo: userRepo,
        orderRepo: orderRepo,
    }
}

func (u *userService) DeleteUser(ctx context.Context, userID, orderID int) error {
    err := u.db.Transaction(context.Background(), func(ctx context.Context, tx db.Tx) error {

        err := u.OrderRepo.DeleteByID(ctx, tx, orderID)
        if err != nil {
            return err
        }

       err = u.UserRepo.DeleteByID(ctx, tx, orderID)
       if err != nil {
            return err
       }

        // return nil if everything worked
        return nil
    })

    if err != nil {
        // commit/rollback error probably happened
    }
}
```