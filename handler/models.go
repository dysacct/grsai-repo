package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"grsai-newapi-go/model"
)

var models = []model.ModelObject{
	// Chat models
	{ID: "gpt-5.5", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "gpt-5.4", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "gemini-3.1-pro", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "gemini-3.1-flash-lite", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "gemini-3-flash", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "gemini-3-pro", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "gemini-2.5-flash", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "gemini-2.5-pro", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	// Image models
	{ID: "gpt-image-2", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "gpt-image-2-vip", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "nano-banana", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "nano-banana-pro", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "nano-banana-pro-vt", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "nano-banana-2", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "nano-banana-fast", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "nano-banana-pro-cl", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "nano-banana-2-cl", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "nano-banana-2-4k-cl", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "nano-banana-pro-vip", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
	{ID: "nano-banana-pro-4k-vip", Object: "model", Created: 1750000000, OwnedBy: "grsai"},
}

func ModelsHandler(w http.ResponseWriter, r *http.Request) {
	resp := model.ModelListResponse{
		Object: "list",
		Data:   models,
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, model.ErrorResponse{
		Error: model.ErrorDetail{Message: message},
	})
}

func now() int64 {
	return time.Now().Unix()
}
