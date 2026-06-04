package model

// --- Draw API (image generation) ---

// Request to /v1/draw/completions or /v1/draw/nano-banana
type DrawRequest struct {
	Model       string `json:"model,omitempty"`
	Prompt      string `json:"prompt"`
	Image       string `json:"image,omitempty"`        // base64 or URL for img2img
	Size        string `json:"size,omitempty"`         // e.g. "1K", "2K", "4K"
	Quality     string `json:"quality,omitempty"`      // "low", "medium", "high"
	AspectRatio string `json:"aspect_ratio,omitempty"` // e.g. "1:1", "16:9", "auto"
}

// grsai error response (code/message pattern, returned as JSON, not SSE)
type GrSAIError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// DrawSSEEvent is one SSE data frame from the draw endpoint.
// grsai returns a stream of progress events; the final event (status="succeeded")
// contains the image URL in results.
type DrawSSEEvent struct {
	ID            string          `json:"id"`
	TaskID        string          `json:"task_id"`
	URL           string          `json:"url"`
	Width         int             `json:"width"`
	Height        int             `json:"height"`
	Progress      int             `json:"progress"`
	Status        string          `json:"status"`
	FailureReason string          `json:"failure_reason"`
	Error         string          `json:"error"`
	Results       []DrawImageItem `json:"results"`
	CallbackURL   string          `json:"callback_url"`
	StartTime     int64           `json:"start_time"`
	EndTime       int64           `json:"end_time"`
}

type DrawImageItem struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
