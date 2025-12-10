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

type registerReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST /auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request")
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

// POST /auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
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
