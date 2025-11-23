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
		Port     string
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
	}
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return def
}
