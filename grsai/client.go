package grsai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type APIError struct {
	StatusCode int
	Message    string
	Type       string
	Param      string
	Code       any
	Body       string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("grsai api error (status %d)", e.StatusCode)
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 180 * time.Second,
		},
	}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

func (c *Client) postJSON(path string, body interface{}) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	return c.do(req)
}

func (c *Client) postRawJSON(path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	return c.do(req)
}

func (c *Client) postJSONStream(path string, body interface{}) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(b))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	// Disable timeout for streaming — the http.Client timeout is for the full round-trip,
	// which would kill long-lived streams. Use a dedicated client with no timeout.
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	streamClient := &http.Client{Timeout: 0}
	return streamClient.Do(req)
}

func (c *Client) postRawJSONStream(path string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	streamClient := &http.Client{Timeout: 0}
	return streamClient.Do(req)
}

func readBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func parseAPIError(statusCode int, body []byte, fallback string) *APIError {
	err := &APIError{
		StatusCode: statusCode,
		Message:    fallback,
		Type:       errorTypeForStatus(statusCode),
		Body:       string(body),
	}

	var openAIResp struct {
		Error struct {
			Message string `json:"message"`
			Type    string `json:"type"`
			Param   string `json:"param"`
			Code    any    `json:"code"`
		} `json:"error"`
	}
	if json.Unmarshal(body, &openAIResp) == nil && openAIResp.Error.Message != "" {
		err.Message = openAIResp.Error.Message
		if openAIResp.Error.Type != "" {
			err.Type = openAIResp.Error.Type
		}
		err.Param = openAIResp.Error.Param
		err.Code = openAIResp.Error.Code
		return err
	}

	var grsaiResp struct {
		Code any    `json:"code"`
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}
	if json.Unmarshal(body, &grsaiResp) == nil && grsaiResp.Msg != "" {
		err.Message = grsaiResp.Msg
		err.Code = grsaiResp.Code
		return err
	}

	var generic map[string]any
	if json.Unmarshal(body, &generic) == nil {
		for _, key := range []string{"message", "error", "detail"} {
			if value, ok := generic[key].(string); ok && value != "" {
				err.Message = value
				return err
			}
		}
	}

	if len(body) > 0 {
		err.Message = string(body)
	}
	return err
}

func errorTypeForStatus(statusCode int) string {
	switch statusCode {
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
		if statusCode >= 500 {
			return "server_error"
		}
		return "api_error"
	}
}
