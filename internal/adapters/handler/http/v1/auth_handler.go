package v1

import (
	"encoding/json"
	"net/http"

	"github.com/shanth1/gotrace/internal/core/ports"
	"github.com/shanth1/gotrace/internal/pkg/http/response"
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

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Register godoc
// @Summary      Register new user
// @Description  Register a new user account
// @Tags         Auth
// @Security     APIKeyAuth
// @Accept       json
// @Produce      json
// @Param        request body RegisterReq true "Registration info"
// @Success      201  {object}  domain.User "Created user"
// @Failure      400  {object}  response.ErrorResponse
// @Failure      409  {object}  response.ErrorResponse "Email already taken"
// @Router       /api/v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	if req.Email == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "email and password are required")
		return
	}

	user, err := h.service.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		// Conflict обычно используется если email уже занят
		response.Error(w, http.StatusConflict, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, user)
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
// @Failure      401  {object}  response.ErrorResponse
// @Router       /api/v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request")
		return
	}

	token, user, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user":  user,
	})
}
