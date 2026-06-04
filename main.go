package main

import (
	"log"
	"net/http"

	"grsai-newapi-go/config"
	"grsai-newapi-go/grsai"
	"grsai-newapi-go/router"
)

func main() {
	cfg := config.Load()

	if cfg.GrSAIAPIKey == "" {
		log.Fatal("GRSAI_API_KEY environment variable is required")
	}

	client := grsai.NewClient(cfg.GrSAIBaseURL, cfg.GrSAIAPIKey)
	handler := router.New(client)

	addr := ":" + cfg.ServerPort
	log.Printf("Starting server on %s", addr)
	log.Printf("grsai backend: %s", cfg.GrSAIBaseURL)
	log.Fatal(http.ListenAndServe(addr, handler))
}
