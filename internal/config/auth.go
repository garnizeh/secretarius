package config

import (
	"time"
)

type AuthConfig struct {
	JWTSecretKey    string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	CleanupInterval time.Duration
}

func LoadAuthConfig() *AuthConfig {
	return &AuthConfig{
		JWTSecretKey:    getEnv("JWT_SECRET_KEY", "your-secret-key-here"),
		AccessTokenTTL:  getDurationEnv("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
		RefreshTokenTTL: getDurationEnv("JWT_REFRESH_TOKEN_TTL", 30*24*time.Hour),
		CleanupInterval: getDurationEnv("JWT_DENYLIST_CLEANUP_INTERVAL", 24*time.Hour),
	}
}
