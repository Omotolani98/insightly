package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     string

	RedisHost string
	RedisPort string

	LLMHost   string
	LLMPort   string
	EngineID  string
	ModelName string
}

func LoadConfig() (*Config, error) {
	// Load .env into process
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: No .env file found, relying on system environment...")
	}

	return &Config{
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresUser:     getEnv("POSTGRES_USER", "postgres"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		PostgresDB:       getEnv("POSTGRES_DB", "insightly"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),

		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),

		LLMHost:   getEnv("LLM_HOST", "localhost"),
		LLMPort:   getEnv("LLM_PORT", "12434"),
		EngineID:  getEnv("ENGINE_ID", "llama.cpp"),
		ModelName: getEnv("MODEL_NAME", "ai/llama3.2"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
