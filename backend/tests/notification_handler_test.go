package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"servixa/internal/auth"
	"servixa/internal/config"
	"servixa/internal/database"
	"servixa/internal/handler"
	"servixa/internal/middleware"
	"servixa/internal/model"
	"servixa/internal/repository"
	"servixa/internal/response"
	"servixa/internal/service"

	"github.com/gin-gonic/gin"
)

func TestNotificationHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Load()
	db := database.Connect(cfg)

	if err := db.AutoMigrate(&model.Notification{}); err != nil {
		t.Fatalf("failed to migrate notification: %v", err)
	}

	const (
		userID      uint = 999993
		otherUserID uint = 999994
	)

	db.Where("user_id IN ?", []uint{userID, otherUserID}).
		Delete(&model.Notification{})

	defer db.Where("user_id IN ?", []uint{userID, otherUserID}).
		Delete(&model.Notification{})

	notifications := []model.Notification{
		{
			UserID:    userID,
			Type:      model.NotificationWorkOrderCreated,
			Title:     "Work order created",
			Message:   "Your work order has been created",
			Reference: "WORK_ORDER:999994",
		},
		{
			UserID:    userID,
			Type:      model.NotificationWorkStarted,
			Title:     "Work started",
			Message:   "Technician has started the work",
			Reference: "WORK_ORDER:999995",
		},
		{
			UserID:    otherUserID,
			Type:      model.NotificationWorkCompleted,
			Title:     "Work completed",
			Message:   "Another user's work is completed",
			Reference: "WORK_ORDER:999996",
		},
	}

	if err := db.Create(&notifications).Error; err != nil {
		t.Fatalf("failed to create notifications: %v", err)
	}

	repo := repository.NewNotificationRepository(db)
	svc := service.NewNotificationService(db, repo)
	h := handler.NewNotificationHandler(svc)

	r := gin.New()
	protected := r.Group("/api/v1")
	protected.Use(middleware.Auth(cfg))

	notification := protected.Group("/notifications")
	{
		notification.GET("", h.GetByUserID)
		notification.GET("/unread", h.GetUnreadByUserID)
		notification.GET("/unread/count", h.CountUnreadByUserID)
		notification.PATCH("/read-all", h.MarkAllAsRead)
		notification.PATCH("/:id/read", h.MarkAsRead)
	}

	createToken := func(id uint) string {
		t.Helper()

		token, err := auth.GenerateAccessToken(
			id,
			model.RoleCustomer,
			cfg.JWTSecret,
			time.Hour,
		)
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		return token
	}

	request := func(method, path, token string) *httptest.ResponseRecorder {
		t.Helper()

		req := httptest.NewRequest(method, path, strings.NewReader(""))
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}

		recorder := httptest.NewRecorder()
		r.ServeHTTP(recorder, req)
		return recorder
	}

	userToken := createToken(userID)
	otherToken := createToken(otherUserID)

	t.Run("rejects request without token", func(t *testing.T) {
		recorder := request(
			http.MethodGet,
			"/api/v1/notifications",
			"",
		)

		if recorder.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401, got %d: %s", recorder.Code, recorder.Body.String())
		}
	})

	t.Run("lists only authenticated user's notifications", func(t *testing.T) {
		recorder := request(
			http.MethodGet,
			"/api/v1/notifications",
			userToken,
		)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var result struct {
			Success bool                           `json:"success"`
			Data    []handler.NotificationResponse `json:"data"`
		}

		if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if !result.Success {
			t.Fatal("expected success response")
		}

		if len(result.Data) != 2 {
			t.Fatalf("expected 2 notifications, got %d", len(result.Data))
		}

		for _, item := range result.Data {
			if item.ID == notifications[2].ID {
				t.Fatal("response contains another user's notification")
			}
		}
	})

	t.Run("returns unread count", func(t *testing.T) {
		recorder := request(
			http.MethodGet,
			"/api/v1/notifications/unread/count",
			userToken,
		)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var result struct {
			Success bool `json:"success"`
			Data    struct {
				Count int64 `json:"count"`
			} `json:"data"`
		}

		if err := json.Unmarshal(recorder.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if result.Data.Count != 2 {
			t.Fatalf("expected unread count 2, got %d", result.Data.Count)
		}
	})

	t.Run("marks own notification as read", func(t *testing.T) {
		path := fmt.Sprintf(
			"/api/v1/notifications/%d/read",
			notifications[0].ID,
		)

		recorder := request(http.MethodPatch, path, userToken)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var saved model.Notification
		if err := db.First(&saved, notifications[0].ID).Error; err != nil {
			t.Fatalf("failed to load notification: %v", err)
		}

		if !saved.IsRead || saved.ReadAt == nil {
			t.Fatal("expected notification to be read with read_at set")
		}
	})

	t.Run("does not allow marking another user's notification", func(t *testing.T) {
		path := fmt.Sprintf(
			"/api/v1/notifications/%d/read",
			notifications[2].ID,
		)

		recorder := request(http.MethodPatch, path, userToken)

		// Current handler returns 500 for service errors.
		// The test checks that another user's data remains unchanged.
		var saved model.Notification
		if err := db.First(&saved, notifications[2].ID).Error; err != nil {
			t.Fatalf("failed to load notification: %v", err)
		}

		if saved.IsRead {
			t.Fatal("another user's notification was unexpectedly marked as read")
		}

		if recorder.Code == http.StatusOK {
			t.Fatal("expected unauthorized notification access to fail")
		}
	})

	t.Run("marks all own notifications as read", func(t *testing.T) {
		recorder := request(
			http.MethodPatch,
			"/api/v1/notifications/read-all",
			userToken,
		)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d: %s", recorder.Code, recorder.Body.String())
		}

		var count int64
		if err := db.Model(&model.Notification{}).
			Where("user_id = ? AND is_read = ?", userID, false).
			Count(&count).Error; err != nil {
			t.Fatalf("failed to count unread notifications: %v", err)
		}

		if count != 0 {
			t.Fatalf("expected no unread notifications, got %d", count)
		}

		var other model.Notification
		if err := db.First(&other, notifications[2].ID).Error; err != nil {
			t.Fatalf("failed to load other user's notification: %v", err)
		}

		if other.IsRead {
			t.Fatal("mark-all changed another user's notification")
		}
	})

	_ = otherToken
	_ = response.Success
}
