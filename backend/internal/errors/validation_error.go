package errors

import (
	"fmt"
	"strings"
)

// ValidationError 聚合字段校验错误。
type ValidationError struct {
	Fields map[string]string
}

// Error 实现 error 接口。
func (e *ValidationError) Error() string {
	parts := make([]string, 0, len(e.Fields))
	for k, v := range e.Fields {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(parts, "; ")
}

// NewValidation 创建校验错误。
func NewValidation(fields map[string]string) *ValidationError {
	return &ValidationError{Fields: fields}
}
