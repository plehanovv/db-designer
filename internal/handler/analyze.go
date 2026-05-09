package handler

import (
	"encoding/json"
	"net/http"

	"db-designer-vkr/internal/service"
)

type AnalyzeRequest struct {
	Text string `json:"text"`
}

func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var req AnalyzeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	result := service.AnalyzeText(req.Text)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
