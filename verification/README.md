# verification

This package provides utilities for generating and verifying codes, primarily for email verification purposes. It uses HMAC-SHA512 for secure code generation.

## Features

- **`EmailVerifier`**: A struct that generates secure, HMAC-SHA512 based verification codes for email addresses.
- **`Verifier` Interface**: Defines a generic interface for code generation, allowing for different types of verification (e.g., email, phone number).

## Installation

To use this package, you need to have Go installed. Then, you can add it to your project:

```bash
go get github.com/meowmix1337/go-core/verification
```

## Usage

### Email Verification

The `EmailVerifier` can be used to generate a unique and verifiable code for a given email address. This code can then be sent to the user (e.g., in an email) and later verified to confirm ownership of the email address.

```go
package main

import (
	"fmt"

	"github.com/meowmix1337/go-core/verification"
)

func main() {
	// Initialize the EmailVerifier with a strong, secret master key.
	// This key should be kept confidential and ideally loaded from environment variables or a secure configuration.
	masterKey := "your-super-secret-master-key-that-is-at-least-32-bytes-long"
	emailVerifier := verification.NewEmailVerifier(masterKey)

	email := "user@example.com"

	// Generate a verification code for the email
	code := emailVerifier.GenerateCode(email)
	fmt.Printf("Verification code for %s: %s\n", email, code)

	// In a real application, you would now send this 'code' to 'user@example.com'
	// via an email service.

	// --- Later, when the user tries to verify their email with the received code ---

	// Regenerate the expected code using the same master key and email
	expectedCode := emailVerifier.GenerateCode(email)

	// Compare the received code with the expected code
	if code == expectedCode {
		fmt.Println("Email successfully verified!")
	} else {
		fmt.Println("Invalid verification code.")
	}

	// Example with a different email (should produce a different code)
	anotherEmail := "another.user@example.com"
	anotherCode := emailVerifier.GenerateCode(anotherEmail)
	fmt.Printf("Verification code for %s: %s\n", anotherEmail, anotherCode)
	if code == anotherCode {
		fmt.Println("Error: Codes for different emails should not match.")
	}
}
```

**Important Security Note**: The `masterKey` used to initialize `EmailVerifier` is crucial for the security of your verification process. It should be a long, random, and securely stored secret. **Never hardcode it in your application code in a production environment.** Use environment variables, a secret management service, or a secure configuration system to provide this key.

## How it Works

The `EmailVerifier` uses HMAC (Hash-based Message Authentication Code) with SHA-512. This ensures that:
1. The generated code is unique for each email address (given the same master key).
2. The code cannot be forged without knowing the `masterKey`.
3. The original email address cannot be easily derived from the code.

When a user attempts to verify their email, the application regenerates the code using the same `masterKey` and the email address provided by the user. If the regenerated code matches the code the user submitted, the email is considered verified.

## Contributing

Feel free to open issues or pull requests if you have suggestions or improvements.

