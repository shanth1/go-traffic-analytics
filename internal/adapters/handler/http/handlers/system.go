package handlers

import (
	"net/http"

	"github.com/shanth1/gotrace/internal/pkg/http/response"
)

type HealthResponse struct {
	Status string `json:"status" example:"OK"`
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Check if the server is running
// @Tags         System
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, HealthResponse{Status: "OK"})
}
