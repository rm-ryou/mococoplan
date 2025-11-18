package config

import "os"

type (
	Config struct {
		Port string
		DB   DB
	}

	DB struct {
		Name     string
		User     string
		Password string
		Host     string
		Port     string
	}
)

func NewConfig() *Config {
	return &Config{
		Port: getEnv("PORT", "8080"),
		DB: DB{
			Name:     getEnv("DB_NAME", "mococoplan"),
			User:     getEnv("DB_USER", "mococoplan"),
			Password: getEnv("DB_PASSWORD", "mococoplan"),
			Host:     getEnv("DB_HOST", "mococoplan"),
			Port:     getEnv("DB_Port", "mococoplan"),
		},
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}
