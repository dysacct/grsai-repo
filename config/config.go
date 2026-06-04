package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GrSAIAPIKey  string
	GrSAIBaseURL string
	ServerPort   string
}

func Load() *Config {
	// Load .env file if present; ignore error (no .env file is fine)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		GrSAIAPIKey:  getEnv("GRSAI_API_KEY", ""),
		GrSAIBaseURL: getEnv("GRSAI_BASE_URL", "https://grsai.dakka.com.cn"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
