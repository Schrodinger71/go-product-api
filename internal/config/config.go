package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	ServerPort  string
	RedisUrl    string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		DatabaseUrl: GetsEnv(
			"DATABASE_URL",
			"host=localhost user=postgres password=postgres dbname=products port=5432 sslmode=disable"),
		RedisUrl:   GetsEnv("REDIS_URL", "localhost:6379"),
		ServerPort: GetsEnv("SERVER_PORT", "8080"),
	}
}

func GetsEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
