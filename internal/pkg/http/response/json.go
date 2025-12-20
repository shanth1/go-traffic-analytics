package response

import (
	"encoding/json"
	"net/http"
)

type Envelope map[string]interface{}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

type ErrorResponse struct {
	Error string `json:"error" example:"message"`
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Envelope{"error": message})
}
