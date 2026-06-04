package grsai

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		// Try parsing as grsai error
		var errResp model.GrSAIError
		if json.Unmarshal(body, &errResp) == nil && errResp.Code != 0 {
			return nil, fmt.Errorf("grsai draw error: %s (code=%d)", errResp.Msg, errResp.Code)
		}
		return nil, fmt.Errorf("grsai draw error (status %d): %s", resp.StatusCode, string(body))
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
			return nil, fmt.Errorf("grsai generation failed: progress=%d failure_reason=%s error=%s", evt.Progress, evt.FailureReason, evt.Error)
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
