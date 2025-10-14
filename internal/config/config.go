package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DatabaseUrl string
	ServerPort  string
	RedisUrl    string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		DatabaseUrl: getEnv(
			"DATABASE_URL",
			"host=localhost user=postgres password=postgres dbname=products port=5432 sslmode=disable"),
		RedisUrl:   getEnv("REDIS_URL", "localhost:6379"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
