package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"grsai-newapi-go/grsai"
	"grsai-newapi-go/model"
	"grsai-newapi-go/storage"
)

type ImageHandler struct {
	Client       *grsai.Client
	ImageStorage *storage.RustFSClient
}

// HandleGenerations handles POST /v1/images/generations (text-to-image).
func (h *ImageHandler) HandleGenerations(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeInvalidRequestError(w, "failed to read request body", "")
		return
	}
	defer r.Body.Close()

	var req model.ImageGenerationRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeInvalidRequestError(w, "invalid request: "+err.Error(), "")
		return
	}

	if req.Prompt == "" {
		writeInvalidRequestError(w, "Missing required parameter: 'prompt'.", "prompt")
		return
	}

	if req.Model == "" {
		req.Model = "gpt-image-2"
	}
	if req.N == 0 {
		req.N = 1
	}
	if req.N != 1 {
		writeInvalidRequestError(w, "Parameter 'n' is not supported for this image backend. Only n=1 is supported.", "n")
		return
	}
	if !validImageResponseFormat(req.ResponseFormat) {
		writeInvalidRequestError(w, "Parameter 'response_format' must be one of: 'url', 'b64_json'.", "response_format")
		return
	}
	if req.Size == "" {
		req.Size = "1024x1024"
	}

	grsaiSize, aspectRatio := mapSize(req.Size)

	drawReq := &model.DrawRequest{
		Model:       req.Model,
		Prompt:      req.Prompt,
		Size:        grsaiSize,
		AspectRatio: aspectRatio,
		Quality:     req.Quality,
	}

	result, err := h.Client.CreateImage(drawReq)
	if err != nil {
		writeErrorFromUpstream(w, err)
		return
	}

	if err := h.writeImageResult(r.Context(), w, result, req.ResponseFormat); err != nil {
		writeErrorDetail(w, http.StatusInternalServerError, err.Error(), "server_error", "", nil)
	}
}

// HandleEdits handles POST /v1/images/edits (image-to-image).
func (h *ImageHandler) HandleEdits(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")

	var imageData, prompt, modelName, size, quality, responseFormat string
	n := 1

	if strings.HasPrefix(ct, "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			writeInvalidRequestError(w, "failed to parse multipart form: "+err.Error(), "")
			return
		}

		prompt = r.FormValue("prompt")
		modelName = r.FormValue("model")
		size = r.FormValue("size")
		quality = r.FormValue("quality")
		responseFormat = r.FormValue("response_format")
		if rawN := r.FormValue("n"); rawN != "" {
			parsedN, err := strconv.Atoi(rawN)
			if err != nil {
				writeInvalidRequestError(w, "Parameter 'n' must be an integer.", "n")
				return
			}
			n = parsedN
		}

		file, _, err := r.FormFile("image")
		if err != nil {
			writeInvalidRequestError(w, "Missing required parameter: 'image'.", "image")
			return
		}
		defer file.Close()

		imageBytes, err := io.ReadAll(file)
		if err != nil {
			writeInvalidRequestError(w, "failed to read image: "+err.Error(), "image")
			return
		}

		imageData = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageBytes)
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeInvalidRequestError(w, "failed to read request body", "")
			return
		}
		defer r.Body.Close()

		var req model.ImageEditRequest
		if err := json.Unmarshal(body, &req); err != nil {
			writeInvalidRequestError(w, "invalid request: "+err.Error(), "")
			return
		}
		imageData = req.Image
		prompt = req.Prompt
		modelName = req.Model
		if req.N != 0 {
			n = req.N
		}
		size = req.Size
		quality = req.Quality
		responseFormat = req.ResponseFormat
	}

	if imageData == "" {
		writeInvalidRequestError(w, "Missing required parameter: 'image'.", "image")
		return
	}
	if prompt == "" {
		writeInvalidRequestError(w, "Missing required parameter: 'prompt'.", "prompt")
		return
	}
	if n != 1 {
		writeInvalidRequestError(w, "Parameter 'n' is not supported for this image backend. Only n=1 is supported.", "n")
		return
	}
	if !validImageResponseFormat(responseFormat) {
		writeInvalidRequestError(w, "Parameter 'response_format' must be one of: 'url', 'b64_json'.", "response_format")
		return
	}
	if modelName == "" {
		modelName = "nano-banana-pro"
	}
	if size == "" {
		size = "1024x1024"
	}

	grsaiSize, aspectRatio := mapSize(size)

	drawReq := &model.DrawRequest{
		Model:       modelName,
		Prompt:      prompt,
		Image:       imageData,
		Size:        grsaiSize,
		AspectRatio: aspectRatio,
		Quality:     quality,
	}

	result, err := h.Client.EditImage(drawReq)
	if err != nil {
		writeErrorFromUpstream(w, err)
		return
	}

	if err := h.writeImageResult(r.Context(), w, result, responseFormat); err != nil {
		writeErrorDetail(w, http.StatusInternalServerError, err.Error(), "server_error", "", nil)
	}
}

