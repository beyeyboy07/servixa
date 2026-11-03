package handler

import (
	"net/http"

	"servixa/internal/middleware"
	"servixa/internal/response"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) Check(c *gin.Context) {
	response.Success(
		c,
		http.StatusOK,
		"SERVIXA API is running",
		nil,
	)
}

func (h *HealthHandler) Protected(c *gin.Context) {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	role, err := middleware.GetUserRole(c)
	if err != nil {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"protected endpoint accessed",
		gin.H{
			"user_id": userID,
			"role":    role,
		},
	)
}
