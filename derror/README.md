# derror
This is just an Error wrapper for the Go Error interface. Please don't actually use this as this is for my own learning and testing purposes.

## Add the dependency
```
go get github.com/meowmix1337/go-core
```

## Usage
Basic Error
```go
import "github.com/meowmix1337/go-core/derror"

err := derror.New(context.Background(), derror.InternalServerCode, derror.InternalType, "failed to do something", errors.New("some error"))
```

Retryable Error
```go
import "github.com/meowmix1337/go-core/derror"

err := derror.NewRetryable(context.Background(), derror.InternalServerCode, derror.InternalType, "failed to do something", errors.New("some error"))
```