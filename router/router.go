package router

import (
	"encoding/json"
	"net/http"

	"grsai-newapi-go/grsai"
	"grsai-newapi-go/handler"
	"grsai-newapi-go/model"
	"grsai-newapi-go/storage"
)

func New(client *grsai.Client, proxyAPIKey string, imageStorage *storage.RustFSClient) http.Handler {
	chatHandler := &handler.ChatHandler{Client: client}
	imageHandler := &handler.ImageHandler{Client: client, ImageStorage: imageStorage}

	routes := map[string]map[string]http.HandlerFunc{
		"/v1/models": {
			http.MethodGet: handler.ModelsHandler,
		},
		"/v1/chat/completions": {
			http.MethodPost: chatHandler.Handle,
		},
		"/v1/images/generations": {
			http.MethodPost: imageHandler.HandleGenerations,
		},
		"/v1/images/edits": {
			http.MethodPost: imageHandler.HandleEdits,
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if proxyAPIKey != "" && r.Header.Get("Authorization") != "Bearer "+proxyAPIKey {
			writeError(w, http.StatusUnauthorized, "Incorrect API key provided.", "authentication_error", "", "invalid_api_key")
			return
		}

		route, ok := routes[r.URL.Path]
		if !ok {
			writeError(w, http.StatusNotFound, "Unknown request URL: "+r.URL.Path+".", "invalid_request_error", "", nil)
			return
		}

		h, ok := route[r.Method]
		if !ok {
			writeError(w, http.StatusMethodNotAllowed, "Method "+r.Method+" not allowed for "+r.URL.Path+".", "invalid_request_error", "", nil)
			return
		}

		h(w, r)
	})
}

func writeError(w http.ResponseWriter, status int, message, errorType, param string, code any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.ErrorResponse{
		Error: model.ErrorDetail{
			Message: message,
			Type:    errorType,
			Param:   param,
			Code:    code,
		},
	})
}
