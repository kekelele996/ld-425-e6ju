package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/utils"
)

// RBACMiddleware 角色校验中间件。
func RBACMiddleware(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		allowed[role] = struct{}{}
	}
	return func(c *gin.Context) {
		role, _ := c.Get(string(ContextRole))
		roleStr, _ := role.(string)
		if _, ok := allowed[roleStr]; !ok {
			utils.AbortError(c, http.StatusForbidden, constants.CodeForbidden, "forbidden: insufficient role")
			return
		}
		c.Next()
	}
}
