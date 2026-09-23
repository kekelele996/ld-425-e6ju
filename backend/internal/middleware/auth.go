package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/service"
	"github.com/home-renovation/platform/internal/utils"
)

// ContextKey 中间件上下文键。
type ContextKey string

const (
	ContextUserID   ContextKey = "user_id"
	ContextUsername ContextKey = "username"
	ContextRole     ContextKey = "role"
)

// AuthMiddleware JWT 鉴权中间件。
func AuthMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			utils.AbortError(c, http.StatusUnauthorized, constants.CodeUnauthorized, "missing or invalid authorization header")
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		user, err := userService.ParseToken(token)
		if err != nil {
			utils.AbortError(c, http.StatusUnauthorized, constants.CodeUnauthorized, "invalid or expired token")
			return
		}
		c.Set(string(ContextUserID), user.ID)
		c.Set(string(ContextUsername), user.Username)
		c.Set(string(ContextRole), user.Role)
		c.Next()
	}
}
