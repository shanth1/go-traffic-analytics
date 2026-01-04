package v1

import (
	"errors"
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
	"github.com/shanth1/gotrace/internal/pkg/request"
)

type AuthHandler struct {
	service ports.AuthService
}

func NewAuthHandler(s ports.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"secret123"`
}

func (r RegisterReq) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if len(r.Password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginReq) Validate() error {
	if r.Email == "" {
		return errors.New("email is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

// Register godoc
// @Summary      Register new user
// @Description  Register a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterReq true "Registration info"
// @Success      201  {object}  domain.User "Created user"
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse "Missing authorization header"
// @Failure      409  {object}  response.ErrorResponse "Email already taken"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterReq

	if err := request.DecodeJSON(w, r, &req); err != nil {
		response.ClientError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		response.ClientError(w, r, http.StatusConflict, err.Error())
		return
	}

	log.FromContext(r.Context()).Info().Str("user_id", string(user.ID)).Str("email", user.Email).Msg("user_registered")

	response.Created(w, r, user)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate and get JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body LoginReq true "Credentials"
// @Success      200  {object}  LoginResponse "Token and User info"
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse "Invalid credentials"
// @Failure      500  {object}  response.ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	if err := request.DecodeJSON(w, r, &req); err != nil {
		response.ClientError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	token, user, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.ClientError(w, r, http.StatusUnauthorized, "invalid credentials")
		return
	}

	log.FromContext(r.Context()).Info().Str("user_id", string(user.ID)).Str("email", user.Email).Msg("user_logged_in")

	response.Success(w, r, response.Envelope{
		"token": token,
		"user":  user,
	})
}
