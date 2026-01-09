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

// DataResponse wraps a successful payload.
type DataResponse[T any] struct {
	Data T `json:"data"`
}

// ErrorResponse wraps an error message.
type ErrorResponse struct {
	Error string `json:"error"`
}

// --- Success Responses ---

// JSON writes a JSON response with a specific status code.
func JSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.Context().Err(); err != nil {
		log.FromContext(r.Context()).Warn().Err(err).Msg("response_write_aborted_context_canceled")
		return
	}

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

func Success(w http.ResponseWriter, r *http.Request, data any) {
	JSON(w, r, http.StatusOK, data)
}

func SuccessData[T any](w http.ResponseWriter, r *http.Request, data T) {
	JSON(w, r, http.StatusOK, DataResponse[T]{Data: data})
}

func Created(w http.ResponseWriter, r *http.Request, data any) {
	JSON(w, r, http.StatusCreated, data)
}

func CreatedData[T any](w http.ResponseWriter, r *http.Request, data T) {
	JSON(w, r, http.StatusCreated, DataResponse[T]{Data: data})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// --- Error Handling ---

// RespondWithError is the primary entry point for handling errors
func RespondWithError(w http.ResponseWriter, r *http.Request, err error) {
	status, msg, level := prepareError(err)

	logRequestError(r, err, status, level)

	if status == 499 {
		return
	}

	ClientError(w, r, status, msg)
}

// BadRequest is a helper for handler-level validation errors (e.g. invalid JSON body)
func BadRequest(w http.ResponseWriter, r *http.Request, msg string, err error) {
	detailedErr := ops.E("http.decode", ops.KindInvalid, err)

	if msg != "" {
		detailedErr = ops.E("http.decode", ops.KindInvalid, errors.New(msg))
	}

	RespondWithError(w, r, detailedErr)
}

// ClientError sends a JSON error response.
// Note: It is generally better to use RespondWithError, but this is useful for
// static errors where no `error` object exists.
func ClientError(w http.ResponseWriter, r *http.Request, status int, message string) {
	JSON(w, r, status, ErrorResponse{Error: message})
}

// --- Internal Logic ---

// prepareError analyzes the error and returns:
// 1. HTTP Status Code
// 2. Safe Error Message (for the client)
// 3. Log Level (for the server)
func prepareError(err error) (int, string, log.Level) {
	if errors.Is(err, context.Canceled) {
		return 499, "Request Canceled", log.LevelInfo
	}

	var e *ops.Error
	if errors.As(err, &e) {
		code := kindToStatus(e.Kind)

		// Case 1: Server Side Error (5xx)
		if code >= 500 {
			return code, "Internal Server Error", log.LevelError
		}

		// Case 2: Security Issues (401, 403)
		if code == http.StatusUnauthorized || code == http.StatusForbidden {
			return code, e.Error(), log.LevelWarn
		}

		// Case 3: Client Errors (400, 404, 409, etc)
		return code, e.Error(), log.LevelInfo
	}

	return http.StatusInternalServerError, "Internal Server Error", log.LevelError
}

func logRequestError(r *http.Request, err error, status int, level log.Level) {
	logger := log.FromContext(r.Context())

	eventLogger := logger.With(
		log.Err(err),
		log.Int(logkeys.HTTPStatus, status),
		log.Str(logkeys.HTTPPath, r.URL.Path),
	)

	msg := "http_request_error"

	switch level {
	case log.LevelInfo:
		eventLogger.Info().Msg(msg)
	case log.LevelWarn:
		eventLogger.Warn().Msg(msg)
	case log.LevelError:
		eventLogger.Error().Msg(msg)
	default:
		eventLogger.Error().Msg(msg)
	}
}

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
	default:
		return http.StatusInternalServerError
	}
}
