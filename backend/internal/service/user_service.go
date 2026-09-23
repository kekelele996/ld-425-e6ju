package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/home-renovation/platform/internal/config"
	"github.com/home-renovation/platform/internal/constants"
	"github.com/home-renovation/platform/internal/dto"
	"github.com/home-renovation/platform/internal/model"
	"github.com/home-renovation/platform/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

// UserService 认证与用户服务接口。
type UserService interface {
	Login(username, password string) (*dto.LoginResponse, error)
	ParseToken(token string) (*model.User, error)
	GetByID(id uint) (*model.User, error)
	SeedIfEmpty() error
}

type userService struct {
	repo   repository.UserRepository
	cfg    config.JWTConfig
	logger *slog.Logger
}

// NewUserService 构造用户服务。
func NewUserService(repo repository.UserRepository, cfg config.JWTConfig, logger *slog.Logger) UserService {
	return &userService{repo: repo, cfg: cfg, logger: logger}
}

func (s *userService) Login(username, password string) (*dto.LoginResponse, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("login: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	token, err := s.issueToken(user)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	return &dto.LoginResponse{
		Token: token,
		User: dto.UserDTO{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
			Name:     user.Name,
		},
	}, nil
}

func (s *userService) issueToken(user *model.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":      user.ID,
		"username": user.Username,
		"role":     user.Role,
		"iat":      now.Unix(),
		"exp":      now.Add(s.cfg.ExpireDuration()).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

func (s *userService) ParseToken(tokenString string) (*model.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	idFloat, _ := claims["sub"].(float64)
	username, _ := claims["username"].(string)
	role, _ := claims["role"].(string)
	return &model.User{ID: uint(idFloat), Username: username, Role: role}, nil
}

func (s *userService) GetByID(id uint) (*model.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (s *userService) SeedIfEmpty() error {
	users, err := s.repo.List()
	if err != nil {
		return fmt.Errorf("seed users list: %w", err)
	}
	if len(users) > 0 {
		return nil
	}
	seedUsers := []struct {
		username, password, role, name string
	}{
		{"admin", "Admin123456", constants.RoleAdmin, "系统管理员"},
		{"designer", "Designer123", constants.RoleDesigner, "设计师小王"},
		{"contractor", "Contractor123", constants.RoleContractor, "施工队长老李"},
		{"owner", "Owner123456", constants.RoleOwner, "业主陈先生"},
		{"pm", "Manager123456", constants.RoleProjectManager, "项目经理赵工"},
	}
	for _, item := range seedUsers {
		hash, err := bcrypt.GenerateFromPassword([]byte(item.password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("seed user %s hash: %w", item.username, err)
		}
		user := &model.User{Username: item.username, PasswordHash: string(hash), Role: item.role, Name: item.name}
		if err := s.repo.Create(user); err != nil {
			return fmt.Errorf("seed user %s: %w", item.username, err)
		}
		s.logger.Info("seeded user", "username", item.username, "role", item.role)
	}
	return nil
}
