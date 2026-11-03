package handler

import (
	"net/http"
	"strconv"

	"servixa/internal/service"

	"github.com/gin-gonic/gin"
)

type WorkCompletionHandler struct {
	service *service.WorkCompletionService
}

func NewWorkCompletionHandler(
	service *service.WorkCompletionService,
) *WorkCompletionHandler {
	return &WorkCompletionHandler{
		service: service,
	}
}

type CreateWorkCompletionRequest struct {
	Summary         string `json:"summary"`
	WorkPerformed   string `json:"work_performed"`
	Recommendations string `json:"recommendations"`
}

func (h *WorkCompletionHandler) Create(c *gin.Context) {
	workOrderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid work order id",
			"data":    nil,
		})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user not authenticated",
			"data":    nil,
		})
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user id",
			"data":    nil,
		})
		return
	}

	var req CreateWorkCompletionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"data":    nil,
		})
		return
	}

	completion, err := h.service.Create(
		uint(workOrderID),
		userIDUint,
		req.Summary,
		req.WorkPerformed,
		req.Recommendations,
	)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "work completion created successfully",
		"data":    completion,
	})
}
