package storage

import (
	"context"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestRustFSIntegrationUploadTinyPNG(t *testing.T) {
	loadDotEnv(t)
	if os.Getenv("RUN_RUSTFS_INTEGRATION") != "1" {
		t.Skip("set RUN_RUSTFS_INTEGRATION=1 to upload a tiny PNG to RustFS")
	}

	cfg := RustFSConfig{
		Endpoint:      os.Getenv("RUSTFS_ENDPOINT"),
		AccessKey:     os.Getenv("RUSTFS_ACCESS_KEY"),
		SecretKey:     os.Getenv("RUSTFS_SECRET_KEY"),
		Bucket:        envOr("RUSTFS_BUCKET", "grsai-images"),
		Region:        envOr("RUSTFS_REGION", "us-east-1"),
		Prefix:        envOr("RUSTFS_PREFIX", "grsai-test"),
		PublicBaseURL: os.Getenv("RUSTFS_PUBLIC_BASE_URL"),
		MaxBytes:      1024 * 1024,
	}
	if cfg.Endpoint == "" || cfg.AccessKey == "" || cfg.SecretKey == "" {
		t.Fatal("RUSTFS_ENDPOINT, RUSTFS_ACCESS_KEY, and RUSTFS_SECRET_KEY are required")
	}

	client, err := NewRustFSClient(cfg)
	if err != nil {
		t.Fatalf("NewRustFSClient: %v", err)
	}

	png, err := base64.StdEncoding.DecodeString("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+/p9sAAAAASUVORK5CYII=")
	if err != nil {
		t.Fatalf("decode png fixture: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stored, err := client.StoreImageBytes(ctx, png, "image/png")
	if err != nil {
		t.Fatalf("StoreImageBytes: %v", err)
	}
	if stored.URL == "" {
		t.Fatal("expected stored URL")
	}
	t.Logf("stored tiny png: %s", stored.URL)
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func loadDotEnv(t *testing.T) {
	t.Helper()
	_ = godotenv.Load(".env", "../.env")
}
