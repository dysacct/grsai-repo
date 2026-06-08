package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GrSAIAPIKey  string
	GrSAIBaseURL string
	ServerPort   string
	ProxyAPIKey  string
	ImageStorage ImageStorageConfig
}

type ImageStorageConfig struct {
	Endpoint      string
	AccessKey     string
	SecretKey     string
	Bucket        string
	Region        string
	Prefix        string
	PublicBaseURL string
	MaxBytes      int64
	PublicRead    bool
}

func (c ImageStorageConfig) Enabled() bool {
	return c.Endpoint != "" && c.AccessKey != "" && c.SecretKey != "" && c.Bucket != ""
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
		ProxyAPIKey:  getEnv("PROXY_API_KEY", ""),
		ImageStorage: ImageStorageConfig{
			Endpoint:      getEnv("RUSTFS_ENDPOINT", ""),
			AccessKey:     getEnv("RUSTFS_ACCESS_KEY", ""),
			SecretKey:     getEnv("RUSTFS_SECRET_KEY", ""),
			Bucket:        getEnv("RUSTFS_BUCKET", "apiLyncr"),
			Region:        getEnv("RUSTFS_REGION", "us-east-1"),
			Prefix:        getEnv("RUSTFS_PREFIX", "lyncr"),
			PublicBaseURL: getEnv("RUSTFS_PUBLIC_BASE_URL", ""),
			MaxBytes:      getEnvInt64("RUSTFS_MAX_IMAGE_BYTES", 50*1024*1024),
			PublicRead:    getEnvBool("RUSTFS_PUBLIC_READ", false),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err == nil && n > 0 {
			return n
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		b, err := strconv.ParseBool(v)
		if err == nil {
			return b
		}
	}
	return fallback
}
