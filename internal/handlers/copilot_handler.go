package handlers

import (
	"encoding/json"
	"net/http"

	"mo-da-backend/internal/intelligence"
)

type CopilotRequestPayload struct {
	Question string                       `json:"question"`
	Context  *intelligence.CopilotContext `json:"context,omitempty"`
}

func AskCopilot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CopilotRequestPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Fallback if empty body
		req.Question = ""
	}

	resp := intelligence.Analyze(req.Question, req.Context)
	JSON(w, resp)
}
