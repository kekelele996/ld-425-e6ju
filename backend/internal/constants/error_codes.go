package constants

// 统一错误码，集中维护。
const (
	CodeOK              = 0
	CodeBadRequest      = 40000
	CodeUnauthorized    = 40100
	CodeForbidden       = 40300
	CodeNotFound        = 40400
	CodeConflict        = 40900
	CodeValidationError = 42200
	CodeInternalError   = 50000
)
