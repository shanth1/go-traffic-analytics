package response

import (
	"encoding/json"
	"net/http"

	"github.com/shanth1/gotools/log"
)

type Envelope map[string]interface{}

type DataResponse[T any] struct {
	Data T
}

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

func ErrorWithLog(w http.ResponseWriter, r *http.Request, err error, status int, userMessage string) {
	logger := log.FromContext(r.Context())
	logger.Error().Err(err).Msg("request failed")
	Error(w, status, userMessage)
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	ErrorWithLog(w, r, err, http.StatusInternalServerError, "Internal server error")
}
