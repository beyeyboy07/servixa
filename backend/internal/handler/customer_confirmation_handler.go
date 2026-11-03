package handler

import (
	"net/http"
	"strconv"

	"servixa/internal/model"
	"servixa/internal/service"

	"github.com/gin-gonic/gin"
)

type CustomerConfirmationHandler struct {
	service *service.CustomerConfirmationService
}

func NewCustomerConfirmationHandler(
	service *service.CustomerConfirmationService,
) *CustomerConfirmationHandler {
	return &CustomerConfirmationHandler{
		service: service,
	}
}

type CustomerConfirmationRequest struct {
	Status model.CustomerConfirmationStatus `json:"status" binding:"required"`
	Notes  string                           `json:"notes"`
}

func (h *CustomerConfirmationHandler) Confirm(c *gin.Context) {
	workOrderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || workOrderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid work order ID",
			"data":    nil,
		})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid user identity",
			"data":    nil,
		})
		return
	}

	var req CustomerConfirmationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"data":    nil,
		})
		return
	}

	confirmation, err := h.service.Confirm(
		uint(workOrderID),
		userID,
		req.Status,
		req.Notes,
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
		"message": "customer confirmation submitted successfully",
		"data":    confirmation,
	})
}