func (h *ImageHandler) writeImageResult(ctx context.Context, w http.ResponseWriter, event *model.DrawSSEEvent, responseFormat string) error {
	var imageDatas []model.ImageData
	for _, item := range event.Results {
		id, err := h.imageDataForURL(ctx, item.URL, responseFormat)
		if err != nil {
			return err
		}
		imageDatas = append(imageDatas, id)
	}

	// Fallback: if results is empty but event.URL has the image
	if len(imageDatas) == 0 && event.URL != "" {
		id, err := h.imageDataForURL(ctx, event.URL, responseFormat)
		if err != nil {
			return err
		}
		imageDatas = append(imageDatas, id)
	}

	if len(imageDatas) == 0 {
		return fmt.Errorf("image generation completed without image data")
	}

	resp := model.ImageResponse{
		Created: now(),
		Data:    imageDatas,
	}
	writeJSON(w, http.StatusOK, resp)
	return nil
}

func (h *ImageHandler) imageDataForURL(ctx context.Context, imageURL, responseFormat string) (model.ImageData, error) {
	if h.ImageStorage == nil {
		if responseFormat == "b64_json" {
			b64, err := toBase64(imageURL)
			if err != nil {
				return model.ImageData{}, err
			}
			return model.ImageData{B64JSON: b64}, nil
		}
		return model.ImageData{URL: imageURL}, nil
	}

	stored, err := h.storeImage(ctx, imageURL)
	if err != nil {
		return model.ImageData{}, err
	}
	if responseFormat == "b64_json" {
		return model.ImageData{B64JSON: base64.StdEncoding.EncodeToString(stored.Bytes)}, nil
	}
	return model.ImageData{URL: stored.URL}, nil
}

func (h *ImageHandler) storeImage(ctx context.Context, imageURL string) (*storage.StoredImage, error) {
	if strings.HasPrefix(imageURL, "data:") {
		contentType, bytes, err := decodeDataURL(imageURL)
		if err != nil {
			return nil, err
		}
		return h.ImageStorage.StoreImageBytes(ctx, bytes, contentType)
	}
	return h.ImageStorage.StoreImageURL(ctx, imageURL)
}

func validImageResponseFormat(format string) bool {
	return format == "" || format == "url" || format == "b64_json"
}

// mapSize converts OpenAI size format to grsai format (1K/2K/4K + aspect_ratio).
// Accepts both formats: "1024x1024" (OpenAI) and "1K"/"2K"/"4K" (grsai native, pass through).
func mapSize(size string) (grsaiSize, aspectRatio string) {
	switch size {
	case "256x256", "512x512", "1024x1024":
		return "1K", "1:1"
	case "1024x1792":
		return "2K", "9:16"
	case "1792x1024":
		return "2K", "16:9"
	case "1K", "2K", "4K":
		return size, ""
	default:
		// Unknown format, pass through and let grsai decide
		return size, ""
	}
}

func toBase64(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("image URL is empty")
	}
	if strings.HasPrefix(url, "data:") {
		_, bytes, err := decodeDataURL(url)
		if err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(bytes), nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image for b64_json response: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to fetch image for b64_json response: status %d", resp.StatusCode)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image for b64_json response: %w", err)
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func decodeDataURL(dataURL string) (string, []byte, error) {
	header, payload, ok := strings.Cut(dataURL, ",")
	if !ok {
		return "", nil, fmt.Errorf("image data URL is invalid")
	}
	if !strings.Contains(header, ";base64") {
		return "", nil, fmt.Errorf("image data URL is not base64 encoded")
	}
	contentType := strings.TrimPrefix(strings.TrimSuffix(header, ";base64"), "data:")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	bytes, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", nil, fmt.Errorf("decode image data URL: %w", err)
	}
	return contentType, bytes, nil
}
