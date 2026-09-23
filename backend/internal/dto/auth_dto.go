package dto

// LoginRequest 登录请求。
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

// LoginResponse 登录响应。
type LoginResponse struct {
	Token string `json:"token"`
	User  UserDTO `json:"user"`
}

// UserDTO 用户展示信息。
type UserDTO struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}
