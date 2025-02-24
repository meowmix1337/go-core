package verification

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
)

type EmailVerifier struct {
	masterKey []byte
}

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

	return hex.EncodeToString([]byte(signature))
}
