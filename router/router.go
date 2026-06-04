package router

import (
	"net/http"

	"grsai-newapi-go/grsai"
	"grsai-newapi-go/handler"
)

func New(client *grsai.Client) http.Handler {
	mux := http.NewServeMux()

	chatHandler := &handler.ChatHandler{Client: client}
	imageHandler := &handler.ImageHandler{Client: client}

	mux.HandleFunc("GET /v1/models", handler.ModelsHandler)
	mux.HandleFunc("POST /v1/chat/completions", chatHandler.Handle)
	mux.HandleFunc("POST /v1/images/generations", imageHandler.HandleGenerations)
	mux.HandleFunc("POST /v1/images/edits", imageHandler.HandleEdits)

	return mux
}
