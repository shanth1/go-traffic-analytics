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

// --- Structs ---

type DataResponse[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// friendlyError acts as a container to separate the user-facing message
// from the technical underlying error.
type friendlyError struct {
	UserMsg string
	Cause   error
}

// Error returns the user-facing message.
func (e *friendlyError) Error() string {
	return e.UserMsg
}

// Unwrap returns the original technical error.
func (e *friendlyError) Unwrap() error {
	return e.Cause
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

	// Do not write response if client disconnected
	if status == 499 {
		return
	}

	ClientError(w, r, status, msg)
}

// BadRequest handles client-side validation errors.
// It ensures that BOTH the sanitized user message AND the original technical error are preserved.
func BadRequest(w http.ResponseWriter, r *http.Request, msg string, err error) {
	var finalErr error

	// If a custom message is provided ("Invalid request"), we wrap the original error
	// so the logger can see 'err', but the client sees 'msg'.
	if msg != "" {
		finalErr = &friendlyError{UserMsg: msg, Cause: err}
	} else {
		finalErr = err
	}

	// We treat this as an Invalid operation in the domain
	detailedErr := ops.E("http.decode", ops.KindInvalid, finalErr)

	RespondWithError(w, r, detailedErr)
}

// ClientError sends a raw JSON error response without internal logging.
func ClientError(w http.ResponseWriter, r *http.Request, status int, message string) {
	JSON(w, r, status, ErrorResponse{Error: message})
}

// --- Internal Logic ---

func prepareError(err error) (int, string, log.Level) {
	if errors.Is(err, context.Canceled) {
		return 499, "Request Canceled", log.LevelInfo
	}

	var e *ops.Error
	if errors.As(err, &e) {
		code := kindToStatus(e.Kind)

		// 1. Server Errors (5xx): Always Mask
		if code >= 500 {
			return code, "Internal Server Error", log.LevelError
		}

		// 2. Client Errors (4xx): Show Safe Message
		// If it's a friendlyError (wrapped in ops.Error), this returns UserMsg.
		// If it's a standard error, it returns err.Error().
		clientMsg := getSafeMessage(e)

		// 3. Determine Log Level
		level := log.LevelInfo
		if code == http.StatusUnauthorized || code == http.StatusForbidden {
			level = log.LevelWarn
		}

		return code, clientMsg, level
	}

	// Fallback for unknown errors
	return http.StatusInternalServerError, "Internal Server Error", log.LevelError
}

func logRequestError(r *http.Request, err error, status int, level log.Level) {
	logger := log.FromContext(r.Context())

	eventLogger := logger.With(
		log.Int(logkeys.HTTPStatus, status),
		log.Str(logkeys.HTTPPath, r.URL.Path),
	)

	// INTELLIGENT ERROR LOGGING:
	// If the error is our friendlyError (or wrapped in ops), we want to make sure
	// we log the *Cause* (technical details), not just the "User Message".

	// Standard log.Err(err) calls err.Error(). For friendlyError, that is just "Invalid Request".
	// We want to see: "json: syntax error at offset 5".

	// We create a helper to find the deepest relevant error text or object
	cause := getDeepCause(err)
	if cause != nil {
		eventLogger = eventLogger.With(log.Err(cause))
	} else {
		eventLogger = eventLogger.With(log.Err(err))
	}

	if fe := new(friendlyError); errors.As(err, &fe) {
		eventLogger = eventLogger.With(log.Str("client_msg", fe.UserMsg))
	}

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

func getDeepCause(err error) error {
	var opsErr *ops.Error
	if errors.As(err, &opsErr) {
		if opsErr.Err != nil {
			return getDeepCause(opsErr.Err)
		}
	}

	var fe *friendlyError
	if errors.As(err, &fe) {
		return getDeepCause(fe.Cause)
	}

	return err
}

func getSafeMessage(e *ops.Error) string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Kind.String()
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
