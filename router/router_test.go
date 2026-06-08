package router

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnknownRouteReturnsOpenAIErrorShape(t *testing.T) {
	router := New(nil, "", nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/unknown", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assertRouterError(t, rec, http.StatusNotFound, "invalid_request_error")
}

func TestWrongMethodReturnsOpenAIErrorShape(t *testing.T) {
	router := New(nil, "", nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/chat/completions", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assertRouterError(t, rec, http.StatusMethodNotAllowed, "invalid_request_error")
}

func TestProxyAPIKeyRequiresBearerToken(t *testing.T) {
	router := New(nil, "proxy-key", nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assertRouterError(t, rec, http.StatusUnauthorized, "authentication_error")
}

func assertRouterError(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, expectedType string) {
	t.Helper()
	if rec.Code != expectedStatus {
		t.Fatalf("unexpected status: got %d want %d", rec.Code, expectedStatus)
	}

	var body struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Error.Type != expectedType {
		t.Fatalf("unexpected error type: %q", body.Error.Type)
	}
	if body.Error.Message == "" {
		t.Fatal("expected error message")
	}
}
