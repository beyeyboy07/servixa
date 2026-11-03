package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"servixa/internal/model"
	"servixa/internal/response"
	"servixa/internal/service"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(
	notificationService *service.NotificationService,
) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

type NotificationResponse struct {
	ID        uint       `json:"id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Message   string     `json:"message"`
	Reference string     `json:"reference"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}

func toNotificationResponse(
	notification model.Notification,
) NotificationResponse {
	return NotificationResponse{
		ID:        notification.ID,
		Type:      string(notification.Type),
		Title:     notification.Title,
		Message:   notification.Message,
		Reference: notification.Reference,
		IsRead:    notification.IsRead,
		ReadAt:    notification.ReadAt,
		CreatedAt: notification.CreatedAt,
	}
}

// GetByUserID returns all notifications belonging to the authenticated user.
func (h *NotificationHandler) GetByUserID(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok || userID == 0 {
		response.Error(
			c,
			http.StatusUnauthorized,
			"invalid user identity",
			"UNAUTHORIZED",
		)
		return
	}

	notifications, err := h.notificationService.GetByUserID(userID)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to retrieve notifications",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	data := make([]NotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		data = append(data, toNotificationResponse(notification))
	}

	response.Success(
		c,
		http.StatusOK,
		"notifications retrieved successfully",
		data,
	)
}

// GetUnreadByUserID returns unread notifications belonging to the authenticated user.
func (h *NotificationHandler) GetUnreadByUserID(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok || userID == 0 {
		response.Error(
			c,
			http.StatusUnauthorized,
			"invalid user identity",
			"UNAUTHORIZED",
		)
		return
	}

	notifications, err := h.notificationService.GetUnreadByUserID(userID)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to retrieve unread notifications",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	data := make([]NotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		data = append(data, toNotificationResponse(notification))
	}

	response.Success(
		c,
		http.StatusOK,
		"unread notifications retrieved successfully",
		data,
	)
}

// CountUnreadByUserID returns the unread notification count.
func (h *NotificationHandler) CountUnreadByUserID(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok || userID == 0 {
		response.Error(
			c,
			http.StatusUnauthorized,
			"invalid user identity",
			"UNAUTHORIZED",
		)
		return
	}

	count, err := h.notificationService.CountUnreadByUserID(userID)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to count unread notifications",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"unread notification count retrieved successfully",
		UnreadCountResponse{
			Count: count,
		},
	)
}

// MarkAsRead marks one notification as read.
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid notification id",
			"VALIDATION_ERROR",
		)
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok || userID == 0 {
		response.Error(
			c,
			http.StatusUnauthorized,
			"invalid user identity",
			"UNAUTHORIZED",
		)
		return
	}

	if err := h.notificationService.MarkAsRead(uint(id), userID); err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to mark notification as read",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"notification marked as read",
		nil,
	)
}

// MarkAllAsRead marks all notifications belonging to the authenticated user as read.
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(
			c,
			http.StatusUnauthorized,
			"unauthorized",
			"UNAUTHORIZED",
		)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok || userID == 0 {
		response.Error(
			c,
			http.StatusUnauthorized,
			"invalid user identity",
			"UNAUTHORIZED",
		)
		return
	}

	if err := h.notificationService.MarkAllAsRead(userID); err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"failed to mark notifications as read",
			"INTERNAL_SERVER_ERROR",
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		"all notifications marked as read",
		nil,
	)
}
