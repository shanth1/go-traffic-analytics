package v1

import (
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/request"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

type AuthHandler struct {
	service ports.AuthService
}

func NewAuthHandler(s ports.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// Register godoc
// @Summary Register new user
// @Description Register a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterReq true "Registration info"
// @Success 201 {object} RegisterResponse "Created user data"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse "Missing authorization header"
// @Failure 409 {object} response.ErrorResponse "Email already taken"
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/register [post]
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

	response.CreatedData(w, r, user.ToPublic())
}

// Login godoc
// @Summary Login user
// @Description Authenticate and get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginReq true "Credentials"
// @Success 200 {object} LoginResponse "Token and User info"
// @Failure 400 {object} response.ErrorResponse
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/auth/login [post]
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

	response.SuccessData(w, r, LoginData{
		Token: token,
		User:  user.ToPublic(),
	})
}
