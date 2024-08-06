# My Client
Just a basic Http Client so I can use this anywhere to interact with other APIs without having to rewrite this everytime.

## Add the dependency
```
go get github.com/meowmix1337/go-core
```

## Usage
Basic Usage
```go
import "github.com/meowmix1337/go-core/http_client"

httpClient := http_client.New("http://dog.ceo", "/api")
resp, err := httpClient.Get(context.Context, "/breeds/list/all", nil)

resp, err := httpClient.Get(context.Context, "/breeds/list/all", map[string]string{
    "breed": "corgi"
})
```