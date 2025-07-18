# httputil

This package provides utility functions for common HTTP operations, making it easier to handle HTTP responses in Go applications.

## Features

- **`JSONResponse` Function**: A helper function to write JSON responses with a specified status code and data. It automatically sets the `Content-Type` header to `application/json` and handles JSON marshaling.

## Installation

To use this package, you need to have Go installed. Then, you can add it to your project:

```bash
go get github.com/meowmix1337/go-core/httputil
```

## Usage

### Sending JSON Responses

The `JSONResponse` function simplifies sending structured JSON data back to the client.

```go
package main

import (
	"net/http"
	"github.com/meowmix1337/go-core/httputil"
)

func handler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "success", "message": "Operation completed."}
	httputil.JSONResponse(w, http.StatusOK, data)
}

func main() {
	http.HandleFunc("/api/data", handler)
	http.ListenAndServe(":8080", nil)
}
```

In this example, `JSONResponse` is used within an HTTP handler to send a `200 OK` response with a JSON payload. If there's an error during JSON encoding, it gracefully falls back to sending a `500 Internal Server Error`.

## Contributing

Feel free to open issues or pull requests if you have suggestions or improvements.
