package verification

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
)

type EmailVerifier struct {
	masterKey []byte
}

const separator = "."

func NewEmailVerifier(masterKey string) *EmailVerifier {
	return &EmailVerifier{
		masterKey: []byte(masterKey),
	}
}

var _ Verifier[string] = (*EmailVerifier)(nil)

// GenerateCode generates a code for the given email.
func (v *EmailVerifier) GenerateCode(email string) string {
	// create the signature
	mac := hmac.New(sha512.New, v.masterKey)
	mac.Write([]byte(email))
	signature := mac.Sum(nil)

	// separate the email and the signature
	encodedEmail := hex.EncodeToString([]byte(email))
	encodedSignature := hex.EncodeToString([]byte(signature))

	// combine the email and the signature
	token := fmt.Sprintf("%s%s%s", encodedEmail, separator, encodedSignature)

	return token
}

func (v *EmailVerifier) VerifyCode(code string) bool {
	// separate the email and the signature
	parts := strings.Split(code, separator)
	if len(parts) != 2 {
		return false
	}

	// decode the email and the signature
	email, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}
	signature, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	// create the signature
	mac := hmac.New(sha512.New, v.masterKey)
	mac.Write(email)
	expectedSignature := mac.Sum(nil)

	// compare the signature
	return hmac.Equal(signature, expectedSignature)
}
