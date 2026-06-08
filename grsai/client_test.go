package grsai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestChatCompletionPassesRawJSONAndReturnsRawResponse(t *testing.T) {
	var upstreamBody map[string]any
	client := NewClient("https://example.test", "test-key")
	client.httpClient.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("unexpected authorization header: %q", got)
		}
		if err := json.NewDecoder(r.Body).Decode(&upstreamBody); err != nil {
			t.Fatalf("decode upstream request: %v", err)
		}
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewBufferString(`{"id":"chatcmpl_test","object":"chat.completion","custom_field":{"kept":true}}`)),
			Request:    r,
		}, nil
	})

	resp, err := client.ChatCompletion([]byte(`{"model":"gpt-test","messages":[],"stream_options":{"include_usage":true}}`))
	if err != nil {
		t.Fatalf("ChatCompletion returned error: %v", err)
	}
	if got := upstreamBody["stream_options"].(map[string]any)["include_usage"]; got != true {
		t.Fatalf("stream_options was not preserved: %#v", upstreamBody)
	}
	if !json.Valid(resp) {
		t.Fatalf("response is not valid JSON: %s", string(resp))
	}
	if string(resp) != `{"id":"chatcmpl_test","object":"chat.completion","custom_field":{"kept":true}}` {
		t.Fatalf("unexpected raw response: %s", string(resp))
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

func TestParseAPIErrorPrefersOpenAIErrorShape(t *testing.T) {
	err := parseAPIError(http.StatusTooManyRequests, []byte(`{
		"error": {
			"message": "Rate limit reached.",
			"type": "rate_limit_error",
			"param": null,
			"code": "rate_limit_exceeded"
		}
	}`), "fallback")

	if err.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("unexpected status: %d", err.StatusCode)
	}
	if err.Message != "Rate limit reached." {
		t.Fatalf("unexpected message: %q", err.Message)
	}
	if err.Type != "rate_limit_error" {
		t.Fatalf("unexpected type: %q", err.Type)
	}
	if err.Code != "rate_limit_exceeded" {
		t.Fatalf("unexpected code: %#v", err.Code)
	}
}

func TestParseAPIErrorReadsGRSAIErrorShape(t *testing.T) {
	err := parseAPIError(http.StatusUnauthorized, []byte(`{"code":40101,"msg":"invalid key","data":null}`), "fallback")

	if err.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", err.StatusCode)
	}
	if err.Message != "invalid key" {
		t.Fatalf("unexpected message: %q", err.Message)
	}
	if err.Type != "authentication_error" {
		t.Fatalf("unexpected type: %q", err.Type)
	}
	if err.Code == nil {
		t.Fatal("expected numeric code to be preserved")
	}
}
