package middleware

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/home-renovation/platform/internal/constants"
	apperrors "github.com/home-renovation/platform/internal/errors"
	"github.com/home-renovation/platform/internal/utils"
)

// ErrorHandlerMiddleware 统一异常处理与 panic 恢复。
func ErrorHandlerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic recovered", "path", c.Request.URL.Path, "panic", rec)
				utils.AbortError(c, http.StatusInternalServerError, constants.CodeInternalError, "internal server error")
			}
		}()
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		handleError(c, logger, err)
	}
}

func handleError(c *gin.Context, logger *slog.Logger, err error) {
	var businessErr *apperrors.BusinessError
	if errors.As(err, &businessErr) {
		utils.AbortError(c, businessErr.HTTP, businessErr.Code, businessErr.Message)
		return
	}
	var validationErr *apperrors.ValidationError
	if errors.As(err, &validationErr) {
		utils.AbortError(c, http.StatusUnprocessableEntity, constants.CodeValidationError, validationErr.Error())
		return
	}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		msg := buildValidationMessage(ve)
		utils.AbortError(c, http.StatusUnprocessableEntity, constants.CodeValidationError, msg)
		return
	}
	var numErr *strconv.NumError
	var typeErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError
	if errors.As(err, &numErr) || errors.As(err, &typeErr) || errors.As(err, &syntaxErr) {
		utils.AbortError(c, http.StatusBadRequest, constants.CodeBadRequest, "invalid request")
		return
	}
	logger.Error("unhandled error", "error", err)
	utils.AbortError(c, http.StatusInternalServerError, constants.CodeInternalError, "internal server error")
}

func buildValidationMessage(ve validator.ValidationErrors) string {
	if len(ve) == 0 {
		return "validation failed"
	}
	field := ve[0].Field()
	if field == "" {
		field = ve[0].StructField()
	}
	return field + " validation failed: " + ve[0].Tag()
}
