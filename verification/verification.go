package verification

type Verifier[T any] interface {
	GenerateCode(val T) string
}