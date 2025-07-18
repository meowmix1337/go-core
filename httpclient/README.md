# httpclient

This package provides a simple and flexible HTTP client for making requests to external APIs. It includes a base `HttpClient` interface and an implementation that handles common HTTP methods (GET, POST, PUT, DELETE), JSON payload marshaling, and query parameter encoding.

## Features

- **`HttpClient` Interface**: Defines the contract for HTTP request operations.
- **`httpClient` Implementation**: A concrete implementation of `HttpClient` that handles:
    - Constructing URLs from a base URL and API prefix.
    - Marshaling Go structs/maps into JSON payloads for POST/PUT/DELETE requests.
    - Encoding query parameters for GET requests.
    - Basic error handling for non-2xx HTTP status codes.
- **`MyClient` Wrapper**: Provides convenience methods (`Get`, `Post`, `Put`, `Delete`) for common HTTP operations, simplifying usage.

## Installation

To use this package, you need to have Go installed. Then, you can add it to your project:

```bash
go get github.com/meowmix1337/go-core/httpclient
```

## Usage

### Basic Usage with `MyClient`

The `MyClient` struct provides a convenient way to interact with the `httpclient` package.

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/meowmix1337/go-core/httpclient"
)

func main() {
	// Initialize MyClient with your API's base URL and prefix
	client := httpclient.New("https://jsonplaceholder.typicode.com", "")

	// Example GET request
	fmt.Println("--- GET Request ---")
	queryParams := map[string]string{"userId": "1"}
	resp, err := client.Get(context.Background(), "/posts", queryParams)
	if err != nil {
		log.Fatalf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Printf("GET Response Status: %s\n", resp.Status)
	fmt.Printf("GET Response Body: %s\n", body)

	// Example POST request
	fmt.Println("\n--- POST Request ---")
	payload := map[string]interface{}{
		"title":  "foo",
		"body":   "bar",
		"userId": 1,
	}
	resp, err = client.Post(context.Background(), "/posts", payload)
	if err != nil {
		log.Fatalf("POST request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Printf("POST Response Status: %s\n", resp.Status)
	fmt.Printf("POST Response Body: %s\n", body)

	// Example PUT request
	fmt.Println("\n--- PUT Request ---")
	putPayload := map[string]interface{}{
		"id":     1,
		"title":  "foo updated",
		"body":   "bar updated",
		"userId": 1,
	}
	resp, err = client.Put(context.Background(), "/posts/1", putPayload)
	if err != nil {
		log.Fatalf("PUT request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Printf("PUT Response Status: %s\n", resp.Status)
	fmt.Printf("PUT Response Body: %s\n", body)

	// Example DELETE request
	fmt.Println("\n--- DELETE Request ---")
	resp, err = client.Delete(context.Background(), "/posts/1", nil)
	if err != nil {
		log.Fatalf("DELETE request failed: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("DELETE Response Status: %s\n", resp.Status)
	// For DELETE, typically no body is returned on success (204 No Content)
	if resp.StatusCode != http.StatusNoContent {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response body: %v", err)
		}
		fmt.Printf("DELETE Response Body: %s\n", body)
	}
}
```

### Using `HttpClient` Directly (Advanced)

You can also use the `httpClient` struct directly if you need more fine-grained control or want to implement your own wrapper.

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/meowmix1337/go-core/httpclient"
)

func main() {
	// Initialize the low-level httpClient
	client := &httpclient.httpClient{
		BaseURL:   "https://jsonplaceholder.typicode.com",
		APIPrefix: "",
	}

	// Make a custom request
	resp, err := client.Request(context.Background(), http.MethodGet, "/users/1", nil, nil)
	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", body)
}
```

## Contributing

Feel free to open issues or pull requests if you have suggestions or improvements.
