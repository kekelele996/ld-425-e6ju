package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

// Config 聚合应用全部环境配置，集中解析。
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Upload   UploadConfig
}

// ServerConfig 服务监听配置。
type ServerConfig struct {
	Port int    `env:"SERVER_PORT" envDefault:"8080"`
	Mode string `env:"GIN_MODE" envDefault:"release"`
}

// Load 从环境变量读取配置。
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.ParseWithOptions(cfg, env.Options{Prefix: "", Environment: nil}); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	if err := cfg.Database.Validate(); err != nil {
		return nil, fmt.Errorf("validate database config: %w", err)
	}
	if cfg.JWT.Secret == "" {
		log.Println("[config] JWT_SECRET is empty, using development fallback")
		cfg.JWT.Secret = "home-renovation-dev-secret-change-me"
	}
	return cfg, nil
}
