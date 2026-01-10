package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/logkeys"
	"github.com/shanth1/gotools/ops"
)

// --- Internal Wrappers ---

// ResponseWrapper standarizes the JSON structure for success responses.
type ResponseWrapper[T any] struct {
	Data T `json:"data"`
}

// ErrorWrapper standarizes the JSON structure for error responses.
type ErrorWrapper struct {
	Error string `json:"error"`
}

// --- Success Handlers ---

// JSON writes a standardized JSON response with status code and data.
func JSON[T any](w http.ResponseWriter, r *http.Request, status int, data T) {
	if err := r.Context().Err(); err != nil {
		log.FromContext(r.Context()).
			Warn().
			Err(err).
			Msg("response_write_aborted_client_disconnected")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.FromContext(r.Context()).Error().Err(err).Msg("response_json_encode_failed")
	}
}

// OK sends a 200 OK response.
func OK[T any](w http.ResponseWriter, r *http.Request, data T) {
	JSON(w, r, http.StatusOK, ResponseWrapper[T]{Data: data})
}

// Created sends a 201 Created response wrapped in {"data": ...}.
func Created[T any](w http.ResponseWriter, r *http.Request, data T) {
	JSON(w, r, http.StatusCreated, ResponseWrapper[T]{Data: data})
}

// NoContent sends a 204 No Content response.
func NoContent(w http.ResponseWriter, r *http.Request) {
	JSON(w, r, http.StatusNoContent, struct{}{})
}

// --- Error Handlers ---

// Error translates an error into a proper HTTP status and JSON response.
// It handles logging automatically using the request context.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		log.FromContext(r.Context()).
			Info().
			Str(logkeys.Reason, "client_canceled_request").
			Msg("request_terminated")
		return
	}

	status, userMsg, logLevel := analyzeError(err)

	log.FromContext(r.Context()).
		WithLevel(logLevel).
		Err(err). // ops.Error.Error() provides "op: kind: inner_err" trace
		Int(logkeys.HTTPStatus, status).
		Str(logkeys.HTTPPath, r.URL.Path).
		Msg("http_request_error")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if encodeErr := json.NewEncoder(w).Encode(ErrorWrapper{Error: userMsg}); encodeErr != nil {
		log.FromContext(r.Context()).
			Error().
			Err(encodeErr).
			Msg("response_error_encode_failed")
	}
}

// --- Internal Logic ---

// analyzeError extracts HTTP metadata from the error.
// Returns: HTTP Status, Safe User Message, Log Level.
func analyzeError(err error) (int, string, log.Level) {
	// Default: Internal Server Error (500)
	status := http.StatusInternalServerError
	msg := "Internal Server Error"
	level := log.LevelError

	// Unwrap ops.Error
	var opErr *ops.Error
	if errors.As(err, &opErr) {
		status = kindToStatus(opErr.Kind)

		// Determine User Message
		if opErr.Message != "" {
			// Explicit safe message provided by domain logic
			msg = opErr.Message
		} else if status < 500 {
			// For client errors (4xx) without explicit message,
			// it is safe to show the Kind string (e.g. "not_found", "invalid_input").
			msg = opErr.Kind.String()
		}
		// For 5xx, we keep "Internal Server Error" to avoid leaking stack traces/SQL errors.

		// Determine Log Level
		level = statusToLevel(status)
	}

	return status, msg, level
}

// kindToStatus maps domain error Kinds to HTTP status codes.
func kindToStatus(k ops.Kind) int {
	switch k {
	case ops.KindInvalid:
		return http.StatusBadRequest
	case ops.KindUnauthorized:
		return http.StatusUnauthorized
	case ops.KindPermission:
		return http.StatusForbidden
	case ops.KindNotFound:
		return http.StatusNotFound
	case ops.KindExist:
		return http.StatusConflict
	case ops.KindUnavailable:
		return http.StatusServiceUnavailable
	case ops.KindTimeout:
		return http.StatusGatewayTimeout
	case ops.KindNotImplemented:
		return http.StatusNotImplemented
	case ops.KindInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// statusToLevel determines the appropriate log level based on the status code.
func statusToLevel(status int) log.Level {
	switch {
	case status >= 500:
		// Server faults are always Errors
		return log.LevelError
	case status == http.StatusNotFound:
		// 404 is usually not an "error" in the system sense, just traffic.
		return log.LevelInfo
	case status >= 400:
		// Other 4xx (400, 401, 403, 409) are Warnings (client faults).
		return log.LevelWarn
	default:
		return log.LevelInfo
	}
}
