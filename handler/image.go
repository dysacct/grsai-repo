package handler

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"grsai-newapi-go/grsai"
	"grsai-newapi-go/model"
)

type ImageHandler struct {
	Client *grsai.Client
}

// HandleGenerations handles POST /v1/images/generations (text-to-image).
func (h *ImageHandler) HandleGenerations(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

	var req model.ImageGenerationRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	if req.Prompt == "" {
		writeError(w, http.StatusBadRequest, "prompt is required")
		return
	}

	if req.Model == "" {
		req.Model = "gpt-image-2"
	}
	if req.N == 0 {
		req.N = 1
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
		writeError(w, http.StatusInternalServerError, "image generation failed: "+err.Error())
		return
	}

	h.writeImageResult(w, result, req.ResponseFormat)
}

// HandleEdits handles POST /v1/images/edits (image-to-image).
func (h *ImageHandler) HandleEdits(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")

	var imageData, prompt, modelName, size, quality, responseFormat string

	if strings.HasPrefix(ct, "multipart/form-data") {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			writeError(w, http.StatusBadRequest, "failed to parse multipart form: "+err.Error())
			return
		}

		prompt = r.FormValue("prompt")
		modelName = r.FormValue("model")
		size = r.FormValue("size")
		quality = r.FormValue("quality")
		responseFormat = r.FormValue("response_format")

		file, _, err := r.FormFile("image")
		if err != nil {
			writeError(w, http.StatusBadRequest, "image file is required: "+err.Error())
			return
		}
		defer file.Close()

		imageBytes, err := io.ReadAll(file)
		if err != nil {
			writeError(w, http.StatusBadRequest, "failed to read image: "+err.Error())
			return
		}

		imageData = "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageBytes)
	} else {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeError(w, http.StatusBadRequest, "failed to read request body")
			return
		}
		defer r.Body.Close()

		var req model.ImageEditRequest
		if err := json.Unmarshal(body, &req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request: "+err.Error())
			return
		}
		imageData = req.Image
		prompt = req.Prompt
		modelName = req.Model
		size = req.Size
		quality = req.Quality
		responseFormat = req.ResponseFormat
	}

	if imageData == "" {
		writeError(w, http.StatusBadRequest, "image is required")
		return
	}
	if prompt == "" {
		writeError(w, http.StatusBadRequest, "prompt is required")
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
		writeError(w, http.StatusInternalServerError, "image edit failed: "+err.Error())
		return
	}

	h.writeImageResult(w, result, responseFormat)
}

func (h *ImageHandler) writeImageResult(w http.ResponseWriter, event *model.DrawSSEEvent, responseFormat string) {
	var imageDatas []model.ImageData
	for _, item := range event.Results {
		id := model.ImageData{}
		if responseFormat == "b64_json" {
			id.B64JSON = toBase64(item.URL)
		} else {
			id.URL = item.URL
		}
		imageDatas = append(imageDatas, id)
	}

	// Fallback: if results is empty but event.URL has the image
	if len(imageDatas) == 0 && event.URL != "" {
		id := model.ImageData{}
		if responseFormat == "b64_json" {
			id.B64JSON = toBase64(event.URL)
		} else {
			id.URL = event.URL
		}
		imageDatas = append(imageDatas, id)
	}

	resp := model.ImageResponse{
		Created: now(),
		Data:    imageDatas,
	}
	writeJSON(w, http.StatusOK, resp)
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

func toBase64(url string) string {
	if url == "" {
		return ""
	}
	if strings.HasPrefix(url, "data:") {
		parts := strings.SplitN(url, ";base64,", 2)
		if len(parts) == 2 {
			return parts[1]
		}
	}
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(bytes)
}
