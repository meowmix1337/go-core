# http util
Provides a utility functions with http responses, etc

## Usage
JSONResponse
```go
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Use the JSONResponse function to return the struct as JSON
http_util.JSONResponse(w, http.StatusOK, user)
```

The output will be:
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john.doe@example.com"
}