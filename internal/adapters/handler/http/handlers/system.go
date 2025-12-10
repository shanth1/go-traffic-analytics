package handlers

import (
	"net/http"

	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

type HealthResponse struct {
	Status string `json:"status" example:"OK"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, HealthResponse{Status: "OK"})
}
