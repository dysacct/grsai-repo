package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"grsai-newapi-go/grsai"
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

func writeRawJSON(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeErrorDetail(w, status, message, errorTypeForStatus(status), "", nil)
}

func writeInvalidRequestError(w http.ResponseWriter, message, param string) {
	writeErrorDetail(w, http.StatusBadRequest, message, "invalid_request_error", param, nil)
}

func writeErrorFromUpstream(w http.ResponseWriter, err error) {
	var apiErr *grsai.APIError
	if errors.As(err, &apiErr) {
		writeErrorDetail(w, apiErr.StatusCode, apiErr.Message, apiErr.Type, apiErr.Param, apiErr.Code)
		return
	}
	writeErrorDetail(w, http.StatusInternalServerError, err.Error(), "server_error", "", nil)
}

func writeErrorDetail(w http.ResponseWriter, status int, message, errorType, param string, code any) {
	writeJSON(w, status, model.ErrorResponse{
		Error: model.ErrorDetail{
			Message: message,
			Type:    errorType,
			Param:   param,
			Code:    code,
		},
	})
}

func errorTypeForStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "invalid_request_error"
	case http.StatusUnauthorized:
		return "authentication_error"
	case http.StatusForbidden:
		return "permission_error"
	case http.StatusNotFound:
		return "invalid_request_error"
	case http.StatusTooManyRequests:
		return "rate_limit_error"
	default:
		if status >= 500 {
			return "server_error"
		}
		return "api_error"
	}
}

func now() int64 {
	return time.Now().Unix()
}
