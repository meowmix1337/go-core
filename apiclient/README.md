# apiclient

This package provides a flexible and robust API client for making HTTP requests to external services. It supports various HTTP methods, custom headers, query parameters, and request bodies, along with configurable options like base URL and timeouts.

## Features

- **`APIClient`**: A configurable HTTP client that simplifies making API calls.
- **`Request` Struct**: Defines the parameters for an API request, including method, endpoint, body, query parameters, and headers.
- **`RequestMethod` Type**: Strongly typed HTTP methods (GET, POST, PUT, DELETE) for clarity and type safety.
- Configurable Options:
    - `WithBaseURL`: Sets the base URL for all requests made by the client.
    - `WithTimeout`: Configures the timeout for HTTP requests.
    - `WithAuthorization`: Adds an Authorization header (e.g., Bearer token) to all requests.
- **Convenience Methods**: `Get`, `Post`, `Put`, `Delete` methods for common HTTP operations.
- **Error Handling**: Provides detailed error messages for network issues, non-2xx status codes, and JSON decoding failures.

## Installation

To use this package, you need to have Go installed. Then, you can add it to your project:

```bash
go get github.com/meowmix1337/go-core/apiclient
```

## Usage

### Initializing the API Client

You can initialize the `APIClient` with various options:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/meowmix1337/go-core/apiclient"
)

func main() {
	// Initialize client with a base URL and a 10-second timeout
	client := apiclient.New(
		apiclient.WithBaseURL("https://api.example.com"),
		apiclient.WithTimeout(10*time.Second),
		apiclient.WithAuthorization("your-auth-token"),
	)

	// Example usage (see below for more details)
	_ = client // Placeholder to avoid unused variable error
	fmt.Println("API Client initialized.")
}
```

### Making Requests

#### GET Request

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/meowmix1337/go-core/apiclient"
)

type Post struct {
	UserID int    `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {
	client := apiclient.New(apiclient.WithBaseURL("https://jsonplaceholder.typicode.com"))

	var posts []Post
	queryParams := map[string]string{"userId": "1"}
	headers := map[string]string{"X-Custom-Header": "value"}

	err := client.Get(context.Background(), "/posts", queryParams, headers, &posts)
	if err != nil {
		log.Fatalf("GET request failed: %v", err)
	}

	fmt.Printf("Fetched %d posts for userId 1:\n", len(posts))
	for _, post := range posts {
		fmt.Printf("  ID: %d, Title: %s\n", post.ID, post.Title)
	}
}
```

#### POST Request

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/meowmix1337/go-core/apiclient"
)

type NewPost struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

type CreatedPost struct {
	ID int `json:"id"`
	NewPost
}

func main() {
	client := apiclient.New(apiclient.WithBaseURL("https://jsonplaceholder.typicode.com"))

	newPost := NewPost{
		Title:  "foo",
		Body:   "bar",
		UserID: 1,
	}
	var createdPost CreatedPost
	headers := map[string]string{"Content-Type": "application/json"}

	err := client.Post(context.Background(), "/posts", newPost, headers, &createdPost)
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}

	fmt.Printf("Created new post with ID: %d, Title: %s\n", createdPost.ID, createdPost.Title)
}
```

#### PUT Request

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/meowmix1337/go-core/apiclient"
)

type UpdatedPost struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	UserID int    `json:"userId"`
}

func main() {
	client := apiclient.New(apiclient.WithBaseURL("https://jsonplaceholder.typicode.com"))

	updateData := map[string]interface{}{
		"id":    1,
		"title": "updated title",
		"body":  "updated body",
	}
	var updatedPost UpdatedPost

	err := client.Put(context.Background(), "/posts/1", updateData, nil, &updatedPost)
	if err != nil {
		log.Fatalf("PUT request failed: %v", err)
	}

	fmt.Printf("Updated post ID: %d, New Title: %s\n", updatedPost.ID, updatedPost.Title)
}
```

#### DELETE Request

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/meowmix1337/go-core/apiclient"
)

func main() {
	client := apiclient.New(apiclient.WithBaseURL("https://jsonplaceholder.typicode.com"))

	// For DELETE, the result can often be nil if no content is expected back
	err := client.Delete(context.Background(), "/posts/1", nil, nil)
	if err != nil {
		log.Fatalf("DELETE request failed: %v", err)
	}

	fmt.Println("Post with ID 1 deleted successfully.")
}
```

### Handling Errors

The `apiclient` package returns errors for various scenarios:

- **Network errors**: If the request fails due to network issues.
- **Non-2xx HTTP status codes**: If the API returns an error status (e.g., 400, 404, 500). The error message will include the status code and the response body (if available).
- **JSON decoding errors**: If the response body cannot be unmarshaled into the provided result interface.

Always check the returned `error` from `Do`, `Get`, `Post`, `Put`, and `Delete` methods.

## Contributing

Feel free to open issues or pull requests if you have suggestions or improvements.
