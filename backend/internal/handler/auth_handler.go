package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/dto"
	"github.com/home-renovation/platform/internal/middleware"
	"github.com/home-renovation/platform/internal/service"
	"github.com/home-renovation/platform/internal/utils"
)

// AuthHandler 认证处理器。
type AuthHandler struct {
	userService service.UserService
}

// NewAuthHandler 构造认证处理器。
func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

// Login 用户登录。
// @Summary 用户登录
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "登录请求"
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	resp, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.AbortError(c, http.StatusUnauthorized, constants.CodeUnauthorized, "invalid username or password")
			return
		}
		c.Error(err)
		return
	}
	utils.Success(c, resp)
}

// Me 获取当前登录用户信息。
// @Summary 当前用户
// @Tags Auth
// @Produce json
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	id, _ := c.Get(string(middleware.ContextUserID))
	user, err := h.userService.GetByID(id.(uint))
	if err != nil {
		c.Error(err)
		return
	}
	utils.Success(c, dto.UserDTO{ID: user.ID, Username: user.Username, Role: user.Role, Name: user.Name})
}
