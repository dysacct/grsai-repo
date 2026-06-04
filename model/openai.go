package model

// --- Models API ---

type ModelObject struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type ModelListResponse struct {
	Object string        `json:"object"`
	Data   []ModelObject `json:"data"`
}

// --- Chat API ---

type ChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // string or []ContentPart for vision
}

type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

type ChatCompletionRequest struct {
	Model       string         `json:"model"`
	Messages    []ChatMessage  `json:"messages"`
	Stream      bool           `json:"stream,omitempty"`
	MaxTokens   int            `json:"max_tokens,omitempty"`
	Temperature *float64       `json:"temperature,omitempty"`
	TopP        *float64       `json:"top_p,omitempty"`
	N           int            `json:"n,omitempty"`
	Stop        []string       `json:"stop,omitempty"`
	User        string         `json:"user,omitempty"`
	Tools       interface{}    `json:"tools,omitempty"`
	ToolChoice  interface{}    `json:"tool_choice,omitempty"`
}

// Non-streaming response
type ChatCompletionResponse struct {
	ID                string                 `json:"id"`
	Object            string                 `json:"object"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	Choices           []ChatChoice           `json:"choices"`
	Usage             *Usage                 `json:"usage,omitempty"`
	SystemFingerprint string                 `json:"system_fingerprint,omitempty"`
}

type ChatChoice struct {
	Index                int                   `json:"index"`
	Message              *ChatResponseMessage  `json:"message"`
	FinishReason         string                `json:"finish_reason"`
	ContentFilterResults *ContentFilterResults `json:"content_filter_results,omitempty"`
}

type ChatResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Streaming chunk
type ChatCompletionChunk struct {
	ID                string              `json:"id"`
	Object            string              `json:"object"`
	Created           int64               `json:"created"`
	Model             string              `json:"model"`
	Choices           []ChatChunkChoice   `json:"choices"`
	Usage             *Usage              `json:"usage,omitempty"`
	SystemFingerprint string              `json:"system_fingerprint,omitempty"`
}

type ChatChunkChoice struct {
	Index                int                   `json:"index"`
	Delta                ChatDelta             `json:"delta"`
	FinishReason         *string               `json:"finish_reason"`
	ContentFilterResults *ContentFilterResults `json:"content_filter_results,omitempty"`
}

type ChatDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ContentFilterResults struct {
	Hate      FilterResult `json:"hate"`
	SelfHarm  FilterResult `json:"self_harm"`
	Sexual    FilterResult `json:"sexual"`
	Violence  FilterResult `json:"violence"`
	Jailbreak JailbreakResult `json:"jailbreak,omitempty"`
	Profanity ProfanityResult `json:"profanity,omitempty"`
}

type FilterResult struct {
	Filtered bool `json:"filtered"`
}

type JailbreakResult struct {
	Filtered  bool `json:"filtered"`
	Detected  bool `json:"detected"`
}

type ProfanityResult struct {
	Filtered bool `json:"filtered"`
	Detected bool `json:"detected"`
}

// --- Error ---

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
	Code    string `json:"code,omitempty"`
}

// --- Images API ---

type ImageGenerationRequest struct {
	Model          string `json:"model,omitempty"`
	Prompt         string `json:"prompt"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Style          string `json:"style,omitempty"`
	User           string `json:"user,omitempty"`
}

type ImageEditRequest struct {
	Image          string `json:"image"`
	Prompt         string `json:"prompt"`
	Mask           string `json:"mask,omitempty"`
	Model          string `json:"model,omitempty"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type ImageResponse struct {
	Created int64        `json:"created"`
	Data    []ImageData  `json:"data"`
}

type ImageData struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}
