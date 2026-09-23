package config

import (
	"fmt"
	"time"
)

// DatabaseConfig 数据库连接配置。
type DatabaseConfig struct {
	Host            string        `env:"DB_HOST" envDefault:"127.0.0.1"`
	Port            int           `env:"DB_PORT" envDefault:"3306"`
	User            string        `env:"DB_USER" envDefault:"home_renovation"`
	Password        string        `env:"DB_PASSWORD" envDefault:"home_renovation"`
	Name            string        `env:"DB_NAME" envDefault:"home_renovation"`
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"50"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"1h"`
}

// DSN 生成 GORM MySQL 数据源名称。
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.User, c.Password, c.Host, c.Port, c.Name)
}

// Validate 校验必填配置。
func (c DatabaseConfig) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	return nil
}
