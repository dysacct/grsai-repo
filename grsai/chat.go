package grsai

import (
	"encoding/json"
	"fmt"
	"net/http"

	"grsai-newapi-go/model"
)

// ChatCompletion sends a chat completion request to grsai.
// If stream=true, returns the raw *http.Response for SSE relay.
func (c *Client) ChatCompletion(req *model.ChatCompletionRequest) (*model.ChatCompletionResponse, error) {
	resp, err := c.postJSON("/v1/chat/completions", req)
	if err != nil {
		return nil, fmt.Errorf("chat request: %w", err)
	}

	body, err := readBody(resp)
	if err != nil {
		return nil, fmt.Errorf("read chat response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("grsai chat error (status %d): %s", resp.StatusCode, string(body))
	}

	var result model.ChatCompletionResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal chat response: %w", err)
	}
	return &result, nil
}

// ChatCompletionStream sends a streaming chat request and returns the raw response body for SSE relay.
func (c *Client) ChatCompletionStream(req *model.ChatCompletionRequest) (*http.Response, error) {
	resp, err := c.postJSONStream("/v1/chat/completions", req)
	if err != nil {
		return nil, fmt.Errorf("chat stream request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := readBody(resp)
		resp.Body.Close()
		return nil, fmt.Errorf("grsai chat stream error (status %d): %s", resp.StatusCode, string(body))
	}
	return resp, nil
}
