package verification

type Verifier[T any] interface {
	GenerateCode(val T) string
	VerifyCode(code string) bool
}

