package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"grsai-newapi-go/grsai"
)

type ChatHandler struct {
	Client *grsai.Client
}

func (h *ChatHandler) Handle(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeInvalidRequestError(w, "failed to read request body", "")
		return
	}
	defer r.Body.Close()

	var req struct {
		Stream bool `json:"stream,omitempty"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		writeInvalidRequestError(w, "invalid request: "+err.Error(), "")
		return
	}

	if req.Stream {
		h.handleStream(w, r, body)
	} else {
		h.handleNormal(w, r, body)
	}
}

func (h *ChatHandler) handleNormal(w http.ResponseWriter, r *http.Request, body []byte) {
	resp, err := h.Client.ChatCompletion(body)
	if err != nil {
		writeErrorFromUpstream(w, err)
		return
	}
	writeRawJSON(w, http.StatusOK, resp)
}

func (h *ChatHandler) handleStream(w http.ResponseWriter, r *http.Request, body []byte) {
	resp, err := h.Client.ChatCompletionStream(body)
	if err != nil {
		writeErrorFromUpstream(w, err)
		return
	}
	defer resp.Body.Close()

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeErrorDetail(w, http.StatusInternalServerError, "streaming not supported", "server_error", "", nil)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return
			}
			flusher.Flush()
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("stream read error: %v\n", err)
			break
		}
	}
}
