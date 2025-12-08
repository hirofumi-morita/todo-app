package config

import (
	"os"
)

type Config struct {
	DatabaseURL       string
	JWTSecret         string
	HasuraEndpoint    string
	HasuraAdminSecret string
	ServerPort        string
}

func Load() *Config {
	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/todo_db?sslmode=disable"),
		JWTSecret:         getEnv("JWT_SECRET", "your-256-bit-secret-key-change-this-in-production"),
		HasuraEndpoint:    getEnv("HASURA_GRAPHQL_ENDPOINT", "http://localhost:8080/v1/graphql"),
		HasuraAdminSecret: getEnv("HASURA_ADMIN_SECRET", "hasura_admin_secret"),
		ServerPort:        getEnv("SERVER_PORT", "8000"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
