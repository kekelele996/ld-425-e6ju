package errors

import "fmt"

// BusinessError 携带业务错误码的错误类型。
type BusinessError struct {
	Code    int
	Message string
	HTTP    int
}

// Error 实现 error 接口。
func (e *BusinessError) Error() string {
	return fmt.Sprintf("business error [%d]: %s", e.Code, e.Message)
}

// NewBusiness 创建业务错误。
func NewBusiness(code, http int, message string) *BusinessError {
	return &BusinessError{Code: code, Message: message, HTTP: http}
}

// NewBadRequest 创建参数错误。
func NewBadRequest(message string) *BusinessError {
	return NewBusiness(40000, 400, message)
}

// NewNotFound 创建资源不存在错误。
func NewNotFound(message string) *BusinessError {
	return NewBusiness(40400, 404, message)
}

// NewForbidden 创建无权限错误。
func NewForbidden(message string) *BusinessError {
	return NewBusiness(40300, 403, message)
}

// NewConflict 创建状态冲突错误。
func NewConflict(message string) *BusinessError {
	return NewBusiness(40900, 409, message)
}

// NewInternal 创建服务内部错误。
func NewInternal(message string) *BusinessError {
	return NewBusiness(50000, 500, message)
}
