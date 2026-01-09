package response

import (
	"encoding/json"
	"net/http"

	"github.com/shanth1/gotools/log"
)

type DataResponse[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Envelope map[string]any

func JSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json")

	if data == nil {
		w.WriteHeader(status)
		return
	}

	buf, err := json.Marshal(data)
	if err != nil {
		log.FromContext(r.Context()).Error().Err(err).Msg("response_marshal_failure")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	if _, err := w.Write(buf); err != nil {
		log.FromContext(r.Context()).Error().Err(err).Msg("response_write_failure")
	}
}

func SuccessData[T any](w http.ResponseWriter, r *http.Request, data T) {
	JSON(w, r, http.StatusOK, DataResponse[T]{Data: data})
}

func Created[T any](w http.ResponseWriter, r *http.Request, data T) {
	JSON(w, r, http.StatusCreated, DataResponse[T]{Data: data})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// ClientError sends a 4xx response. It does NOT log an error internally.
func ClientError(w http.ResponseWriter, r *http.Request, status int, message string) {
	JSON(w, r, status, ErrorResponse{Error: message})
}

// ServerError logs the internal error and sends a generic 500 response to the client.
func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.FromContext(r.Context()).Error().Err(err).Msg("internal_server_error")

	ClientError(w, r, http.StatusInternalServerError, "Internal server error")
}
