package derror

type Type string

const (
	InternalType   Type = "INTERNAL_ERROR"
	BadRequestType Type = "BAD_REQUEST"
)
