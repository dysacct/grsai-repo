package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestImageGenerationMissingPromptUsesOpenAIErrorShape(t *testing.T) {
	handler := &ImageHandler{}
	req := httptest.NewRequest(http.MethodPost, "/v1/images/generations", strings.NewReader(`{"model":"gpt-image-2"}`))
	rec := httptest.NewRecorder()

	handler.HandleGenerations(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var body struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Param   string `json:"param"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error.Type != "invalid_request_error" {
		t.Fatalf("unexpected error type: %q", body.Error.Type)
	}
	if body.Error.Param != "prompt" {
		t.Fatalf("unexpected error param: %q", body.Error.Param)
	}
	if body.Error.Message == "" {
		t.Fatal("expected error message")
	}
}

func TestImageGenerationRejectsUnsupportedN(t *testing.T) {
	handler := &ImageHandler{}
	req := httptest.NewRequest(http.MethodPost, "/v1/images/generations", strings.NewReader(`{"prompt":"draw","n":2}`))
	rec := httptest.NewRecorder()

	handler.HandleGenerations(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var body struct {
		Error struct {
			Type  string `json:"type"`
			Param string `json:"param"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error.Type != "invalid_request_error" {
		t.Fatalf("unexpected error type: %q", body.Error.Type)
	}
	if body.Error.Param != "n" {
		t.Fatalf("unexpected error param: %q", body.Error.Param)
	}
}

func TestImageGenerationRejectsInvalidResponseFormat(t *testing.T) {
	handler := &ImageHandler{}
	req := httptest.NewRequest(http.MethodPost, "/v1/images/generations", strings.NewReader(`{
		"prompt": "draw",
		"response_format": "json"
	}`))
	rec := httptest.NewRecorder()

	handler.HandleGenerations(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	var body struct {
		Error struct {
			Type  string `json:"type"`
			Param string `json:"param"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error.Type != "invalid_request_error" {
		t.Fatalf("unexpected error type: %q", body.Error.Type)
	}
	if body.Error.Param != "response_format" {
		t.Fatalf("unexpected error param: %q", body.Error.Param)
	}
}
