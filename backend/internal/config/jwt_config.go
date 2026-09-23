package config

import "time"

// JWTConfig JWT 签发与校验配置。
type JWTConfig struct {
	Secret      string        `env:"JWT_SECRET"`
	ExpireHours time.Duration `env:"JWT_EXPIRE_HOURS" envDefault:"24h"`
}

// ExpireDuration 返回 token 有效期。
func (c JWTConfig) ExpireDuration() time.Duration {
	if c.ExpireHours <= 0 {
		return 24 * time.Hour
	}
	return c.ExpireHours
}
