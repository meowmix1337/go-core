package verification

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type EmailVerificationSuite struct {
	suite.Suite
}

func TestEmailVerificationSuite(t *testing.T) {
	suite.Run(t, new(EmailVerificationSuite))
}

func (s *EmailVerificationSuite) SetupSuite() {}

func (s *EmailVerificationSuite) TestValidCode() {
	emailVerifier := NewEmailVerifier("test")
	code := emailVerifier.GenerateCode("test@example.com")
	s.Equal("15b4b2821571fa31c96101da5466be49d079c3eeb54b85437fbf2200303e39d64c7fbe319d8a54e3327fbc2bcedb16c9abd18461a7fe0d69f33430d914a1bdd7", code)
}
