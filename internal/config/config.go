package config

import "os"

type (
	Config struct {
		Port  string
		DB    DB
		Redis Redis
	}

	DB struct {
		Name     string
		User     string
		Password string
		Port     string
	}

	Redis struct {
		Host     string
		Port     string
		Password string
		DB       string
	}
)

func NewConfig() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		DB: DB{
			Name:     getEnv("DB_NAME", "mococoplan"),
			User:     getEnv("DB_USER", "user"),
			Password: getEnv("DB_PASSWORD", "password"),
			Port:     getEnv("DB_Port", "3306"),
		},
		Redis: Redis{
			Host:     getEnv("REDIS_HOST", "redis"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnv("REDIS_DB", "0"),
		},
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}
