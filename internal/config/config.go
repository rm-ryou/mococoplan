package config

import (
	"fmt"
	"os"
	"time"
)

type (
	Config struct {
		Port  string
		DB    DB
		Token Token
		Redis Redis
	}

	DB struct {
		Name     string
		User     string
		Password string
		Host     string
		Port     string
	}

	Token struct {
		AccessTokenSecret string
		AccessTokenIssuer string
		AccessTokenTTL    time.Duration
		RefreshTokenTTL   time.Duration
	}

	Redis struct {
		Host     string
		Port     string
		Password string
	}
)

func NewConfig() (*Config, error) {
	accessTokenTTL, err := parseDurationEnv("ACCESS_TOKEN_TTL", 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshTokenTTL, err := parseDurationEnv("REFRESH_TOKEN_TTL", 30*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &Config{
		Port: getEnv("PORT", "8080"),
		DB: DB{
			Name:     getEnv("DB_NAME", "mococoplan"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			Host:     getEnv("DB_HOST", "mysql"),
			Port:     getEnv("DB_Port", "3306"),
		},
		Token: Token{
			AccessTokenSecret: os.Getenv("JWT_SECRET"),
			AccessTokenIssuer: getEnv("JWT_ISSUER", "mococoplan"),
			AccessTokenTTL:    accessTokenTTL,
			RefreshTokenTTL:   refreshTokenTTL,
		},
		Redis: Redis{
			Host:     getEnv("REDIS_HOST", "redis"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
	}, nil
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}

func parseDurationEnv(key string, def time.Duration) (time.Duration, error) {
	val := os.Getenv(key)
	if val == "" {
		return def, nil
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		return 0, err
	}
	if duration <= 0 {
		return 0, fmt.Errorf("ttl must be positive duration")
	}

	return duration, nil
}
