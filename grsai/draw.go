package grsai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"grsai-newapi-go/model"
)

// CreateImage sends a text-to-image request and reads the SSE stream until completion.
func (c *Client) CreateImage(req *model.DrawRequest) (*model.DrawSSEEvent, error) {
	path := "/v1/draw/completions"
	if strings.HasPrefix(req.Model, "nano-banana") {
		path = "/v1/draw/nano-banana"
	}
	return c.drawAndWait(path, req)
}

// EditImage sends an image-to-image request and reads the SSE stream until completion.
func (c *Client) EditImage(req *model.DrawRequest) (*model.DrawSSEEvent, error) {
	path := "/v1/draw/completions"
	if strings.HasPrefix(req.Model, "nano-banana") {
		path = "/v1/draw/nano-banana"
	}
	return c.drawAndWait(path, req)
}

func (c *Client) drawAndWait(path string, req *model.DrawRequest) (*model.DrawSSEEvent, error) {
	resp, err := c.postJSON(path, req)
	if err != nil {
		return nil, fmt.Errorf("draw request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, parseAPIError(resp.StatusCode, body, "grsai draw request failed")
	}

	// Read SSE stream, collect the last (terminal) event
	var lastEvent *model.DrawSSEEvent
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		payload := strings.TrimPrefix(line, "data: ")

		var evt model.DrawSSEEvent
		if err := json.Unmarshal([]byte(payload), &evt); err != nil {
			continue
		}
		lastEvent = &evt

		switch evt.Status {
		case "succeeded":
			return &evt, nil
		case "failed", "error":
			message := evt.FailureReason
			if message == "" {
				message = evt.Error
			}
			if message == "" {
				message = "grsai image generation failed"
			}
			return nil, &APIError{
				StatusCode: 500,
				Message:    message,
				Type:       "server_error",
			}
		}
		// "running" — continue reading
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read draw SSE stream: %w", err)
	}

	// Stream ended without explicit succeeded/failed; check if last event has results
	if lastEvent != nil && len(lastEvent.Results) > 0 {
		return lastEvent, nil
	}
	return nil, fmt.Errorf("grsai draw stream ended without result")
}
