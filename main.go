package main

import (
	"log"
	"net/http"

	"grsai-newapi-go/config"
	"grsai-newapi-go/grsai"
	"grsai-newapi-go/router"
	"grsai-newapi-go/storage"
)

func main() {
	cfg := config.Load()

	if cfg.GrSAIAPIKey == "" {
		log.Fatal("GRSAI_API_KEY environment variable is required")
	}

	client := grsai.NewClient(cfg.GrSAIBaseURL, cfg.GrSAIAPIKey)
	var imageStorage *storage.RustFSClient
	if cfg.ImageStorage.Enabled() {
		var err error
		imageStorage, err = storage.NewRustFSClient(storage.RustFSConfig{
			Endpoint:      cfg.ImageStorage.Endpoint,
			AccessKey:     cfg.ImageStorage.AccessKey,
			SecretKey:     cfg.ImageStorage.SecretKey,
			Bucket:        cfg.ImageStorage.Bucket,
			Region:        cfg.ImageStorage.Region,
			Prefix:        cfg.ImageStorage.Prefix,
			PublicBaseURL: cfg.ImageStorage.PublicBaseURL,
			MaxBytes:      cfg.ImageStorage.MaxBytes,
			PublicRead:    cfg.ImageStorage.PublicRead,
		})
		if err != nil {
			log.Fatalf("invalid RustFS image storage config: %v", err)
		}
		log.Printf("image storage enabled: %s/%s", cfg.ImageStorage.Endpoint, cfg.ImageStorage.Bucket)
	}

	handler := router.New(client, cfg.ProxyAPIKey, imageStorage)

	addr := ":" + cfg.ServerPort
	log.Printf("Starting server on %s", addr)
	log.Printf("grsai backend: %s", cfg.GrSAIBaseURL)
	log.Fatal(http.ListenAndServe(addr, handler))
}
