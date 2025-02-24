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
	s.True(emailVerifier.VerifyCode(code))
}

func (s *EmailVerificationSuite) TestInvalidCode() {
	emailVerifier := NewEmailVerifier("test")
	s.False(emailVerifier.VerifyCode("some invalid code"))
}
