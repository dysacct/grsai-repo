package grsai

import (
	"fmt"
	"net/http"
)

// ChatCompletion sends a chat completion request to grsai.
func (c *Client) ChatCompletion(body []byte) ([]byte, error) {
	resp, err := c.postRawJSON("/v1/chat/completions", body)
	if err != nil {
		return nil, fmt.Errorf("chat request: %w", err)
	}

	respBody, err := readBody(resp)
	if err != nil {
		return nil, fmt.Errorf("read chat response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseAPIError(resp.StatusCode, respBody, "grsai chat request failed")
	}
	return respBody, nil
}

// ChatCompletionStream sends a streaming chat request and returns the raw response body for SSE relay.
func (c *Client) ChatCompletionStream(body []byte) (*http.Response, error) {
	resp, err := c.postRawJSONStream("/v1/chat/completions", body)
	if err != nil {
		return nil, fmt.Errorf("chat stream request: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := readBody(resp)
		return nil, parseAPIError(resp.StatusCode, body, "grsai chat stream request failed")
	}
	return resp, nil
}
