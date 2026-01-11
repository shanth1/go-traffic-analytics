package v1

import (
	"net/http"

	"github.com/shanth1/gotools/log"
	"github.com/shanth1/gotools/logkeys"
	"github.com/shanth1/gotools/ops"
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
// @Failure 400 {object} response.ErrorWrapper
// @Failure 401 {object} response.ErrorWrapper "Missing authorization header"
// @Failure 409 {object} response.ErrorWrapper "Email already taken"
// @Failure 500 {object} response.ErrorWrapper
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "v1.AuthHandler.Register"

	var req RegisterReq

	if err := request.DecodeJSON(w, r, &req); err != nil {
		err := ops.WrapMsg(op, ops.KindInvalid, err, err.Error())
		response.Error(w, r, err)
		return
	}

	user, err := h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Str(logkeys.UserID, string(user.ID)).Str("email", user.Email).Msg("user_registered")

	response.Created(w, r, user.ToPublic())
}

// Login godoc
// @Summary Login user
// @Description Authenticate and get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginReq true "Credentials"
// @Success 200 {object} LoginResponse "Token and User info"
// @Failure 400 {object} response.ErrorWrapper
// @Failure 401 {object} response.ErrorWrapper "Invalid credentials"
// @Failure 500 {object} response.ErrorWrapper
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "v1.AuthHandler.Login"

	var req LoginReq
	if err := request.DecodeJSON(w, r, &req); err != nil {
		err := ops.WrapMsg(op, ops.KindInvalid, err, err.Error())
		response.Error(w, r, err)
		return
	}

	token, user, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(w, r, err)
		return
	}

	log.FromContext(r.Context()).Info().Str(logkeys.UserID, string(user.ID)).Str("email", user.Email).Msg("user_logged_in")

	response.OK(w, r, LoginData{
		Token: token,
		User:  user.ToPublic(),
	})
}
