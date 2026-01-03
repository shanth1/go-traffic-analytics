package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/shanth1/gotrace/internal/core/domain"
)

type Validator interface {
	Validate() error
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	if r.Header.Get("Content-Type") != "" {
		if val := r.Header.Get("Content-Type"); !strings.HasPrefix(val, "application/json") {
			return errors.New("content-type must be application/json")
		}
	}

	// Prevent large payloads (DoS protection)
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("badly-formed json (at position %d)", syntaxError.Offset)
		case errors.As(err, &unmarshalTypeError):
			return fmt.Errorf("invalid value for field %q (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			return fmt.Errorf("unknown field %s", strings.TrimPrefix(err.Error(), "json: unknown field "))
		case errors.Is(err, io.EOF):
			return errors.New("request body must not be empty")
		case err.Error() == "http: request body too large":
			return errors.New("body too large (max 1MB)")
		default:
			return err
		}
	}

	if v, ok := dst.(Validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func GetUserClaims(r *http.Request) *domain.JwtCustomClaims {
	if claims, ok := r.Context().Value(domain.CtxKeyUser).(*domain.JwtCustomClaims); ok {
		return claims
	}
	return nil
}

func GetUserID(r *http.Request) domain.UserID {
	if claims := GetUserClaims(r); claims != nil {
		return claims.UserID
	}
	return ""
}
